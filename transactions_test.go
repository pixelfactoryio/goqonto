package goqonto

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	transactionsFixture = `{
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

	transactionFixture = `{
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
	}`

	trx1SettledAt, _ = time.Parse(time.RFC3339, "2019-05-29T05:28:00.191Z")
	trx1EmittedAt, _ = time.Parse(time.RFC3339, "2019-05-29T05:27:51.353Z")
	trx1UpdatedAt, _ = time.Parse(time.RFC3339, "2019-05-29T05:29:38.068Z")
	trx1             = Transaction{
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

	transactions = transactionsRoot{
		Transactions: []Transaction{trx1},
	}
)

func TestTransaction_marshall(t *testing.T) {
	testJSONMarshal(t, &Transaction{}, "{}")
	testJSONMarshal(t, &trx1, transactionFixture)

}

func TestTransactionsService_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc(fmt.Sprintf("/%s", transactionsBasePath), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		testBody(t, r, `{"slug":"mycompany-9134","iban":"FR761679800001000000123456"}`+"\n")
		fmt.Fprint(w, transactionsFixture)
	})

	params := &TransactionsOptions{
		Slug: "mycompany-9134",
		IBAN: "FR761679800001000000123456",
	}

	got, resp, err := client.Transactions.List(ctx, params)

	if err != nil {
		t.Errorf("Transactions.List returned error: %v", err)
	}

	want := transactions.Transactions

	if !reflect.DeepEqual(got, want) {
		t.Errorf("Transactions.List \n got %v\n want %v\n", got, want)
	}

	testResponseMeta(t, resp.Meta, &ResponseMeta{
		CurrentPage: 2,
		NextPage:    3,
		PrevPage:    1,
		TotalPages:  3,
		TotalCount:  30,
		PerPage:     10,
	})
}

func TestTransactionsService_List_Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		testHeader(t, r, "Accept", mediaType)
		testHeader(t, r, "Content-Type", mediaType)
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, `{ "message": "Not found" }`)
	})

	got, resp, err := client.Transactions.List(ctx, &TransactionsOptions{})

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
