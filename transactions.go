package goqonto

import (
	"encoding/json"
	"fmt"

	"github.com/amine7536/goqonto/context"
)

const transactionsBasePath = "v1/transactions"

type TransactionsOptions struct {
	Slug string `json:"slug"`
	IBAN string `json:"iban"`
}

type TransactionsService interface {
	List(context.Context, *TransactionsOptions, *ListOptions) ([]Transaction, *Response, error)
	Get(context.Context, string) (*Transaction, *Response, error)
}

type Transaction struct {
	Amout           float32 `json:"amount"`
	AmoutCents      int     `json:"amount_cents"`
	LocalAmout      float32 `json:"local_amount"`
	LocalAmoutCents int     `json:"local_amount_cents"`
	Side            string  `json:"side"`
	OperationType   string  `json:"operation_type"`
	Currency        string  `json:"currency"`
	LocalCurrency   string  `json:"local_currency"`
	Label           string  `json:"label"`
	SettledAt       string  `json:"settled_at"`
}

type TransactionsServiceOp struct {
	client *Client
}

var _ TransactionsService = &TransactionsServiceOp{}

type transactionsRoot struct {
	Transactions []Transaction `json:"transactions"`
}

// Convert Droplet to a string
func (t Transaction) String() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

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
		ResponseMeta
	}

	root := new(respWithMeta)
	resp, err := s.client.Do(ctx, req, root)
	if err != nil {
		return nil, nil, err
	}

	if m := &root.ResponseMeta; m != nil {
		resp.Meta = m
	}

	return root.Transactions, resp, nil
}

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
