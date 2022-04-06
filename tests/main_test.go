package tests

import (
	"ms/sdk/client"
	"ms/sdk/client/products"
	"testing"
)

func TestOurClient(t *testing.T) {
	cfg := client.DefaultTransportConfig().WithHost("localhost:3000")
	c := client.NewHTTPClientWithConfig(nil, cfg)

	params := products.NewListProductsParams()
	prod, err := c.Products.ListProducts(params)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(prod.Payload)
}
