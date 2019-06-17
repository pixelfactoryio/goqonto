[![Travis branch](https://img.shields.io/travis/pixelfactoryio/goqonto/v2.svg?style=flat-square)](https://travis-ci.org/pixelfactoryio/goqonto)
[![Coverage Status](https://coveralls.io/repos/github/pixelfactoryio/goqonto/badge.svg)](https://coveralls.io/github/pixelfactoryio/goqonto)

# GoQonto
Qonto API (v2) Go client

## Installation

The import path for the package is gopkg.in/pixelfactoryio/goqonto.v2

To install it, run:

```bash
go get gopkg.in/pixelfactoryio/goqonto.v2
```

## API documentation

Package Documentation is located at : https://godoc.org/gopkg.in/pixelfactoryio/goqonto.v2

Qonto API v2 documentation is located at : https://api-doc.qonto.eu/2.0/welcome

## Usage

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"gopkg.in/pixelfactoryio/goqonto.v2"
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

	// Get Organisation
	orga, resp, err := qonto.Organizations.Get(ctx, orgID)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	prettyPrint(orga)

	// List Transactions
	params := &goqonto.TransactionsOptions{
		Slug:   orga.Slug,
		IBAN:   orga.BankAccounts[0].IBAN,
		Status: []string{"completed"},
	}

	transactions, resp, err := qonto.Transactions.List(ctx, params)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}

	for _, trx := range transactions {
		prettyPrint(trx)
	}
	prettyPrint(resp.Meta)

	// Get an attachment
	attachement, resp, err := qonto.Attachments.Get(ctx, "1812345c-cf62-49a0-bbb0-f654321678")
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}
	prettyPrint(attachement)

	// List memberships
	memberships, resp, err := qonto.Memberships.List(ctx, nil)
	if err != nil && resp.StatusCode != http.StatusOK {
		panic(err.Error())
	}

	for _, member := range memberships {
		prettyPrint(member)
	}
	prettyPrint(resp.Meta)

}

func prettyPrint(v interface{}) {
	b, _ := json.MarshalIndent(v, "", "  ")
	println(string(b))
}
```

## Credits

This client is heavily inspired by :
- DigitalOcean GoDo : https://github.com/digitalocean/godo
- Google Go-Github : https://github.com/google/go-github
