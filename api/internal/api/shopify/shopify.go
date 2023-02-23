package shopify

import (
	"context"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/softcery/shopify-app-template-go/config"
	"github.com/softcery/shopify-app-template-go/internal/entity"
	"github.com/softcery/shopify-app-template-go/internal/service"
	"github.com/softcery/shopify-app-template-go/pkg/logging"
)

type Options struct {
	Config *config.Config
	Logger logging.Logger
}

var _ service.PlatformAPI = (*shopifyAPI)(nil)

type shopifyAPI struct {
	client  *resty.Client
	logger  logging.Logger
	cfg     *config.Config
	retries int
}

func NewAPI(opts Options) *shopifyAPI {
	restyClient := resty.New()

	return &shopifyAPI{
		logger: opts.Logger.Named("shopifyAPI"),
		cfg:    opts.Config,
		client: restyClient,
	}
}
func (s *shopifyAPI) WithConfig(ctx context.Context, store *entity.Store) service.PlatformAPI {
	var h *resty.Client

	if s.retries != 0 {
		c := retryablehttp.NewClient()
		c.RetryMax = s.retries
		c.RetryWaitMax = time.Second * 30
		h = resty.NewWithClient(c.StandardClient())
	} else {
		h = resty.New()
	}

	h = h.
		SetBaseURL(fmt.Sprintf(`https://%s`, store.Name)).
		SetHeader("X-Shopify-Access-Token", store.AccessToken).
		SetHeader("Content-Type", "application/json")

	return &shopifyAPI{
		client:  h,
		logger:  s.logger,
		cfg:     s.cfg,
		retries: s.retries,
	}
}
