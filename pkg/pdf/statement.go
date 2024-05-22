package pdf

import (
	"assignment21-gopdf/types"
	"bytes"
	"fmt"
	"github.com/johnfercher/maroto/v2"
	"github.com/johnfercher/maroto/v2/pkg/components/code"
	"github.com/johnfercher/maroto/v2/pkg/components/col"
	"github.com/johnfercher/maroto/v2/pkg/components/row"
	"github.com/johnfercher/maroto/v2/pkg/components/text"
	"github.com/johnfercher/maroto/v2/pkg/config"
	"github.com/johnfercher/maroto/v2/pkg/consts/align"
	"github.com/johnfercher/maroto/v2/pkg/consts/border"
	"github.com/johnfercher/maroto/v2/pkg/consts/fontstyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/linestyle"
	"github.com/johnfercher/maroto/v2/pkg/consts/pagesize"
	"github.com/johnfercher/maroto/v2/pkg/core"
	"github.com/johnfercher/maroto/v2/pkg/props"
	"io"
	"strings"
	"time"
)

// GenerateStatement creates the PDF for the provided statement and writes it to the io.Writer if supplied. A bytes.Reader is returned to be used for possible content serving.
func GenerateStatement(statement types.Statement, writer io.Writer) (*bytes.Reader, error) {
	cfg := config.NewBuilder().
		WithPageSize(pagesize.A4).
		WithMargins(10, 15, 10).
		WithPageNumber("Page {current} of {total}", props.RightBottom).
		Build()

	mrt := maroto.New(cfg)

	err := mrt.RegisterHeader(getPageHeader(statement.Currency))
	if err != nil {
		return nil, err
	}

	err = mrt.RegisterFooter(getPageFooter())
	if err != nil {
		return nil, err
	}

	mrt.AddRows(getAccountInfo(statement))

	mrt.AddRows(getBankAccounts(statement.BankAccounts)...)

	mrt.AddRows(getSummary(statement.Products, statement.Currency)...)

	mrt.AddRows(getTransactions(statement)...)

	document, err := mrt.Generate()
	if err != nil {
		return nil, err
	}
	b := document.GetBytes()
	if writer != nil {
		_, err = writer.Write(b)
	}
	return bytes.NewReader(b), err
}

func getAccountInfo(statement types.Statement) core.Row {
	var lines []core.Component
	lines = append(lines, text.New(strings.ToUpper(statement.FullName), props.Text{Size: 14, Style: fontstyle.Bold}))
	addressLines := getAddress(statement.Address)
	lines = append(lines, addressLines...)
	return row.New(33).Add(
		col.New(12).Add(
			lines...,
		),
	)
}

func getAddress(address string) []core.Component {
	lines := strings.Split(address, "\n")
	texts := make([]core.Component, len(lines))
	for i, l := range lines {
		var top = 4 + float64(i+1)*5
		texts[i] = text.New(l, props.Text{Size: 8, Style: fontstyle.Bold, Top: top})
	}
	return texts
}

func getBankAccounts(accounts []types.BankAccount) []core.Row {
	rows := make([]core.Row, len(accounts))
	for i, account := range accounts {
		var height float64
		var lines = make([]core.Component, 2, 3)
		lines[0] = text.New(account.Iban, props.Text{Size: 7})
		lines[1] = text.New(account.Bic, props.Text{Size: 7, Top: 4})
		if account.ValidTransfer {
			height = 12
		} else {
			height = 18
			lines = append(lines, text.New("(You cannot use this IBAN for bank transfers. Please use the IBAN found in the app.", props.Text{Size: 7, Top: 8}))
		}
		r := row.New(height).Add(
			col.New(7),
			col.New(1).Add(
				text.New("IBAN", props.Text{Size: 7, Style: fontstyle.Bold}),
				text.New("BIC", props.Text{Size: 7, Style: fontstyle.Bold, Top: 4}),
			),
			col.New(4).Add(lines...),
		)
		rows[i] = r
	}

	return rows
}

