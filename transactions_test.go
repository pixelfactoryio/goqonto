package goqonto

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
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
					"transaction_id": "mycompany-bank-account-1-transaction-491",
					"amount": 126.0,
					"amount_cents": 12600,
					"attachment_ids": [],
					"local_amount": 126.0,
					"local_amount_cents": 12600,
					"side": "debit",
					"operation_type": "transfer",
					"currency": "EUR",
					"local_currency": "EUR",
					"label": "Qonto",
					"settled_at": "2019-05-29T05:28:00.191Z",
					"emitted_at": "2019-05-29T05:27:51.353Z",
					"updated_at": "2019-05-29T05:29:38.068Z",
					"status": "completed",
					"note": "Memo added by the user on the transaction",
					"reference": "Message sent along income, transfer and direct_debit transactions",
					"vat_amount": 21.0,
					"vat_amount_cents": 2100,
					"vat_rate": 20.0,
					"initiator_id": "ID of the membership who initiated the transaction",
					"label_ids": [
						"f4c39147-9c1f-43b0-4720-bd6ef6ac372d",
						"dee2cfdb-8147-444c-9b8d-5c0aa25b8dd9"
					],
					"attachment_lost": false,
					"attachment_required": true
				},
				{
					"transaction_id": "mycompany-bank-account-1-transaction-490",
					"amount": 136.8,
					"amount_cents": 13680,
					"attachment_ids": [
					  "b324f133-187c-4684-818d-530110a76521"
					],
					"local_amount": 136.8,
					"local_amount_cents": 13680,
					"side": "debit",
					"operation_type": "income",
					"currency": "EUR",
					"local_currency": "EUR",
					"label": "Qonto",
					"settled_at": "2019-05-28T15:18:12.102Z",
					"emitted_at": "2019-05-28T15:18:04.938Z",
					"updated_at": "2019-05-29T05:29:01.420Z",
					"status": "completed",
					"note": "",
					"reference": "",
					"vat_amount": 22.8,
					"vat_amount_cents": 2280,
					"vat_rate": 20.0,
					"initiator_id": "",
					"label_ids": [],
					"attachment_lost": false,
					"attachment_required": true
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

		_, err := fmt.Fprint(w, response)
		if err != nil {
			t.Errorf("Unable to write response error: %v", err)
		}
	})

	params := &TransactionsOptions{
		Slug: "mycompany-9134",
		IBAN: "FR761679800001000000123456",
	}

	transactions, resp, err := client.Transactions.List(context.Background(), params)
	if err != nil {
		t.Errorf("Transactions.Get returned error: %v", err)
	}

	trx1SettledAt, _ := time.Parse(time.RFC3339, "2019-05-29T05:28:00.191Z")
	trx1EmittedAt, _ := time.Parse(time.RFC3339, "2019-05-29T05:27:51.353Z")
	trx1UpdatedAt, _ := time.Parse(time.RFC3339, "2019-05-29T05:29:38.068Z")
	trx1 := Transaction{
		TransactionID:    "mycompany-bank-account-1-transaction-491",
		Amount:           126.0,
		AmountCents:      12600,
		AttachmentIds:    []string{},
		LocalAmount:      126.0,
		LocalAmountCents: 12600,
		Side:             "debit",
		OperationType:    "transfer",
		Currency:         "EUR",
		LocalCurrency:    "EUR",
		Label:            "Qonto",
		SettledAt:        trx1SettledAt,
		EmittedAt:        trx1EmittedAt,
		UpdatedAt:        trx1UpdatedAt,
		Status:           "completed",
		Note:             "Memo added by the user on the transaction",
		Reference:        "Message sent along income, transfer and direct_debit transactions",
		VatAmount:        21.0,
		VatAmountCents:   2100,
		VatRate:          20.0,
		InitiatorID:      "ID of the membership who initiated the transaction",
		LabelIds: []string{
			"f4c39147-9c1f-43b0-4720-bd6ef6ac372d",
			"dee2cfdb-8147-444c-9b8d-5c0aa25b8dd9",
		},
		AttachmentLost:     false,
		AttachmentRequired: true,
	}

	trx2SettledAt, _ := time.Parse(time.RFC3339, "2019-05-28T15:18:12.102Z")
	trx2EmittedAt, _ := time.Parse(time.RFC3339, "2019-05-28T15:18:04.938Z")
	trx2UpdatedAt, _ := time.Parse(time.RFC3339, "2019-05-29T05:29:01.420Z")
	trx2 := Transaction{
		TransactionID:      "mycompany-bank-account-1-transaction-490",
		Amount:             136.8,
		AmountCents:        13680,
		AttachmentIds:      []string{"b324f133-187c-4684-818d-530110a76521"},
		LocalAmount:        136.8,
		LocalAmountCents:   13680,
		Side:               "debit",
		OperationType:      "income",
		Currency:           "EUR",
		LocalCurrency:      "EUR",
		Label:              "Qonto",
		SettledAt:          trx2SettledAt,
		EmittedAt:          trx2EmittedAt,
		UpdatedAt:          trx2UpdatedAt,
		Status:             "completed",
		Note:               "",
		Reference:          "",
		VatAmount:          22.8,
		VatAmountCents:     2280,
		VatRate:            20.0,
		InitiatorID:        "",
		LabelIds:           []string{},
		AttachmentLost:     false,
		AttachmentRequired: true,
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
		t.Errorf("Transactions.Get \n returned: %+v\n expected: %+v\n", transactions, expectedTrx)
	}

	if !reflect.DeepEqual(resp.Meta, expectedMeta) {
		t.Errorf("Transactions.Get \n returned: %+v\n expected: %+v\n", resp, expectedMeta)
	}
}
