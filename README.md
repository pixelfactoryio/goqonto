# GoQonto
Qonto API Go client

https://api-doc.qonto.eu/1.0/welcome

(Heavily inspired by DigitalOcean GoDo : https://github.com/digitalocean/godo)

!! Work In Progress !!

## Usage

```go
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/amine7536/goqonto"
)

type AuthTransport struct {
	*http.Transport
	Slug   string
	Secret string
}

func (t AuthTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", fmt.Sprintf("%s:%s", t.Slug, t.Secret))
	return t.Transport.RoundTrip(r)
}

func main() {

	apiURL := os.Getenv("QONTO_API")
	orgID := os.Getenv("QONTO_ORD_ID")
	orgSlug := os.Getenv("QONTO_ORG_SLUG")
	orgSecret := os.Getenv("QONTO_ORG_SECRET")

	client := http.Client{
		Transport: AuthTransport{
			&http.Transport{},
			orgSlug,
			orgSecret,
		},
	}

	qonto := goqonto.New(&client, apiURL)

	ctx := context.Background()

	orga, _, err := qonto.Organizations.Get(ctx, orgID)
	if err != nil {
		log.Fatalln(err.Error())
	}
	fmt.Println(orga)

	params := goqonto.TransactionsOptions{
		Slug: orga.Slug,
		IBAN: orga.BankAccounts[0].IBAN,
	}

	list := goqonto.ListOptions{
		Page:    1,
		PerPage: 10,
	}

	transactions, resp, err := qonto.Transactions.List(ctx, &params, &list)
	if err != nil {
		log.Fatalln(err.Error())
	}

	for _, trx := range transactions {
		fmt.Println(trx)
	}

	fmt.Println(resp.Meta)
}
```