func getSummary(products []types.Product, currency types.Currency) []core.Row {
	rows := make([]core.Row, 0, len(products)+3)
	rows = append(rows, text.NewRow(12, "Balance summary", props.Text{Size: 14, Style: fontstyle.Bold}))
	total := types.Product{
		Name:           "Total",
		OpeningBalance: 0,
		MoneyOut:       0,
		MoneyIn:        0,
	}

	rows = append(rows, row.New(9).WithStyle(&props.Cell{
		BackgroundColor: nil,
		BorderColor: &props.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
		BorderType:      border.Bottom,
		BorderThickness: 0.4,
		LineStyle:       linestyle.Solid,
	}).Add(
		text.NewCol(5, "Product", props.Text{Style: fontstyle.Bold, Top: 1.5}),
		text.NewCol(2, "Opening balance", props.Text{Style: fontstyle.Bold, Top: 1.5}),
		text.NewCol(2, "Money out", props.Text{Style: fontstyle.Bold, Top: 1.5}),
		text.NewCol(2, "Money in", props.Text{Style: fontstyle.Bold, Top: 1.5}),
		text.NewCol(1, "Closing balance", props.Text{Style: fontstyle.Bold}),
	))

	productFunc := func(product types.Product, totalRow bool) core.Row {
		var style *props.Cell
		if totalRow {
			style = nil
		} else {
			style = &props.Cell{
				BackgroundColor: nil,
				BorderColor: &props.Color{
					Red:   0,
					Green: 0,
					Blue:  0,
				},
				BorderType:      border.Bottom,
				BorderThickness: 0.2,
				LineStyle:       linestyle.Solid,
			}
		}
		return row.New(6).WithStyle(style).Add(
			text.NewCol(5, product.Name, props.Text{Top: 0.5}),
			text.NewCol(2, currency.Format(product.OpeningBalance), props.Text{Top: 0.5}),
			text.NewCol(2, currency.Format(product.MoneyOut), props.Text{Top: 0.5}),
			text.NewCol(2, currency.Format(product.MoneyIn), props.Text{Top: 0.5}),
			text.NewCol(1, currency.Format(product.OpeningBalance+product.MoneyIn-product.MoneyOut), props.Text{Top: 0.5}),
		)
	}

	for _, product := range products {
		total.OpeningBalance += product.OpeningBalance
		total.MoneyOut += product.MoneyOut
		total.MoneyIn += product.MoneyIn
		rows = append(rows, productFunc(product, false))
	}

	rows = append(rows, productFunc(total, true))

	rows = append(rows, text.NewRow(20, "The balance on your statement might differ from the balance shown in your app. The statement valance only reflects completed transactions while the app shows the balance available for use, which accounts for pending transactions.", props.Text{Size: 5}))

	return rows
}

