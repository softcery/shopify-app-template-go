package shopify

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/softcery/shopify-app-template-go/internal/service"
)

type createProductRequestBody struct {
	Product *product `json:"product"`
}

type product struct {
	Title       string `json:"title"`
	BodyHTML    string `json:"body_html"`
	Vendor      string `json:"vendor"`
	ProductType string `json:"product_type"`
	Status      string `json:"status"`
}

func (s *shopifyAPI) CreateProducts(ctx context.Context) error {
	logger := s.logger.
		Named("CreateProducts").
		WithContext(ctx)

	for i := 0; i < service.DEFAULT_PRODUCT_COUNT; i++ {
		title := generateRandomProductTitle()
		product := &createProductRequestBody{
			Product: &product{
				Title:       title,
				BodyHTML:    fmt.Sprintf("<p>Product %s</p>", title),
				Vendor:      "Vendor",
				ProductType: "Type",
				Status:      "active",
			},
		}
		res, err := s.client.R().
			SetBody(product).
			Post("/admin/api/2022-07/products.json")
		if err != nil {
			logger.Error("failed to create product", "err", err)
			return err
		}
		if res.StatusCode() != http.StatusCreated {
			logger.Error("failed to create product", "status", res.StatusCode())
			return fmt.Errorf("failed to create product")
		}
	}
	logger.Info("created products")

	return nil
}

func generateRandomProductTitle() string {
	return fmt.Sprintf("Product %d", rand.Intn(1000))
}

func (s *shopifyAPI) GetProductsCount(ctx context.Context) (int, error) {
	logger := s.logger.
		Named("GetProductsCount").
		WithContext(ctx)

	var responseBody struct {
		Count int `json:"count"`
	}

	resp, err := s.client.R().
		SetResult(&responseBody).
		Get("/admin/api/2022-07/products/count.json")
	if err != nil {
		logger.Error("failed to get products count", "err", err)
		return 0, err
	}

	if resp.StatusCode() != 200 {
		logger.Error("failed to get products count", "status", resp.StatusCode())
		return 0, fmt.Errorf("failed to get products count")
	}

	return responseBody.Count, nil
}
