package stripe

import (
	"github.com/stripe/stripe-go/v74"
	"net/http"
)

type Client struct {
	B   stripe.Backend
	Key string
}

func (c Client) CreateProduct(params *stripe.ProductParams) (*stripe.Product, error) {
	product := &stripe.Product{}
	err := c.B.Call(http.MethodPost, "/v1/products", c.Key, params, product)
	return product, err
}

func (c Client) CreatePrice(params *stripe.PriceParams) (*stripe.Price, error) {
	price := &stripe.Price{}
	err := c.B.Call(http.MethodPost, "/v1/prices", c.Key, params, price)
	return price, err
}

func (c Client) CreateCustomer(params *stripe.CustomerParams) (*stripe.Customer, error) {
	customer := &stripe.Customer{}
	err := c.B.Call(http.MethodPost, "/v1/customers", c.Key, params, customer)
	return customer, err
}
