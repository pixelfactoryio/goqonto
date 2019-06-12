package goqonto

import (
	"context"
	"fmt"
	"net/http"
)

// organizationsBasePath Qonto API Organizations Endpoint
const organizationsBasePath = "v2/organizations"

// OrganizationsService interface
// Get: get organizations details
type OrganizationsService interface {
	Get(context.Context, string) (*Organization, *Response, error)
}

// Organization struct
// https://api-doc.qonto.eu/2.0/organizations/show-organization-1
type Organization struct {
	Slug         string        `json:"slug"`
	BankAccounts []BankAccount `json:"bank_accounts"`
}

// BankAccount struct
// https://api-doc.qonto.eu/2.0/organizations/show-organization-1
type BankAccount struct {
	Slug                   string  `json:"slug,omitempty"`
	IBAN                   string  `json:"iban"`
	BIC                    string  `json:"bic"`
	Currency               string  `json:"currency"`
	Balance                float32 `json:"balance"`
	BalanceCents           int     `json:"balance_cents"`
	AuthorizedBalance      float32 `json:"authorized_balance"`
	AuthorizedBalanceCents int     `json:"authorized_balance_cents"`
}

// OrganizationsServiceOp struct used to embed *Client
type OrganizationsServiceOp struct {
	client *Client
}

var _ OrganizationsService = &OrganizationsServiceOp{}

// organizationsRoot root key in the JSON response for organizations
type organizationsRoot struct {
	Organization Organization `json:"organization"`
}

// Get Organization
func (o *OrganizationsServiceOp) Get(ctx context.Context, id string) (*Organization, *Response, error) {

	path := fmt.Sprintf("%s/%s", organizationsBasePath, id)

	req, err := o.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationsRoot)
	resp, err := o.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	return &root.Organization, resp, nil
}
