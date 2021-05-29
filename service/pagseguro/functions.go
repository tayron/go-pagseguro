package pagseguro

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/paulrosania/go-charset/charset"
	_ "github.com/paulrosania/go-charset/data"
)

func NewPaymentRequest(sellerToken, sellerEmail, referenceID, redirectURL, notificationURL string) *PaymentRequest {
	req := &PaymentRequest{
		Email:           sellerEmail,
		Token:           sellerToken,
		Currency:        "BRL",
		ReferenceID:     referenceID,
		RedirectURL:     redirectURL,
		NotificationURL: notificationURL,
		MaxUses:         "10",
		MaxAge:          "7200",
	}
	return req
}

func (r *PaymentRequest) AddItem(id string, description string, amount float64, quantity int) *PaymentItem {
	item := &PaymentItem{
		Id:          id,
		Description: description,
		PriceAmount: toPriceAmountStr(amount),
		Quantity:    strconv.Itoa(quantity),
	}
	if r.Items == nil {
		r.Items = make([]*PaymentItem, 0)
	}
	r.Items = append(r.Items, item)

	return item
}

func (r *PaymentItem) SetWeight(grams int) *PaymentItem {
	r.Weight = strconv.Itoa(grams)
	return r
}

func (r *PaymentItem) SetAmount(amount float64) *PaymentItem {
	r.PriceAmount = toPriceAmountStr(amount)
	return r
}

func (r *PaymentItem) SetQuantity(quantity int) *PaymentItem {
	r.Quantity = strconv.Itoa(quantity)
	return r
}

func (r *PaymentItem) SetShippingCost(cost float64) *PaymentItem {
	r.ShippingCost = toPriceAmountStr(cost)
	return r
}

func (r *PaymentRequest) SetBuyer(name, email string) *Buyer {
	buyer := &Buyer{
		Name:  name,
		Email: email,
	}
	r.Buyer = buyer
	return buyer
}

func (r *Buyer) SetPhone(areaCode string, phone string) *Buyer {
	r.Phone = &Phone{
		AreaCode:    areaCode,
		PhoneNumber: phone,
	}
	return r
}

func (r *Buyer) SetCPF(cpf string) *Buyer {
	if r.Documents == nil {
		r.Documents = make([]*BuyerDocument, 0)
		r.Documents = append(r.Documents, &BuyerDocument{Type: "CPF"})
	}
	for i := 0; i < len(r.Documents); i++ {
		if r.Documents[i].Type == "CPF" {
			r.Documents[i].Value = cpf
			break
		}
	}
	return r
}

func (r *PaymentRequest) SetShipping(shippingType int, shippingCost float64) *Shipping {
	shipping := &Shipping{
		Type: strconv.Itoa(shippingType),
		Cost: toPriceAmountStr(shippingCost),
	}
	r.Shipping = shipping
	return shipping
}

func (r *Shipping) SetAddress(state, city, postalCode, district, street, number, complement string) *Shipping {
	addr := &ShippingAddress{
		Country:    "BRA",
		State:      state,
		City:       city,
		PostalCode: postalCode,
		District:   district,
		Street:     street,
		Number:     number,
		Complement: complement,
	}
	r.Address = addr
	return r
}

func (r *Shipping) SetAddressStateCity(state, city string) *Shipping {
	if r.Address == nil {
		r.SetAddress(state, city, "", "", "", "", "")
		return r
	}
	r.Address.State = state
	r.Address.City = city
	return r
}

func (r *Shipping) SetAddressCountry(country string) *Shipping {
	if r.Address == nil {
		r.SetAddress("", "", "", "", "", "", "")
	}
	r.Address.Country = country
	return r
}

