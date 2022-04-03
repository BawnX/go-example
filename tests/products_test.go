package tests

import (
	"ms/data"
	"testing"
)

func TestChecksValidation(t *testing.T) {
	p := &data.Product{
		Name:  "test",
		Price: 1.00,
		SKU:   "abs-asd-adde",
	}
	err := p.Validate()

	if err != nil {
		t.Fatal(err)
	}
}
