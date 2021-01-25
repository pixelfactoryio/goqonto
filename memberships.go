package goqonto

import (
	"context"
	"net/http"
)

// membershipsBasePath Qonto API Memberships Endpoint
const membershipsBasePath = "v2/memberships"

// MembershipsService provides access to to the memberships in Qonto API
type MembershipsService service

// MembershipsOptions Qonto API Memberships query strings
// https://api-doc.qonto.eu/2.0/memberships/list-memberships
type MembershipsOptions struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
}

// Membership struct
// https://api-doc.qonto.eu/2.0/memberships/list-memberships
type Membership struct {
	ID       string `json:"id"`
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
}

// membershipsRoot root key in the JSON response for memberships
type membershipsRoot struct {
	Memberships []Membership `json:"memberships"`
}

// List all the memberships
func (s *MembershipsService) List(ctx context.Context, opt *MembershipsOptions) ([]Membership, *Response, error) {

	req, err := s.client.NewRequest(ctx, http.MethodGet, membershipsBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	type respWithMeta struct {
		membershipsRoot
		metaRoot
	}

	root := new(respWithMeta)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if m := &root.metaRoot; m != nil {
		resp.Meta = &m.Meta
	}

	return root.Memberships, resp, nil
}
