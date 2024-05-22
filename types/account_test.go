package types

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var statementJson = []byte(`
{
	"fullName": "Sandra Saulgrieze", 
	"address":  "14 The Dale\nWhitefield hall\nBettystown\nMeath\nA92N27C",
	"currency": "EUR",
	"bankAccounts": [
		{
			"iban": "IE30REVO99036022547749",
			"bic": "REVOIE23",
			"transferAccount": true
		},
		{
			"iban": "LT093250041069208595",
			"bic": "REVOLT21"
		},
		{
			"iban": "LT087070024346246713",
			"bic": "RELBLT21"
		}
	],
	"startDate": "2023-02-01T00:00:00.000Z",
	"endDate": "2023-03-29T23:59:59.999Z",
	"products": [
		{
			"name": "Account (Current Account)",
			"openingBalance": 2.52,
			"moneyOut": 1944.09,
			"moneyIn": 1978.00
		}
	],
	"transactions": [
		{
			"date": "2023-02-03T09:12:25.000Z",
			"description": "Apple Pay Top-Up by *5453",
			"moneyIn": 50,
			"balance": 52.52
		}
	]
}`)

func TestStatementJson(t *testing.T) {
	var s Statement
	err := json.Unmarshal(statementJson, &s)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, "Sandra Saulgrieze", s.FullName)
	assert.Equal(t, "14 The Dale\nWhitefield hall\nBettystown\nMeath\nA92N27C", s.Address)

	assert.Equal(t, Eur, s.Currency)

	assert.Len(t, s.BankAccounts, 3)
	assert.Equal(t, "IE30REVO99036022547749", s.BankAccounts[0].Iban)
	assert.Equal(t, "REVOIE23", s.BankAccounts[0].Bic)
	assert.True(t, s.BankAccounts[0].ValidTransfer)
	assert.False(t, s.BankAccounts[1].ValidTransfer)

	assert.Equal(t, NewDate(2023, 2, 1, 0, 0, 0, 0, time.UTC), s.StartDate)
	assert.Equal(t, NewDate(2023, 3, 29, 23, 59, 59, 999000000, time.UTC), s.EndDate)

	assert.Len(t, s.Products, 1)
	assert.Equal(t, "Account (Current Account)", s.Products[0].Name)
	assert.Equal(t, 2.52, s.Products[0].OpeningBalance)
	assert.Equal(t, 1944.09, s.Products[0].MoneyOut)
	assert.Equal(t, 1978.0, s.Products[0].MoneyIn)

	assert.Len(t, s.Transactions, 1)
	assert.Equal(t, NewDate(2023, 2, 3, 9, 12, 25, 0, time.UTC), s.Transactions[0].Date)
	assert.Equal(t, "Apple Pay Top-Up by *5453", s.Transactions[0].Description)
	assert.Equal(t, "", s.Transactions[0].DescriptionDetail)
	assert.Equal(t, 0.0, s.Transactions[0].MoneyOut)
	assert.Equal(t, 50.0, s.Transactions[0].MoneyIn)
	assert.Equal(t, 52.52, s.Transactions[0].Balance)

}
