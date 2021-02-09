package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/pixelfactoryio/goqonto/v2"
)

// AuthTransport structs holds company Slug and  Secret key
type AuthTransport struct {
	*http.Transport
	Slug   string
	Secret string
}

// RoundTrip set "Authorization" header
func (t AuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("%s:%s", t.Slug, t.Secret))
	return t.Transport.RoundTrip(r)
}

func main() {
	orgID := os.Getenv("QONTO_ORG_ID")
	userLogin := os.Getenv("QONTO_USER_LOGIN")
	userSecretKey := os.Getenv("QONTO_SECRET_KEY")

	client := http.Client{
		Transport: AuthTransport{
			&http.Transport{},
			userLogin,
			userSecretKey,
		},
	}

	qonto := goqonto.NewClient(&client)
	ctx := context.Background()

	// Get Organization
	orga, resp, err := qonto.Organizations.Get(ctx, orgID)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	prettyPrint(orga)

	// List Transactions
	params := &goqonto.TransactionsOptions{
		Slug:          orga.Slug,
		IBAN:          orga.BankAccounts[0].IBAN,
		Side:          goqonto.TransactionSideCredit,
		SortBy:        goqonto.TransactionSortByUpdatedAtAsc,
		OperationType: []string{goqonto.TransactionOperationTypeCard},
		Status:        []string{goqonto.TransactionStatusCompleted},
	}

	transactions, resp, err := qonto.Transactions.List(ctx, params)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}

	for _, trx := range transactions {
		prettyPrint(trx)
	}
	prettyPrint(resp.Meta)

}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(b))
}