func (r *PaymentRequest) Submit() (result *PaymentPreSubmitResult) {
	result = &PaymentPreSubmitResult{}

	// Conectar com timeout caso o PagSeguro esteja morgando
	functimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(30*time.Second))
	}

	// create a custom http client that ignores https cert validity, so we don't have to install PagSeguro CAs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            functimeout,
	}
	client := &http.Client{Transport: tr}

	// generate xml
	xmlb, err := xml.Marshal(r)

	if err != nil {
		log.Println("XML MARSHAL ERROR: " + err.Error())
		return
	}

	var clBuffer bytes.Buffer
	clBuffer.WriteString(XMLHeader)
	clBuffer.Write(xmlb)

	urlCheckout := CheckoutURL

	if r.IsSandbox == true {
		urlCheckout = CheckoutSandBoxURL
	}

	checkoutURL := fmt.Sprintf("%s?email=%s&token=%s&charset=%s", urlCheckout, r.Email, r.Token, "UTF-8")

	// send the request (this goroutine is blocked until a response is received)
	resp, err := client.Post(checkoutURL, "application/xml", &clBuffer)

	if err != nil {
		log.Println("Client.Post ERROR: " + err.Error())
		return
	}

	defer resp.Body.Close()
	clBuffer.Truncate(0)

	// io.Copy has a 32kB max buffer size, so no extra memory is consumed
	io.Copy(&clBuffer, resp.Body)
	respBytes := clBuffer.Bytes()
	result.XML = string(respBytes)
	var decoder *xml.Decoder

	errors := &ErrorResponse{}

	clBuffer.Truncate(0)
	clBuffer.Write(respBytes)
	decoder = xml.NewDecoder(&clBuffer)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(errors)

	if err != nil {
		// an error was not found!
		//log.Println("^~PAGSEGO~^ Unmarshal(errors)  ERROR: " + err.Error())
		//return
	} else {
		if errors.Errors != nil {
			if len(errors.Errors) > 0 {
				//log.Println("LOL ERRORS")
				//log.Println(errors.Errors[0].Message)
				result.Error = errors
				result.Success = false
				return
			}
		}
	}

	success := &PaymentPreResponse{}

	clBuffer.Truncate(0)
	clBuffer.Write(respBytes)
	decoder = xml.NewDecoder(&clBuffer)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(success)

	if err != nil {
		log.Println("Unmarshal(success)  ERROR: " + err.Error())
		result.Success = false
		return
	}

	result.CheckoutResponse = success
	result.Success = true
	return
}

func FetchTransactionInfo(sellerToken, sellerEmail, notificationCode string, sandbox bool) (result *Transaction, err error) {
	result = &Transaction{}

	// Conectar com timeout caso o PagSeguro esteja morgando
	functimeout := func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, time.Duration(30*time.Second))
	}

	// create a custom http client that ignores https cert validity, so we don't have to install PagSeguro CAs
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial:            functimeout,
	}
	client := &http.Client{Transport: tr}

	urlTransaction := TransactionsURL

	if sandbox == true {
		urlTransaction = TransactionsSandBoxURL
	}

	transactionsURL := fmt.Sprintf("%s/%s?email=%s&token=%s", urlTransaction, notificationCode, sellerEmail, sellerToken)

	//&charset=%s  , "UTF-8"
	resp, err := client.Get(transactionsURL)

	if err != nil {
		log.Println("Client.Get ERROR: " + err.Error())
		return
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	io.Copy(&buffer, resp.Body)

	respBytes := buffer.Bytes()
	result.XML = string(respBytes)

	decoder := xml.NewDecoder(&buffer)
	decoder.CharsetReader = charset.NewReader
	err = decoder.Decode(result)
	if err != nil {
		log.Println("Decoder.Decode ERROR: " + err.Error())
		return
	}
	err = nil
	return
}

func (result *PaymentPreSubmitResult) GetUrlToPagseguro(request PaymentRequest) string {
	linkRedirecionamento := ""

	urlCheckout := CheckoutRedirectURL

	if request.IsSandbox == true {
		urlCheckout = CheckoutSandBoxRedirectURL
	}

	if result.Success {
		linkRedirecionamento = fmt.Sprintf(urlCheckout, result.CheckoutResponse.Code)
	}

	return linkRedirecionamento
}

func GetNameStatusPaymentByCode(code int) string {
	if code == TransactionStatusAwaitingPayment {
		return "Aguardando Pagamento"
	}

	if code == TransactionStatusInAnalysis {
		return "Em analise"
	}

	if code == TransactionStatusPaid {
		return "Pago"
	}

	if code == TransactionStatusAvailable {
		return "Disponível"
	}

	if code == TransactionStatusInDispute {
		return "Em disputa"
	}

	if code == TransactionStatusReturned {
		return "Estornado"
	}

	if code == TransactionStatusCanceled {
		return "Cancelado"
	}

	if code == TransactionStatusDebited {
		return "Debitado"
	}

	if code == TransactionStatusTemporaryHold {
		return "Retenção temporária"
	}

	return fmt.Sprintf("Status desconhecido para o código: %d", code)
}

func GetURLPagseguro(IsSandbox bool) string {
	if IsSandbox == true {
		return UrlPagSeguroSandBox
	}

	return UrlPagSeguro
}