func getTransactions(statement types.Statement) []core.Row {
	var rows []core.Row
	rows = append(rows, text.NewRow(12, fmt.Sprintf("Account transactions from %s to %s", statement.StartDate.Format("January 2, 2006"), statement.EndDate.Format("January 2, 2006")), props.Text{Size: 14, Style: fontstyle.Bold}))
	rows = append(rows, row.New(6).WithStyle(&props.Cell{
		BackgroundColor: nil,
		BorderColor: &props.Color{
			Red:   0,
			Green: 0,
			Blue:  0,
		},
		BorderType:      border.Bottom,
		BorderThickness: 0.4,
		LineStyle:       linestyle.Solid,
	}).Add(
		text.NewCol(2, "DateTime", props.Text{Style: fontstyle.Bold}),
		text.NewCol(5, "Description", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Money out", props.Text{Style: fontstyle.Bold}),
		text.NewCol(2, "Money in", props.Text{Style: fontstyle.Bold}),
		text.NewCol(1, "Balance", props.Text{Style: fontstyle.Bold}),
	))

	transactionFunc := func(tx types.Transaction, lastRow bool) core.Row {
		var style *props.Cell
		if lastRow {
			style = nil
		} else {
			style = &props.Cell{
				BackgroundColor: nil,
				BorderColor: &props.Color{
					Red:   0,
					Green: 0,
					Blue:  0,
				},
				BorderType:      border.Bottom,
				BorderThickness: 0.2,
				LineStyle:       linestyle.Solid,
			}
		}

		var desc = make([]core.Component, 1, 2)
		var height float64 = 6
		desc[0] = text.New(tx.Description)
		if len(tx.DescriptionDetail) > 0 {
			desc = append(desc, text.New(tx.DescriptionDetail, props.Text{Size: 5, Top: 5}))
			height = 8
		}

		moneyIn := ""
		if tx.MoneyIn > 0 {
			moneyIn = statement.Currency.Format(tx.MoneyIn)
		}
		moneyOut := ""
		if tx.MoneyOut > 0 {
			moneyOut = statement.Currency.Format(tx.MoneyOut)
		}

		return row.New(height).WithStyle(style).Add(
			text.NewCol(2, tx.Date.Format("Jan 2, 2006"), props.Text{Top: 0.5}),
			col.New(5).Add(desc...),
			text.NewCol(2, moneyOut, props.Text{Top: 0.5}),
			text.NewCol(2, moneyIn, props.Text{Top: 0.5}),
			text.NewCol(1, statement.Currency.Format(tx.Balance), props.Text{Top: 0.5}),
		)
	}

	for i, transaction := range statement.Transactions {
		rows = append(rows, transactionFunc(transaction, i == len(statement.Transactions)-1))
	}

	return rows
}

func getPageHeader(currency types.Currency) core.Row {
	return row.New(35).Add(
		text.NewCol(5, "ReGolute", props.Text{Size: 24, Style: fontstyle.Bold, Align: align.Left}),
		col.New(7).Add(
			text.New(fmt.Sprintf("%s Statement", currency), props.Text{Size: 22, Style: fontstyle.Bold, Align: align.Right}),
			text.New(fmt.Sprintf("Generated on the %s", time.Now().Format("02 Jan 2006")), props.Text{Size: 10, Top: 10, Align: align.Right}),
			text.New("Revolut Bank UAB", props.Text{Size: 10, Top: 14, Align: align.Right}),
		),
	)
}

func getPageFooter() core.Row {
	return row.New(20).Add(
		code.NewQrCol(1, "https://regolute.net/lost-card", props.Rect{
			Top:     1.2,
			Left:    1,
			Percent: 70,
			Center:  false,
		}),
		col.New(3).Add(
			text.New("Report lost or stolen card", props.Text{
				Size:  8,
				Style: fontstyle.Bold,
			}),
			text.New("+547 9 258 3482", props.Text{
				Size: 7,
				Top:  3.5,
			}),
			text.New("Get help directly in app", props.Text{
				Size:  8,
				Top:   7.5,
				Style: fontstyle.Bold,
			}),
			text.New("Scan the QR code", props.Text{
				Size: 7,
				Top:  11,
			}),
		),
		text.NewCol(8, "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Etiam semper justo enim, ac consectetur lectus dignissim nec. Aliquam rhoncus ac lorem sit amet lobortis. Etiam consectetur mi eu arcu condimentum blandit. Cras ut magna eu tortor consequat tempor nec at dui. Cras sit amet interdum nunc, egestas porta turpis. Ut in libero vulputate, interdum metus a, consequat nulla. Nulla erat magna, sollicitudin sed felis sit amet, malesuada convallis turpis. Mauris viverra blandit mattis. Maecenas venenatis sollicitudin elit, eget ullamcorper nulla sollicitudin in.Donec velit purus, mattis commodo suscipit ut, euismod vitae ante. Cras enim justo, consequat nec ultricies ut, hendrerit non orci. Morbi porttitor efficitur ipsum ut auctor. Curabitur leo libero, aliquam quis molestie vitae, tempus et purus. Fusce sed faucibus mauris. Nulla sagittis sem eget felis dictum, malesuada placerat orci viverra. Proin suscipit euismod mollis. Praesent sed augue elit. In maximus eu est ac ultricies. Maecenas a molestie mi. Mauris ultricies fermentum elit vel convallis.", props.Text{Size: 5}),
	)
}
