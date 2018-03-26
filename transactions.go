package goqonto

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/amine7536/goqonto/context"
)

// transactionsBasePath Qonto API Transactions Endpoint
const transactionsBasePath = "v2/transactions"

// TransactionsOptions Qonto API Transactions query strings
// https://api-doc.qonto.eu/2.0/transactions/list-transactions
type TransactionsOptions struct {
	Slug   string   `json:"slug"`
	IBAN   string   `json:"iban"`
	Status []string `json:"status"`
}

// TransactionsService interface
// List: list all the transactions
// Get: get one transaction by id
type TransactionsService interface {
	List(context.Context, *TransactionsOptions, *ListOptions) ([]Transaction, *Response, error)
	Get(context.Context, string) (*Transaction, *Response, error)
}

// Transaction struct
// https://api-doc.qonto.eu/2.0/transactions/list-transactions
type Transaction struct {
	TransactionID   string    `json:"transaction_id"`
	Amount          float64   `json:"amount"`
	AmountCents     int64     `json:"amount_cents"`
	LocalAmout      float64   `json:"local_amount"`
	LocalAmoutCents int64     `json:"local_amount_cents"`
	Side            string    `json:"side"`
	OperationType   string    `json:"operation_type"`
	Currency        string    `json:"currency"`
	LocalCurrency   string    `json:"local_currency"`
	SettledAt       time.Time `json:"settled_at"`
	EmittedAt       time.Time `json:"emitted_at"`
	Status          string    `jsont:"status"`
	Note            string    `json:"note"`
	Label           string    `json:"label"`
}

// TransactionsServiceOp struct used to embed *Client
type TransactionsServiceOp struct {
	client *Client
}

var _ TransactionsService = &TransactionsServiceOp{}

// transactionsRoot root key in the JSON response for transactions
type transactionsRoot struct {
	Transactions []Transaction `json:"transactions"`
}

// metaRoot root key in the JSON response for meta
type metaRoot struct {
	Meta ResponseMeta `json:"meta"`
}

// Convert Transaction to a string
// TODO: shouldn't Panic here
func (t Transaction) String() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

// List all the transactions for a given Org.Slug and BankAccount.IBAN
func (s *TransactionsServiceOp) List(ctx context.Context, trxOpt *TransactionsOptions, listOpt *ListOptions) ([]Transaction, *Response, error) {

	opt := struct {
		*TransactionsOptions
		*ListOptions
	}{trxOpt, listOpt}

	req, err := s.client.NewRequest(ctx, "GET", transactionsBasePath, opt)
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
		return nil, nil, err
	}

	if m := &root.metaRoot; m != nil {
		resp.Meta = &m.Meta
	}

	return root.Transactions, resp, nil
}

// Get a transaction by its id
func (s *TransactionsServiceOp) Get(ctx context.Context, id string) (*Transaction, *Response, error) {

	path := fmt.Sprintf("%s/%s", transactionsBasePath, id)

	req, err := s.client.NewRequest(ctx, "GET", path, nil)
	if err != nil {
		return nil, nil, err
	}

	v := new(Transaction)
	resp, err := s.client.Do(ctx, req, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}
