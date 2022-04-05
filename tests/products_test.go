package tests

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"ms/data"
	"testing"
)

func TestProductMissingNameReturnsErr(t *testing.T) {
	p := data.Product{
		Price: 1.22,
	}

	v := data.NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 2)
}

func TestProductMissingPriceReturnsErr(t *testing.T) {
	p := data.Product{
		Name:  "abc",
		Price: -1,
	}

	v := data.NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 2)
}

func TestProductInvalidSKUReturnsErr(t *testing.T) {
	p := data.Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc",
	}

	v := data.NewValidation()
	err := v.Validate(p)
	assert.Len(t, err, 1)
}

func TestValidProductDoesNOTReturnsErr(t *testing.T) {
	p := data.Product{
		Name:  "abc",
		Price: 1.22,
		SKU:   "abc-efg-hji",
	}

	v := data.NewValidation()
	err := v.Validate(p)
	assert.Nil(t, err)
}

func TestProductsToJSON(t *testing.T) {
	ps := []*data.Product{
		&data.Product{
			Name: "abc",
		},
	}

	b := bytes.NewBufferString("")
	err := data.ToJSON(ps, b)
	assert.NoError(t, err)
}
