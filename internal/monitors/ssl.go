package monitors

import (
	"crypto/tls"
	"fmt"
	"net"
	"strconv"
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

	dialer := &net.Dialer{
		Timeout: time.Duration(timeout) * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", domain,
		&tls.Config{
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		return models.SSLInfo{
			Valid: string(models.SSLStatusUnknown),
			Error: fmt.Sprintf("falha ao conectar: %v", err),
		}
	}
	defer conn.Close()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return models.SSLInfo{
			Valid: string(models.SSLStatusUnknown),
			Error: "nenhum certificado encontrado",
		}
	}

	cert := certs[0] // pega o primeiro cert da cadeia

	switch cert {

	}

	return models.SSLInfo{
		Valid:          string(models.SSLStatusValid),
		ExpirationDate: cert.NotAfter,
		Issuer:         cert.Issuer.CommonName,
	}
}
