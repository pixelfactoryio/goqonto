package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestOrganizationsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s/9134", organizationsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
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

		_, err := fmt.Fprint(w, response)
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}	})

	orga, _, err := client.Organizations.Get(ctx, "9134")
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	bankAccount := BankAccount{
		Slug:                   "croissant-bank-account-1",
		IBAN:                   "FR7616798000010000004321396",
		BIC:                    "TRZOFR21XXX",
		Currency:               "EUR",
		Balance:                225.3,
		BalanceCents:           22530,
		AuthorizedBalance:      213.2,
		AuthorizedBalanceCents: 21320,
	}

	expected := &Organization{
		Slug:         "croissant-9134",
		BankAccounts: []BankAccount{bankAccount},
	}

	if !reflect.DeepEqual(orga, expected) {
		t.Errorf("Organizations.Get returned %+v, expected %+v", orga, expected)
	}

}
