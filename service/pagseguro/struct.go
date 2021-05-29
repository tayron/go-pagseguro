package pagseguro

import "encoding/xml"

type PaymentRequest struct {
	XMLName         xml.Name       `xml:"checkout"`
	Email           string         `xml:"email"`
	Token           string         `xml:"token"`
	Currency        string         `xml:"currency"`
	Items           []*PaymentItem `xml:"items>item"`
	ReferenceID     string         `xml:"reference"`
	Buyer           *Buyer         `xml:"sender"`
	Shipping        *Shipping      `xml:"shipping"`
	ExtraAmount     string         `xml:"extraAmount,omitempty"` // use this for discounts or taxes
	RedirectURL     string         `xml:"redirectURL,omitempty"`
	NotificationURL string         `xml:"notificationURL,omitempty"`
	MaxUses         string         `xml:"maxUses,omitempty"`  // from 0 to 999 (the amount of tries a user can do with the same reference ID)
	MaxAge          string         `xml:"maxAge,omitempty"`   // time (in seconds) that the returned payment code is valid (30-999999999)
	Metadata        []*Metadata    `xml:"metadata,omitempty"` // https://pagseguro.uol.com.br/v2/guia-de-integracao/api-de-pagamentos.html#v2-item-api-de-pagamentos-parametros-http
	IsSandbox       bool           `xml:"-"`                  // o PagSeguro não tem um modo sandbox no momento (╯°□°）╯︵ ┻━┻
}

type PaymentItem struct {
	XMLName      xml.Name `xml:"item"`
	Id           string   `xml:"id"`
	Description  string   `xml:"description"`
	PriceAmount  string   `xml:"amount"`
	Quantity     string   `xml:"quantity"`
	ShippingCost string   `xml:"shippingCost,omitempty"`
	Weight       string   `xml:"weight,omitempty"`
}

type Buyer struct {
	Email     string           `xml:"email"`
	Name      string           `xml:"name"`
	Phone     *Phone           `xml:"phone,omitempty"`
	Documents []*BuyerDocument `xml:"documents>document,omitempty"`
	BornDate  string           `xml:"bornDate,omitempty"` //dd/MM/yyyy optional
}

type Phone struct {
	AreaCode    string `xml:"areaCode,omitempty"` // optional
	PhoneNumber string `xml:"number,omitempty"`   // optional
}

type BuyerDocument struct {
	Type  string `xml:"type"` // It's always "CPF" ¯\_(ツ)_/¯
	Value string `xml:"value"`
}

type Shipping struct {
	Type    string           `xml:"type"`
	Cost    string           `xml:"cost"`
	Address *ShippingAddress `xml:"address,omitempty"`
}

type ShippingAddress struct {
	Country    string `xml:"country"`              // It's always "BRA" ¯\_(ツ)_/¯
	State      string `xml:"state,omitempty"`      // "SP"
	City       string `xml:"city,omitempty"`       // max 60 min 2
	PostalCode string `xml:"postalCode,omitempty"` // XXXXXXXX
	District   string `xml:"district,omitempty"`   // Bairro | max chars: 60
	Street     string `xml:"street,omitempty"`     // max: 80
	Number     string `xml:"number,omitempty"`     // max: 20
	Complement string `xml:"complement,omitempty"` // max: 40
}

type Metadata struct {
	Key   string     `xml:"key"`
	Value string     `xml:"value,omitempty"`
	Group []Metadata `xml:"group,omitempty"`
}

type ErrorResponse struct {
	Errors []XMLError `xml:"errors"`
}

type XMLError struct {
	XMLName xml.Name `xml:"error"`
	Code    int      `xml:"code"`
	Message string   `xml:"message"`
}

type PaymentPreResponse struct {
	XMLName xml.Name `xml:"checkout"`
	Code    string   `xml:"code"`
	Data    string   `xml:"data"`
}

type PaymentPreSubmitResult struct {
	XML              string
	CheckoutResponse *PaymentPreResponse
	Error            *ErrorResponse
	Success          bool
}

type Transaction struct {
	XMLName            xml.Name                  `xml:"transaction"`
	Date               string                    `xml:"date,omitempty"`
	Code               string                    `xml:"code,omitempty"`
	Reference          string                    `xml:"reference,omitempty"`
	Type               int                       `xml:"type,omitempty"`
	Status             int                       `xml:"status,omitempty"`
	LastEventDate      string                    `xml:"lastEventDate,omitempty"`
	PaymentMethod      *TransactionPaymentMethod `xml:"paymentMethod,omitempty"`
	GrossAmount        string                    `xml:"grossAmount,omitempty"`
	DiscountAmount     string                    `xml:"discountAmount,omitempty"`
	FeeAmount          string                    `xml:"feeAmount,omitempty"`
	NetAmount          string                    `xml:"netAmount,omitempty"`
	EscrowEndDate      string                    `xml:"escrowEndDate,omitempty"`
	ExtraAmount        string                    `xml:"extraAmount,omitempty"`
	InstallmentCount   int                       `xml:"installmentCount,omitempty"`
	ItemCount          int                       `xml:"itemCount,omitempty"`
	Items              []*PaymentItem            `xml:"items>item,omitempty"`
	Buyer              *Buyer                    `xml:"sender,omitempty"`
	Shipping           *Shipping                 `xml:"shipping,omitempty"`
	CancellationSource string                    `xml:"cancellationSource,omitempty"`
	XML                string
}

type TransactionPaymentMethod struct {
	Type int `xml:"type"`
	Code int `xml:"code"`
}
