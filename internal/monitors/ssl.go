package monitors

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/models"
	"github.com/brunohfonseca/ratatoskr/internal/services"
	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/redis/go-redis/v9"
)

func ProcessSSLCheck(msg redis.XMessage) {
	logger.DebugLog("✅ SSL check started")
	uuid := msg.Values["uuid"].(string)
	domain := msg.Values["domain"].(string)
	timeoutStr, _ := msg.Values["timeout"].(string)

	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		logger.ErrLog("Erro ao converter timeout", err)
		timeout = 30
	}
	sslInfo := FetchSSL(domain, timeout)
	sslInfo.UUID = uuid
	err = services.RegisterSslInfo(sslInfo)
	if err != nil {
		logger.ErrLog("Erro ao registrar SSL info", err)
		return
	}

}

func FetchSSL(domain string, timeout int) models.SSLInfo {
	// Remove protocolo se presente (http:// ou https://)
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")

	// Remove trailing slash se houver
	domain = strings.TrimSuffix(domain, "/")

	addr := domain
	if !hasPort(domain) {
		addr = domain + ":443"
	}
	dialer := &net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", addr,
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		return models.SSLInfo{
			Valid:     models.SSLStatusError,
			Error:     fmt.Sprintf("falha ao conectar: %v", err),
			LastCheck: time.Now(),
		}
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return models.SSLInfo{
			Valid:     models.SSLStatusError,
			Error:     "nenhum certificado encontrado",
			LastCheck: time.Now(),
		}
	}

	cert := certs[0] // pega o primeiro cert da cadeia
	now := time.Now()

	// Determina o status baseado na data de expiração
	var status models.SSLStatus
	var errorMsg string

	if now.After(cert.NotAfter) {
		// Certificado expirado
		status = models.SSLStatusExpired
		errorMsg = "certificado expirado"
	} else if now.Add(30 * 24 * time.Hour).After(cert.NotAfter) {
		// Certificado expira em menos de 30 dias
		status = models.SSLStatusWarning
		daysRemaining := int(cert.NotAfter.Sub(now).Hours() / 24)
		errorMsg = fmt.Sprintf("certificado expira em %d dias", daysRemaining)
	} else {
		// Certificado válido
		status = models.SSLStatusValid
	}

	return models.SSLInfo{
		Valid:          status,
		ExpirationDate: cert.NotAfter,
		Issuer:         cert.Issuer.CommonName,
		Error:          errorMsg,
		LastCheck:      now,
	}
}

func hasPort(host string) bool {
	for _, c := range host {
		if c == ':' {
			return true
		}
	}
	return false
}
