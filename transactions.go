package goqonto

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// TransactionsService provides access to the transactions in Qonto API
type TransactionsService service

// transactionsBasePath Qonto API Transactions Endpoint
const transactionsBasePath = "v2/transactions"

// TransactionStatusPending is a transaction that is processing and has impacted
// the bank account's auth_balance but not its balance.
const TransactionStatusPending = "pending"

// TransactionStatusReversed is a transaction that used to be processing, but has then been reversed.
const TransactionStatusReversed = "reversed"

// TransactionStatusDeclined is a transaction that has been declined.
const TransactionStatusDeclined = "declined"

// TransactionStatusCompleted is a transaction that is completed, and has impacted the bank account's balance.
const TransactionStatusCompleted = "completed"

// TransactionSideDebit is a debit transaction.
const TransactionSideDebit = "debit"

// TransactionSideCredit is a credit transaction.
const TransactionSideCredit = "credit"

// TransactionOperationTypeCard is a card transaction.
const TransactionOperationTypeCard = "card"

// TransactionOperationTypeTransfer is a transfer transaction.
const TransactionOperationTypeTransfer = "transfer"

// TransactionOperationTypeIncome is an income transaction.
const TransactionOperationTypeIncome = "income"

// TransactionSortByUpdatedAtDesc sort transactions by descending updated_at field.
const TransactionSortByUpdatedAtDesc = "updated_at:desc"

// TransactionSortByUpdatedAtAsc sort transactions by ascending updated_at field.
const TransactionSortByUpdatedAtAsc = "updated_at:asc"

// TransactionSortBySettledAtDesc sort transactions by descending settled_at field.
const TransactionSortBySettledAtDesc = "settled_at:desc"

// TransactionSortBySettledAtAsc sort transactions by ascending settled_at field.
const TransactionSortBySettledAtAsc = "settled_at:asc"

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
	Side          string   `json:"side,omitempty"`
	OperationType []string `json:"operation_type,omitempty"`
	SortBy        string   `json:"sort_by,omitempty"`
	CurrentPage   int64    `json:"current_page,omitempty"`
	PerPage       int64    `json:"per_page,omitempty"`
}

// Transaction struct
// https://api-doc.qonto.eu/2.0/transactions/list-transactions
type Transaction struct {
	// Transaction ID.
	TransactionID string `json:"transaction_id"`

	// Transaction amount in euros.
	Amount float64 `json:"amount"`

	// Transaction amount in euro cents.
	AmountCents int `json:"amount_cents"`

	// Slice of Attachment ids.
	AttachmentIds []string `json:"attachment_ids,omitempty"`

	// Transaction local currency in ISO 4217 currency code.
	LocalCurrency string `json:"local_currency"`

	// Transaction amount in local currency.
	LocalAmount float64 `json:"local_amount"`

	// Transaction amount in local currency cents.
	LocalAmountCents int `json:"local_amount_cents"`

	// Transaction side.
	// Allowed values: debit, credit.
	Side string `json:"side"`

	// Transaction operation type
	// Allowed values: transfer, card, direct_debit, income, qonto_fee, cheque, recall, swift_income.
	OperationType string `json:"operation_type"`

	// Transaction currency in ISO 4217 currency (can only be EUR, currently).
	Currency string `json:"currency"`

	// Transaction counterparty label.
	Label string `json:"label"`

	// Date the transaction impacted the balance of the account in ISO8601 (yyyy-MM-dd'T'HH:mm:ss.SSSZ).
	SettledAt time.Time `json:"settled_at,omitempty"`

	// Date at which the transaction impacted the authorized balance of the account
	// in ISO8601 (yyyy-MM-dd'T'HH:mm:ss.SSSZ).
	EmittedAt time.Time `json:"emitted_at"`

	// Date at which the transaction was last updated in ISO8601 (yyyy-MM-dd'T'HH:mm:ss.SSSZ).
	UpdatedAt time.Time `json:"updated_at"`

	// Transaction status
	// Allowed values: pending, reversed, declined, completed.
	Status string `json:"status"`

	// Note added by the user on the transaction.
	Note string `json:"note,omitempty"`

	// Message sent along income, transfer, direct_debit and swift_income transactions.
	Reference string `json:"reference,omitempty"`

	// Amount of VAT filled in on the transaction in euros.
	VatAmount float64 `json:"vat_amount,omitempty"`

	// Amount of VAT filled in on the transaction in euro cents.
	VatAmountCents int `json:"vat_amount_cents,omitempty"`

	// Rate of VAT.
	// Allowed values: -1, 0, 2.1, 5.5, 10, 20
	VatRate float64 `json:"vat_rate,omitempty"`

	// ID of the membership who initiated the transaction.
	InitiatorID string `json:"initiator_id,omitempty"`

	// Slive of Labels' id.
	LabelIds []string `json:"label_ids,omitempty"`

	// Transaction labels.
	Labels []Label `json:"labels"`

	// Cards PAN last 4 digits (null if not a card transaction).
	CardLastDigits string `json:"card_last_digits,omitempty"`

	// Transaction category.
	Category string `json:"category,omitempty"`

	// Indicates if the transaction's attachment was lost (default: false).
	AttachmentLost bool `json:"attachment_lost,omitempty"`

	// Indicates if the transaction's attachment is required (default: true).
	AttachmentRequired bool `json:"attachment_required,omitempty"`

	// Transaction attachments.
	Attachments []Attachment `json:"attachements"`
}

// MarshalJSON custom marshaler to handle null Json arrays
func (t Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction

	a := struct {
		Alias
	}{
		Alias: (Alias)(t),
	}

	if a.Attachments == nil {
		a.Attachments = make([]Attachment, 0)
	}

	if a.Labels == nil {
		a.Labels = make([]Label, 0)
	}

	return json.Marshal(a)
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
