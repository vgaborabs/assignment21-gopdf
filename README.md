# PDF Invoice statement generator written in Golang


## Pre-requisites 
Go 1.22

## Run
```go run main.go```

### Endpoints
- /statement/pdf - POST - Generate and serve
- /statement/pdf/stream - POST - Generate and stream as response

### Test data
```json
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
}
```