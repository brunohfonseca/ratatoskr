package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/brunohfonseca/ratatoskr/internal/config"
	"github.com/brunohfonseca/ratatoskr/internal/utils/responses"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gin-gonic/gin"
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

	// TODO: Buscar ou criar usuário no banco baseado no email
	// userID, userUUID := getOrCreateUser(claims.Email, claims.Name)

	// Gera JWT interno da aplicação
	// token, _ := utils.GenerateJWT(userID, userUUID, claims.Email)

	// Por enquanto, retorna as informações
	responses.Success(c, http.StatusOK, gin.H{
		"email":   claims.Email,
		"name":    claims.Name,
		"sub":     claims.Sub,
		"message": "SSO authentication successful. TODO: Generate internal JWT",
	})
}

// generateRandomState gera um state aleatório para OAuth2
func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}
