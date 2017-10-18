package goqonto

import (
	"encoding/json"
	"fmt"

	"github.com/amine7536/goqonto/context"
)

const organizationsBasePath = "v1/organizations"

// OrganizationsService Service
type OrganizationsService interface {
	Get(context.Context, string) (*Organization, *Response, error)
}

// Organization struct
// {
// 	"organization": {
// 		"slug": "croissant-9134",
// 		"bank_accounts": [
// 			{
// 				"slug": "croissant-bank-account-1",
// 				"iban": "FR7616798000010000004321396",
// 				"bic": "TRZOFR21XXX",
// 				"currency": "EUR",
// 				"balance": 24.94
// 			}
// 		]
// 	}
// }
type Organization struct {
	Slug         string        `json:"slug"`
	BankAccounts []BankAccount `json:"bank_accounts"`
}

// BankAccount struct
type BankAccount struct {
	Slug     string  `json:"slug"`
	IBAN     string  `json:"iban"`
	BIC      string  `json:"bic"`
	Currency string  `json:"currency"`
	Balance  float32 `json:"balance"`
}

// OrganizationsServiceOp Service
type OrganizationsServiceOp struct {
	client *Client
}

var _ OrganizationsService = &OrganizationsServiceOp{}

type organizationsRoot struct {
	Organization Organization `json:"organization"`
}

// Convert Organization to a string
func (o Organization) String() string {
	bytes, err := json.Marshal(o)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// Get Organization
func (s *OrganizationsServiceOp) Get(ctx context.Context, id string) (*Organization, *Response, error) {

	path := fmt.Sprintf("%s/%s", organizationsBasePath, id)

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	return &root.Organization, resp, nil
}
