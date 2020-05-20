package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	membershipsFixture = `{
		"memberships": [
			{
				"id": "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
				"first_name": "Bob",
				"last_name": "Foo"
			},
			{
				"id": "88e4a3e6-5012-4c01-8507-362a88712f77",
				"first_name": "Emmett",
				"last_name": "Brown"
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

	member1 = Membership{
		ID:       "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
		FistName: "Bob",
		LastName: "Foo",
	}

	member2 = Membership{
		ID:       "88e4a3e6-5012-4c01-8507-362a88712f77",
		FistName: "Emmett",
		LastName: "Brown",
	}

	memberships = membershipsRoot{
		Memberships: []Membership{member1, member2},
	}
)

func TestMembership_marshall(t *testing.T) {
	testJSONMarshal(t, &Membership{}, "{}")

	want := `{
		"id": "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
		"first_name": "Bob",
		"last_name": "Foo"
	}`

	testJSONMarshal(t, member1, want)
}

func TestMembershipsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s", membershipsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		testBody(t, r, `{"current_page":1,"per_page":10}`+"\n")
		fmt.Fprint(w, membershipsFixture)
	})

	params := &MembershipsOptions{
		CurrentPage: 1,
		PerPage:     10,
	}

	got, resp, err := client.Memberships.List(ctx, params)
	if err != nil {
		t.Errorf("Memberships.List returned error: %v", err)
	}

	want := memberships.Memberships

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Memberships.List \n got %v\n want %v\n", got, want)
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

func TestMembershipsService_List_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{ "message": "Not found" }`)
	})

	got, resp, err := client.Memberships.List(ctx, &MembershipsOptions{})

	if err.Error() == "" {
		t.Errorf("Expected non-empty err.Error()")
	}

	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Expected 404 Status")
	}

	if got != nil {
		t.Errorf("Expected empty body")
	}
}
