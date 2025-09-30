package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	postgres "github.com/brunohfonseca/ratatoskr/internal/infrastructure/db/postgres"
	"github.com/brunohfonseca/ratatoskr/internal/utils"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"golang.org/x/oauth2"
)

var (
	oauthConfig  *oauth2.Config
	oidcVerifier *oidc.IDTokenVerifier
)

// InitKeycloak inicializa a configuração do Keycloak
func InitKeycloak() error {
	cfg := config.Get()

	if !cfg.Keycloak.Enabled {
		return nil // SSO desabilitado
	}

	ctx := context.Background()

	provider, err := oidc.NewProvider(ctx, cfg.Keycloak.URL)
	if err != nil {
		return err
	}

	// Configura OAuth2
	oauthConfig = &oauth2.Config{
		ClientID:     cfg.Keycloak.ClientID,
		ClientSecret: cfg.Keycloak.ClientSecret,
		RedirectURL:  cfg.Keycloak.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email"},
	}

	// Verifier para validar tokens
	oidcVerifier = provider.Verifier(&oidc.Config{ClientID: cfg.Keycloak.ClientID})

	return nil
}

// KeycloakLogin redireciona para página de login do Keycloak
func KeycloakLogin(c *gin.Context) {
	cfg := config.Get()

	if !cfg.Keycloak.Enabled {
		responses.ErrorMsg(c, http.StatusNotImplemented, "SSO not enabled")
		return
	}

	if oauthConfig == nil {
		responses.ErrorMsg(c, http.StatusInternalServerError, "Keycloak not initialized")
		return
	}

	// Gera state aleatório para segurança
	state := generateRandomState()

	// Salva state na sessão (simplificado, você pode usar Redis depois)
	c.SetCookie("oauth_state", state, 300, "/", "", false, true)

	// Redireciona para Keycloak
	authURL := oauthConfig.AuthCodeURL(state)
	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// KeycloakCallback processa o callback do Keycloak
func KeycloakCallback(c *gin.Context) {
	cfg := config.Get()

	if !cfg.Keycloak.Enabled {
		responses.ErrorMsg(c, http.StatusNotImplemented, "SSO not enabled")
		return
	}

	if oauthConfig == nil || oidcVerifier == nil {
		responses.ErrorMsg(c, http.StatusInternalServerError, "Keycloak not initialized")
		return
	}

	// Valida state
	savedState, _ := c.Cookie("oauth_state")
	if c.Query("state") != savedState {
		responses.ErrorMsg(c, http.StatusBadRequest, "Invalid state")
		return
	}

	// Troca code por token
	oauth2Token, err := oauthConfig.Exchange(context.Background(), c.Query("code"))
	if err != nil {
		responses.ErrorMsg(c, http.StatusInternalServerError, "Failed to exchange token")
		return
	}

	// Extrai ID token
	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		responses.ErrorMsg(c, http.StatusInternalServerError, "No id_token in response")
		return
	}

	// Verifica ID token
	idToken, err := oidcVerifier.Verify(context.Background(), rawIDToken)
	if err != nil {
		responses.ErrorMsg(c, http.StatusUnauthorized, "Invalid ID token")
		return
	}

	// Extrai claims
	var claims struct {
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		Sub           string `json:"sub"`
	}
	if err := idToken.Claims(&claims); err != nil {
		responses.ErrorMsg(c, http.StatusInternalServerError, "Failed to parse claims")
		return
	}

	// Busca ou cria usuário no banco
	userID, userUUID, err := getOrCreateUser(claims.Email, claims.Name)
	if err != nil {
		log.Error().Err(err).Msg("Failed to get or create user")
		responses.ErrorMsg(c, http.StatusInternalServerError, "Failed to create user account")
		return
	}

	// Gera JWT interno da aplicação
	token, err := utils.GenerateJWT(userID, userUUID, claims.Email)
	if err != nil {
		log.Error().Err(err).Msg("Failed to generate JWT")
		responses.ErrorMsg(c, http.StatusInternalServerError, "Failed to generate token")
		return
	}

	// Retorna o token
	responses.Success(c, http.StatusOK, gin.H{
		"token": token,
		"user": gin.H{
			"id":    userID,
			"uuid":  userUUID,
			"email": claims.Email,
			"name":  claims.Name,
		},
	})
}

// getOrCreateUser busca ou cria um usuário no banco baseado no email do SSO
func getOrCreateUser(email, name string) (int, string, error) {
	db := postgres.PostgresConn

	// Tenta buscar usuário existente por email
	var userID int
	var userUUID string
	var enabled bool

	err := db.QueryRow(
		"SELECT id, uuid, enabled FROM users WHERE email = $1",
		email,
	).Scan(&userID, &userUUID, &enabled)

	if err == nil {
		// Usuário existe
		if !enabled {
			return 0, "", sql.ErrNoRows // Usuário desabilitado
		}
		return userID, userUUID, nil
	}

	if err != sql.ErrNoRows {
		// Erro inesperado
		return 0, "", err
	}

	// Usuário não existe, cria novo
	err = db.QueryRow(`
		INSERT INTO users (email, full_name, enabled, password_hash, auth_provider)
		VALUES ($1, $2, true, '', 'keycloak')
		RETURNING id, uuid`,
		email,
		name,
	).Scan(&userID, &userUUID)

	if err != nil {
		return 0, "", err
	}

	log.Info().
		Int("user_id", userID).
		Str("email", email).
		Msg("New user created via SSO")

	return userID, userUUID, nil
}

// generateRandomState gera um state aleatório para OAuth2
func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
