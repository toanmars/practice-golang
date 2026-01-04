package main

type PaymentMethod interface {
	Pay(amount int) error
}

type CreditCard struct {
	Number string
	CVV    string
	Expire string
}

func (c *CreditCard) Pay(amount int) error {
	return nil
}

type Momo struct {
	Number string
}

func (m *Momo) Pay(amount int) error {
	return nil
}

func ProcessPayment(p PaymentMethod) error {
	return p.Pay(100)
}

type BitCoin struct {
	Address string
}

func (b *BitCoin) Pay(amount int) error {
	return nil
}

func main() {
	creditCard := CreditCard{}
	momo := Momo{}
	bitcoin := BitCoin{}

	ProcessPayment(&creditCard)
	ProcessPayment(&momo)
	ProcessPayment(&bitcoin)
}
