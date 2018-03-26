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
					"transaction_id": "croissant-bank-account-1-transaction-491",
					"amount": 0.01,
					"amount_cents": 1,
					"local_amount": 0.01,
					"local_amount_cents": 1,
					"side": "debit",
					"operation_type": "transfer",
					"currency": "EUR",
					"local_currency": "EUR",
					"label": "Amine",
					"settled_at": "2018-03-23T08:23:18.000Z",
					"emitted_at": "2018-03-22T14:47:07.909Z",
					"status": "completed",
					"note": null
				},
				{
					"transaction_id": "croissant-bank-account-1-transaction-490",
					"amount": 0.01,
					"amount_cents": 1,
					"local_amount": 0.01,
					"local_amount_cents": 1,
					"side": "credit",
					"operation_type": "card",
					"currency": "EUR",
					"local_currency": "EUR",
					"label": "Amine",
					"settled_at": "2018-03-22T08:23:43.000Z",
					"emitted_at": "2018-03-22T14:47:07.909Z",
					"status": "completed",
					"note": null
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

	trx1 := Transaction{
		TransactionID:   "croissant-bank-account-1-transaction-491",
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
		EmittedAt:       "2018-03-22T14:47:07.909Z",
		Status:          "completed",
		Note:            "",
	}

	trx2 := Transaction{
		TransactionID:   "croissant-bank-account-1-transaction-490",
		Amount:          0.01,
		AmountCents:     1,
		LocalAmout:      0.01,
		LocalAmoutCents: 1,
		Side:            "credit",
		OperationType:   "card",
		Currency:        "EUR",
		LocalCurrency:   "EUR",
		Label:           "Amine",
		SettledAt:       "2018-03-22T08:23:43.000Z",
		EmittedAt:       "2018-03-22T14:47:07.909Z",
		Status:          "completed",
		Note:            "",
	}

	expectedTrx := new(transactionsRoot).Transactions
	expectedTrx = append(expectedTrx, trx1)
	expectedTrx = append(expectedTrx, trx2)

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
