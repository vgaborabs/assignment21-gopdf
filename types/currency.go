package types

import "fmt"

type Currency string

const (
	Eur Currency = "EUR"
	Usd Currency = "USD"
)

func (c Currency) Symbol() (string, bool) {
	switch c {
	case Eur:
		return "€", true
	case Usd:
		return "$", true
	default:
		return "", false
	}
}

func (c Currency) Format(amount float64) string {
	s, leading := c.Symbol()
	if leading {
		return fmt.Sprintf("%s%.2f", s, amount)
	}
	return fmt.Sprintf("%.2f%s", amount, s)
}
