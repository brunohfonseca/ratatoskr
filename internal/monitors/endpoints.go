package monitors

import (
	"context"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog/log"
)

func ProcessEndpoint(ctx context.Context, redisClient *redis.Client, stream, group string, msg redis.XMessage) {
	uuid := msg.Values["uuid"].(string)
	domain := msg.Values["domain"].(string)
	path := msg.Values["path"].(string)
	checkSSLStr, _ := msg.Values["check_ssl"].(string)

	doHealthCheck(domain, path)

	if checkSSLStr == "true" {
		_, err := FetchSSL(domain)
		if err != nil {
			return
		}
	}

	log.Info().Msgf("✅ Health check completed in %s", uuid)

}

func doHealthCheck(domain, path string) (string, int64) {
	return "", 0
}
