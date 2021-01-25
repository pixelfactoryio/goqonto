package goqonto

import (
	"context"
	"net/http"
	"time"
)

// transactionsBasePath Qonto API Transactions Endpoint
const transactionsBasePath = "v2/transactions"

// TransactionsService provides access to the transactions in Qonto API
type TransactionsService service

// TransactionsOptions Qonto API Transactions query strings
// https://api-doc.qonto.eu/2.0/transactions/list-transactions
type TransactionsOptions struct {
	Slug          string   `json:"slug"`
	IBAN          string   `json:"iban"`
	Status        []string `json:"status,omitempty"`
	UpdatedAtFrom string   `json:"updated_at_from,omitempty"`
	UpdatedAtTo   string   `json:"updated_at_to,omitempty"`
	SettledAtFrom string   `json:"settled_at_from,omitempty"`
	SettledAtTo   string   `json:"settled_at_to,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	CurrentPage   int64    `json:"current_page,omitempty"`
	PerPage       int64    `json:"per_page,omitempty"`
}

// Transaction struct
// https://api-doc.qonto.eu/2.0/transactions/list-transactions
type Transaction struct {
	TransactionID      string    `json:"transaction_id"`
	Amount             float64   `json:"amount"`
	AmountCents        int       `json:"amount_cents"`
	AttachmentIds      []string  `json:"attachment_ids,omitempty"`
	LocalAmount        float64   `json:"local_amount"`
	LocalAmountCents   int       `json:"local_amount_cents"`
	Side               string    `json:"side"`
	OperationType      string    `json:"operation_type"`
	Currency           string    `json:"currency"`
	LocalCurrency      string    `json:"local_currency"`
	Label              string    `json:"label"`
	SettledAt          time.Time `json:"settled_at,omitempty"`
	EmittedAt          time.Time `json:"emitted_at"`
	UpdatedAt          time.Time `json:"updated_at"`
	Status             string    `json:"status"`
	Note               string    `json:"note,omitempty"`
	Reference          string    `json:"reference,omitempty"`
	VatAmount          float64   `json:"vat_amount,omitempty"`
	VatAmountCents     int       `json:"vat_amount_cents,omitempty"`
	VatRate            float64   `json:"vat_rate,omitempty"`
	InitiatorID        string    `json:"initiator_id,omitempty"`
	LabelIds           []string  `json:"label_ids,omitempty"`
	AttachmentLost     bool      `json:"attachment_lost,omitempty"`
	AttachmentRequired bool      `json:"attachment_required,omitempty"`
}

// transactionsRoot root key in the JSON response for transactions
type transactionsRoot struct {
	Transactions []Transaction `json:"transactions"`
}

// List all the transactions for a given Org.Slug and BankAccount.IBAN
func (s *TransactionsService) List(ctx context.Context, opt *TransactionsOptions) ([]Transaction, *Response, error) {

	req, err := s.client.NewRequest(ctx, http.MethodGet, transactionsBasePath, opt)
	if err != nil {
		return nil, nil, err
	}

	type respWithMeta struct {
		transactionsRoot
		metaRoot
	}

	root := new(respWithMeta)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, resp, err
	}

	if m := &root.metaRoot; m != nil {
		resp.Meta = &m.Meta
	}

	return root.Transactions, resp, nil
}
