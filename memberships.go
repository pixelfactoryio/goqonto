package goqonto

import (
	"context"
	"encoding/json"
)

// membershipsBasePath Qonto API Memberships Endpoint
const membershipsBasePath = "v2/memberships"

// MembershipsOptions Qonto API Memberships query strings
// https://api-doc.qonto.eu/2.0/memberships/list-memberships
type MembershipsOptions struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
}

// MembershipsService interface
// List: list all the memberships
type MembershipsService interface {
	List(context.Context, *MembershipsOptions) ([]Membership, *Response, error)
}

// Membership struct
// https://api-doc.qonto.eu/2.0/memberships/list-memberships
type Membership struct {
	ID       string `json:"id"`
	FistName string `json:"first_name"`
	LastName string `json:"last_name"`
}

// MembershipsServiceOp struct used to embed *Client
type MembershipsServiceOp struct {
	client *Client
}

var _ MembershipsService = &MembershipsServiceOp{}

// membershipsRoot root key in the JSON response for memberships
type membershipsRoot struct {
	Memberships []Membership `json:"memberships"`
}

// Convert Membership to a string
// TODO: shouldn't Panic here
func (m Membership) String() string {
	bytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// List all the memberships
func (m *MembershipsServiceOp) List(ctx context.Context, memOpt *MembershipsOptions) ([]Membership, *Response, error) {

	req, err := m.client.NewRequest(ctx, "GET", membershipsBasePath, memOpt)
	if err != nil {
		return nil, nil, err
	}

	type respWithMeta struct {
		membershipsRoot
		metaRoot
	}

	root := new(respWithMeta)
	resp, err := m.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	if m := &root.metaRoot; m != nil {
		resp.Meta = &m.Meta
	}

	return root.Memberships, resp, nil
}
