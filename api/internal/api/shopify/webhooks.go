package shopify

import (
	"fmt"
	"net/http"

	"github.com/softcery/shopify-app-template-go/internal/service"
)

type subscribeToWebhookRequestBody struct {
	Address string `json:"address"`
	Topic   string `json:"topic"`
	Format  string `json:"format"`
}

func (s *shopifyAPI) SubscribeToAppUninstallWebhook(opts service.SubscribeToAppUninstallWebhookOptions) error {
	logger := s.logger.
		Named("SubscribeToAppUninstallWebhook").
		With("opts", opts)

	var res, err = s.client.R().
		SetBody(map[string]interface{}{
			"webhook": subscribeToWebhookRequestBody{
				Address: opts.RedirectURL,
				Topic:   "app/uninstalled",
				Format:  "json",
			},
		}).
		SetHeaders(map[string]string{
			"Content-Type":           "application/json",
			"X-Shopify-Access-Token": opts.AccessToken,
		}).
		Post(fmt.Sprintf("https://%s/admin/api/2022-04/webhooks.json", opts.StoreName))
	if err != nil {
		logger.Error("failed to subscribe to shopify app/uninstalled webhook", "err", err)
		return fmt.Errorf("failed to subscribe to shopify app/uninstalled webhook: %w", err)
	}
	if res.StatusCode() != http.StatusCreated {
		logger.Error("failed to subscribe to shopify app/uninstalled webhook", "resBody", res.String())
		return fmt.Errorf("failed to subscribe to shopify app/uninstalled webhook: http status %d, body %s", res.StatusCode(), res.String())
	}
	logger = logger.With("resBody", res.Body())

	logger.Info("successfully subscribed to shopify app/uninstalled webhook")
	return nil
}
