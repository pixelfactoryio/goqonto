package goqonto

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestMembershipsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s", membershipsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
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

		_, err := fmt.Fprint(w, response)
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}
	})

	params := &MembershipsOptions{
		CurrentPage: 1,
		PerPage:     10,
	}

	memberships, resp, err := client.Memberships.List(context.Background(), params)
	if err != nil {
		t.Errorf("Memberships.Get returned error: %v", err)
	}

	member1 := Membership{
		ID:       "6dbdc8ad-2c89-483c-b696-781c86fa1db4",
		FistName: "Bob",
		LastName: "Foo",
	}

	member2 := Membership{
		ID:       "88e4a3e6-5012-4c01-8507-362a88712f77",
		FistName: "Emmett",
		LastName: "Brown",
	}

	expectedMembers := new(membershipsRoot).Memberships
	expectedMembers = append(expectedMembers, member1)
	expectedMembers = append(expectedMembers, member2)

	expectedMeta := &ResponseMeta{
		CurrentPage: 1,
		NextPage:    0,
		PrevPage:    0,
		TotalPages:  1,
		TotalCount:  2,
		PerPage:     10,
	}

	if !reflect.DeepEqual(memberships, expectedMembers) {
		t.Errorf("Memberships.Get \n returned: %+v\n expected: %+v\n", memberships, expectedMembers)
	}

	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Memberships.Get \n returned: %+v\n expected: %+v\n", resp, expectedMeta)
	}
}
