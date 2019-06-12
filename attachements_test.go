package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestAttachmentsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/1ec373a5-e30d-4a70-948d-c8d49e4a4d31", attachmentsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
			"attachment": {
				"id": "1ec373a5-e30d-4a70-948d-c8d49e4a4d31",
				"created_at": "2019-01-07T16:36:25.862Z",
				"file_name": "doc.pdf",
				"file_size": "49599",
				"file_content_type": "application/pdf",
				"url": "https://mybucket.s3.eu-central-1.amazonaws.com/doc.pdf"
			}
		}`

		_, err := fmt.Fprint(w, response)
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}
	})

	attachment, _, err := client.Attachments.Get(ctx, "1ec373a5-e30d-4a70-948d-c8d49e4a4d31")
	if err != nil {
		t.Errorf("Attachments.Get returned error: %v", err)
	}

	createdAt, _ := time.Parse(time.RFC3339, "2019-01-07T16:36:25.862Z")

	expected := &Attachment{
		ID:              "1ec373a5-e30d-4a70-948d-c8d49e4a4d31",
		CreatedAt:       createdAt,
		FileName:        "doc.pdf",
		FileSize:        "49599",
		FileContentType: "application/pdf",
		URL:             "https://mybucket.s3.eu-central-1.amazonaws.com/doc.pdf",
	}

	if !reflect.DeepEqual(attachment, expected) {
		t.Errorf("Attachments.Get \n returned: %+v\n expected: %+v\n", attachment, expected)
	}

}
