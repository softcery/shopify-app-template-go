package shopify

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/softcery/shopify-app-template-go/internal/service"
)

func (s *shopifyAPI) VerifySession(ctx context.Context) (*service.VerifySessionOutput, error) {
	logger := s.logger.
		Named("VerifySession").
		WithContext(ctx)

	// get the session token from the header authorization
	tokenStringRaw := ctx.Value("Authorization")
	if tokenStringRaw == "" {
		logger.Info("missing auth token")
		return &service.VerifySessionOutput{IsVerified: false}, errors.New("missing auth token")
	}

	tokenString, ok := tokenStringRaw.(string)
	if !ok {
		logger.Info("malformed auth token")
		return &service.VerifySessionOutput{IsVerified: false}, errors.New("malformed auth token")
	}

	// split Bearer and token
	tokenStringArr := strings.Split(tokenString, " ")
	if len(tokenStringArr) != 2 {
		logger.Info("malformed auth token", "tokenStringArr", tokenStringArr)
		return &service.VerifySessionOutput{IsVerified: false}, errors.New("malformed auth token")
	}

	// get token
	token := tokenStringArr[1]

	isSessionTokenValid, err := s.verifySessionToken(token)
	if err != nil {
		logger.Error("failed to verify session token", "err", err)
		return &service.VerifySessionOutput{IsVerified: false}, err
	}

	err = s.verifySignature(token)
	if err != nil {
		logger.Error("failed to verify signature", "err", err)
		return &service.VerifySessionOutput{IsVerified: false}, err
	}

	storeName, err := s.getStoreName(token)
	if err != nil {
		logger.Error("failed to get store name", "err", err)
		return &service.VerifySessionOutput{IsVerified: false}, err
	}

	if isSessionTokenValid != "" && err == nil {
		return &service.VerifySessionOutput{IsVerified: true, StoreName: storeName}, nil
	}

	return &service.VerifySessionOutput{IsVerified: false}, nil
}

// Claims is a struct that represents the JWT claims payload
type Claims struct {
	Issuer string `json:"iss"`
	Dest   string `json:"dest"`
	Aud    string `json:"aud"`
	Sub    string `json:"sub"`
	jwt.StandardClaims
}

// verifySessionToken verifies the session details of the provided JWT token
// https://shopify.dev/docs/apps/auth/oauth/session-tokens/getting-started#obtain-and-verify-session-details
func (s *shopifyAPI) verifySessionToken(tokenString string) (string, error) {
	// Parse the token
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Shopify.ApiSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %v", err)
	}

	// Verify the exp value
	now := time.Now().UTC().Unix()
	if claims.ExpiresAt <= now {
		return "", errors.New("JWT token has expired")
	}

	// Verify the nbf value
	if claims.NotBefore > now {
		return "", errors.New("JWT token not yet valid")
	}

	// Verify the iss and dest values
	if !strings.Contains(claims.Issuer, claims.Dest) {
		return "", errors.New("JWT token contains incorrect issuer value")
	}

	// Verify the aud value
	if claims.Aud != s.cfg.Shopify.ApiKey {
		return "", errors.New("JWT token contains incorrect audience value")
	}

	// Return the sub value, which is the user ID
	return claims.Sub, nil
}

// verifySignature takes a JWT token and the app's secret and returns an error if the signature is invalid
// https://shopify.dev/docs/apps/auth/oauth/session-tokens/getting-started#verify-the-session-tokens-signature
func (s *shopifyAPI) verifySignature(token string) error {
	// Split the token into its three parts
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return errors.New("Invalid token")
	}

	// Combine the header and payload to generate the message
	message := parts[0] + "." + parts[1]

	// Generate the HMAC
	mac := hmac.New(sha256.New, []byte(s.cfg.Shopify.ApiSecret))
	mac.Write([]byte(message))
	expectedSignature := mac.Sum(nil)

	// Decode the signature
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}

	// Compare the expected signature with the provided signature
	if !hmac.Equal(expectedSignature, signature) {
		return errors.New("Invalid signature")
	}

	return nil
}

func (s *shopifyAPI) getStoreName(token string) (string, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return "", errors.New("Invalid token")
	}

	// Parse the payload
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.cfg.Shopify.ApiSecret), nil
	})
	if err != nil {
		return "", fmt.Errorf("failed to parse JWT token: %v", err)
	}

	return strings.Replace(claims.Dest, "https://", "", 1), nil
}
