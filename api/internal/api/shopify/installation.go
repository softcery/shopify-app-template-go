package shopify

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/softcery/shopify-app-template-go/internal/service"
)

func (s *shopifyAPI) HandleInstall(opts service.HandleInstallOptions) (service.APIHandleInstallOutput, error) {
	logger := s.logger.
		Named("HandleInstall").
		With("opts", opts)

	storeNonce := s.generateNonce()
	// Build redirection URL
	values := url.Values{
		"client_id":       {s.cfg.Shopify.ApiKey},
		"scope":           {s.cfg.Shopify.Scopes},
		"redirect_uri":    {opts.RedirectURL},
		"state":           {storeNonce},
		"grant_options[]": {"offline"}, // https://shopify.dev/concepts/about-apis/authentication#api-access-modes
	}

	logger = logger.With("values", values)
	logger.Debug("built values")
	return service.APIHandleInstallOutput{
		Nonce:       storeNonce,
		RedirectURL: fmt.Sprintf("https://%s/admin/oauth/authorize?%s", opts.StoreName, values.Encode()),
	}, nil
}

func (s *shopifyAPI) HandleRedirect(opts service.APIHandleRedirectOptions) (string, error) {
	logger := s.logger.
		Named("HandleRedirect").
		With("opts", opts)

	// Verify redirected URL
	parsedURL, err := url.Parse(opts.RedirectedURL)
	if err != nil {
		logger.Info("failed to parse redirected url", "err", err)
		return "", service.ErrHandleRedirectInvalidRedirectedURL
	}
	if !s.verifyNonce(opts.Nonce, parsedURL) {
		logger.Info("nonce is incorrect")
		return "", service.ErrHandleRedirectInvalidRedirectedURL
	}
	logger.Debug("verified redirected url")

	// Getting access token
	query := parsedURL.Query()
	var credentials map[string]string
	res, err := s.client.R().
		SetQueryParams(map[string]string{
			"client_id":     s.cfg.Shopify.ApiKey,
			"client_secret": s.cfg.Shopify.ApiSecret,
			"code":          query.Get("code"),
		}).
		SetResult(&credentials).
		Post(fmt.Sprintf("https://%s/admin/oauth/access_token", opts.StoreName))
	if err != nil {
		logger.Error("failed to get shopifyAPI access token", "err", err)
		return "", fmt.Errorf("failed to get shopifyAPI access token: %w", err)
	}
	if res.StatusCode() != http.StatusOK {
		logger.Error("failed to get shopifyAPI access token", "resBody", res.String())
		return "", fmt.Errorf("failed to get shopifyAPI access token: http status %d, body %s", res.StatusCode(), res.String())
	}
	if credentials["scope"] != s.cfg.Shopify.Scopes {
		logger.Info("scopes are different", "resBody", res.String())
		return "", service.ErrHandleRedirectInvalidScopes
	}
	logger = logger.With("resBody", res.String())
	logger.Info("got credentials")

	return credentials["access_token"], nil
}

// verifyNonce verifies nonce from given url with the actual one.
func (s *shopifyAPI) verifyNonce(actualNonce string, url *url.URL) bool {
	q := url.Query()
	nonce := q.Get("state")
	if nonce == "" {
		return true
	}
	return nonce == actualNonce
}

// generateNonce is used to generate random nonce.
func (s *shopifyAPI) generateNonce() string {
	nonce, err := uuid.NewUUID()
	if err != nil {
		// If nonce couldn't be generated, use default string
		return "jg6yf3JAdB5hvFaG1o"
	}
	return nonce.String()
}
