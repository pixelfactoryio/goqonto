package goqonto

import (
	"context"
	"fmt"
	"net/http"
)

// organizationsBasePath Qonto API Organizations Endpoint
const organizationsBasePath = "v2/organizations"

// OrganizationsService provides access to the organisations in Qonto API
type OrganizationsService service

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

// organizationsRoot root key in the JSON response for organizations
type organizationsRoot struct {
	Organization *Organization `json:"organization"`
}

// Get Organization
func (s *OrganizationsService) Get(ctx context.Context, id string) (*Organization, *Response, error) {

	path := fmt.Sprintf("%s/%s", organizationsBasePath, id)

	req, err := s.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(organizationsRoot)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Organization, resp, nil
}
