package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

var (
	organizationFixture = `{
		"organization": {
			"slug": "croissant-9134",
			"bank_accounts": [
				{
					"slug": "croissant-bank-account-1",
					"iban": "FR7616798000010000004321396",
					"bic": "TRZOFR21XXX",
					"currency": "EUR",
					"balance": 225.3,
					"balance_cents": 22530,
					"authorized_balance": 213.2,
					"authorized_balance_cents": 21320
				}
			]
		}
	}`

	bankAccount = BankAccount{
		Slug:                   "croissant-bank-account-1",
		IBAN:                   "FR7616798000010000004321396",
		BIC:                    "TRZOFR21XXX",
		Currency:               "EUR",
		Balance:                225.3,
		BalanceCents:           22530,
		AuthorizedBalance:      213.2,
		AuthorizedBalanceCents: 21320,
	}

	organization = Organization{
		Slug:         "croissant-9134",
		BankAccounts: []BankAccount{bankAccount},
	}
)

func TestOrganization_marshall(t *testing.T) {
	testJSONMarshal(t, Organization{}, "{}")

	want := `{
		"slug": "croissant-9134",
		"bank_accounts": [
			{
				"slug": "croissant-bank-account-1",
				"iban": "FR7616798000010000004321396",
				"bic": "TRZOFR21XXX",
				"currency": "EUR",
				"balance": 225.3,
				"balance_cents": 22530,
				"authorized_balance": 213.2,
				"authorized_balance_cents": 21320
			}
		]
	}`

	testJSONMarshal(t, organization, want)
}

func TestOrganizationsService_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/9134", organizationsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		fmt.Fprint(w, organizationFixture)
	})

	got, _, err := client.Organizations.Get(ctx, "9134")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	want := &organization

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Organizations.Get \n got %v\n want %v\n", got, want)
	}

}

func TestOrganizationsService_Get_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{ "message": "Not found" }`)
	})

	got, resp, err := client.Organizations.Get(ctx, "9134")

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
