package types

type Statement struct {
	FullName     string        `json:"fullName"`
	Address      string        `json:"address"`
	Currency     Currency      `json:"currency"`
	BankAccounts []BankAccount `json:"bankAccounts"`
	StartDate    DateTime      `json:"startDate"`
	EndDate      DateTime      `json:"endDate"`
	Products     []Product     `json:"products"`
	Transactions []Transaction `json:"transactions"`
}

type BankAccount struct {
	Iban          string `json:"iban"`
	Bic           string `json:"bic"`
	ValidTransfer bool   `json:"transferAccount,omitempty"`
}

type Product struct {
	Name           string  `json:"name"`
	OpeningBalance float64 `json:"openingBalance"`
	MoneyOut       float64 `json:"moneyOut"`
	MoneyIn        float64 `json:"moneyIn"`
}

type Transaction struct {
	Date              DateTime `json:"date"`
	Description       string   `json:"description"`
	DescriptionDetail string   `json:"descriptionDetail,omitempty"`
	MoneyOut          float64  `json:"moneyOut,omitempty"`
	MoneyIn           float64  `json:"moneyIn,omitempty"`
	Balance           float64  `json:"balance"`
}
