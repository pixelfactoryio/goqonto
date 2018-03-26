package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestTransactionsGet(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s", transactionsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)

		response := `
		{
			"transactions": 
			[
				{
					"amount":0.01,
					"amount_cents":1,
					"local_amount":0.01,
					"local_amount_cents":1,
					"side":"debit",
					"operation_type":"transfer",
					"currency":"EUR",
					"local_currency":"EUR",
					"label":"Amine",
					"settled_at":"2018-03-23T08:23:18.000Z"
				},
				{
					"amount":0.01,
					"amount_cents":1,
					"local_amount":0.01,
					"local_amount_cents":1,
					"side":"debit",
					"operation_type":"transfer",
					"currency":"EUR",
					"local_currency":"EUR",
					"label":"Amine",
					"settled_at":"2018-03-23T08:23:18.000Z"
				}
			],
			"meta": {
				"current_page": 2,
				"next_page": 3,
				"prev_page": 1,
				"total_pages": 3,
				"total_count": 30,
				"per_page": 10
			}
		}`

		fmt.Fprint(w, response)
	})

	params := &TransactionsOptions{
		Slug: "croissant-9134",
		IBAN: "FR7616798000010000004321396",
	}

	list := &ListOptions{
		Page:    1,
		PerPage: 10,
	}

	transactions, resp, err := client.Transactions.List(ctx, params, list)
	if err != nil {
		t.Errorf("Organizations.Get returned error: %v", err)
	}

	trx := Transaction{
		Amount:          0.01,
		AmountCents:     1,
		LocalAmout:      0.01,
		LocalAmoutCents: 1,
		Side:            "debit",
		OperationType:   "transfer",
		Currency:        "EUR",
		LocalCurrency:   "EUR",
		Label:           "Amine",
		SettledAt:       "2018-03-23T08:23:18.000Z",
	}

	expectedTrx := new(transactionsRoot).Transactions

	for i := 0; i < 2; i++ {
		expectedTrx = append(expectedTrx, trx)
	}

	expectedMeta := &ResponseMeta{
		CurrentPage: 2,
		NextPage:    3,
		PrevPage:    1,
		TotalPages:  3,
		TotalCount:  30,
		PerPage:     10,
	}

	if !reflect.DeepEqual(transactions, expectedTrx) {
		t.Errorf("Organizations.Get returned %+v, expected %+v", transactions, expectedTrx)
	}

	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Organizations.Get returned %+v, expected %+v", resp, expectedMeta)
	}
}
