package pdf

import (
	"assignment21-gopdf/types"
	"errors"
	"io/fs"
	"os"
	"testing"
	"time"
)

const fileName = "test.pdf"

var statement = types.Statement{
	FullName: "Sandra Saulgrieze",
	Address:  "14 The Dale\nWhitefield hall\nBettystown\nMeath\nA92N27C",
	Currency: types.Eur,
	BankAccounts: []types.BankAccount{
		{
			Iban:          "IE30REVO99036022547749",
			Bic:           "REVOIE23",
			ValidTransfer: true,
		},
		{
			Iban:          "LT093250041069208595",
			Bic:           "REVOLT21",
			ValidTransfer: false,
		},
		{
			Iban:          "LT087070024346246713",
			Bic:           "RELBLT21",
			ValidTransfer: false,
		},
	},
	StartDate: types.NewDate(2023, 02, 01, 0, 0, 0, 0, time.Local),
	EndDate:   types.NewDate(2023, 03, 29, 23, 59, 59, 999, time.Local),
	Products: []types.Product{
		{
			Name:           "Account (Current Account)",
			OpeningBalance: 2.52,
			MoneyOut:       1944.09,
			MoneyIn:        1978.00,
		},
	},
	Transactions: []types.Transaction{
		{
			Date:        types.NewDate(2023, 02, 03, 9, 12, 25, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     50,
			Balance:     52.52,
		},
		{
			Date:        types.NewDate(2023, 02, 03, 11, 5, 40, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     100,
			Balance:     152.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          100,
			Balance:           52.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          10,
			Balance:           42.52,
		},
		{
			Date:        types.NewDate(2023, 02, 03, 9, 12, 25, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     50,
			Balance:     52.52,
		},
		{
			Date:        types.NewDate(2023, 02, 03, 11, 5, 40, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     100,
			Balance:     152.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          100,
			Balance:           52.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          10,
			Balance:           42.52,
		},
		{
			Date:        types.NewDate(2023, 02, 03, 9, 12, 25, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     50,
			Balance:     52.52,
		},
		{
			Date:        types.NewDate(2023, 02, 03, 11, 5, 40, 0, time.Local),
			Description: "Apple Pay Top-Up by *5453",
			MoneyIn:     100,
			Balance:     152.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          100,
			Balance:           52.52,
		},
		{
			Date:              types.NewDate(2023, 02, 03, 13, 41, 19, 0, time.Local),
			Description:       "To LIMA MILLER SAULGRIEZE",
			DescriptionDetail: "To LIMA MILLER SAULGRIEZE",
			MoneyOut:          10,
			Balance:           42.52,
		},
	},
}

func TestGetStatement(t *testing.T) {
	f, err := os.Open(fileName)
	if errors.Is(err, fs.ErrNotExist) {
		f, err = os.Create(fileName)
	}
	if err != nil {
		t.Error(err)
	}
	_, err = GenerateStatement(statement, f)
	if err != nil {
		t.Error("could not get statement pdf", err)
	}

	err = f.Close()
	if err != nil {
		t.Error(err)
	}
	//
	//err = os.Remove(fileName)
	//if err != nil {
	//	t.Error("failed to delete file", err)
	//}
}
