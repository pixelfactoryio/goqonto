package goqonto

import (
	"context"
	"fmt"
	"net/http"
	"time"
)

// transactionsBasePath Qonto API Attachments Endpoint
const attachmentsBasePath = "v2/attachments"

// AttachmentsService provides access to the attachements in Qonto API
type AttachmentsService service

// Attachment struct
// https://api-doc.qonto.eu/2.0/attachments/show-attachment-1
type Attachment struct {
	ID              string    `json:"id"`
	CreatedAt       time.Time `json:"created_at"`
	FileName        string    `json:"file_name"`
	FileSize        string    `json:"file_size"`
	FileContentType string    `json:"file_content_type"`
	URL             string    `json:"url"`
}

// attachmentsRoot root key in the JSON response for attachments
type attachmentsRoot struct {
	Attachment *Attachment `json:"attachment"`
}

// Get Attachment
func (a *AttachmentsService) Get(ctx context.Context, id string) (*Attachment, *Response, error) {

	path := fmt.Sprintf("%s/%s", attachmentsBasePath, id)

	req, err := a.client.NewRequest(ctx, http.MethodGet, path, nil)
	if err != nil {
		return nil, nil, err
	}

	root := new(attachmentsRoot)
	resp, err := a.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	return root.Attachment, resp, nil
}
