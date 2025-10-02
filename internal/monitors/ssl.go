package monitors

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"time"

	"github.com/brunohfonseca/ratatoskr/internal/utils/logger"
	"github.com/redis/go-redis/v9"
)

func ProcessSSLCheck(ctx context.Context, redisClient *redis.Client, stream, group string, msg redis.XMessage) {
	logger.DebugLog("✅ SSL check started")
}

func FetchSSL(domain string) (time.Time, error) {
	host := domain

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	conn, err := tls.DialWithDialer(dialer, "tcp", host, &tls.Config{
		InsecureSkipVerify: true,
	})
	if err != nil {
		return time.Time{}, err
	}
	defer func() {
		if err := conn.Close(); err != nil {
			// Log the error but don't override the main function's return error
			// since connection close errors are typically not critical
			fmt.Printf("Warning: failed to close SSL connection: %v\n", err)
		}
	}()

	certs := conn.ConnectionState().PeerCertificates
	if len(certs) == 0 {
		return time.Time{}, fmt.Errorf("nenhum certificado encontrado")
	}

	return certs[0].NotAfter, nil
}
