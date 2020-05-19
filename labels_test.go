package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	labelsFixture = `
	{
		"labels": [
			{
				"id": "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
				"name": "compta",
				"parent_id": "88e4a3e6-5012-4c01-8507-362a88712f77"
			},
			{
				"id": "88e4a3e6-5012-4c01-8507-362a88712f77",
				"name": "lunch",
				"parent_id": ""
			}
		],
		"meta": {
			"current_page": 1,
			"next_page": null,
			"prev_page": null,
			"total_pages": 1,
			"total_count": 2,
			"per_page": 10
		}
	}`

	label1 = Label{
		ID:       "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
		Name:     "compta",
		ParentID: "88e4a3e6-5012-4c01-8507-362a88712f77",
	}

	label2 = Label{
		ID:       "88e4a3e6-5012-4c01-8507-362a88712f77",
		Name:     "lunch",
		ParentID: "",
	}

	labels = labelsRoot{
		Labels: []Label{label1, label2},
	}
)

func TestLabel_marshall(t *testing.T) {
	testJSONMarshal(t, &Label{}, "{}")

	want := `{
		"id": "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
		"name": "compta",
		"parent_id": "88e4a3e6-5012-4c01-8507-362a88712f77"
	}`

	testJSONMarshal(t, label1, want)
}

func TestLabelsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s", labelsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		_, err := fmt.Fprint(w, labelsFixture)
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}
	})

	params := &LabelsOptions{
		CurrentPage: 1,
		PerPage:     10,
	}

	got, resp, err := client.Labels.List(ctx, params)
	if err != nil {
		t.Errorf("Labels.List returned error: %v", err)
	}

	want := labels.Labels

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Labels.List \n got %v\n want %v\n", got, want)
	}

	testResponseMeta(t, resp.Meta, &ResponseMeta{
		CurrentPage: 1,
		NextPage:    0,
		PrevPage:    0,
		TotalPages:  1,
		TotalCount:  2,
		PerPage:     10,
	})
}

func TestLabelsService_List_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/foo", labelsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		_, err := fmt.Fprint(w, "")
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}
	})

	got, resp, err := client.Labels.List(ctx, &LabelsOptions{})
	if err.Error() == "" {
		t.Errorf("Expected non-empty ErrorResponse.Error()")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Status")
	}

	if got != nil {
		t.Errorf("Expected empty body")
	}
}
