package goqonto

import (
	"context"
	"net/http"
)

// labelsBasePath Qonto API Labels Endpoint
const labelsBasePath = "v2/labels"

// LabelsService provides access to the labels in Qonto API
type LabelsService service

// LabelsOptions Qonto API Labels query strings
// https://api-doc.qonto.eu/2.0/labels/list-labels
type LabelsOptions struct {
	CurrentPage int64 `json:"current_page,omitempty"`
	PerPage     int64 `json:"per_page,omitempty"`
}

// Label struct
// https://api-doc.qonto.eu/2.0/labels/list-labels
type Label struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	ParentID string `json:"parent_id"`
}

// labelsRoot root key in the JSON response for labels
type labelsRoot struct {
	Labels []Label `json:"labels"`
}

// List all the labels
func (s *LabelsService) List(ctx context.Context, opt *LabelsOptions) ([]Label, *Response, error) {

	req, err := s.client.NewRequest(ctx, http.MethodGet, labelsBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	type respWithMeta struct {
		labelsRoot
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

	return root.Labels, resp, nil
}
