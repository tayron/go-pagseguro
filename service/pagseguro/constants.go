package pagseguro

const (
	ShippingPAC   = 1
	ShippingSEDEX = 2
	ShippingOther = 3
	XMLHeader     = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>`

	CheckoutURL     = "https://ws.pagseguro.uol.com.br/v2/checkout"
	TransactionsURL = "https://ws.pagseguro.uol.com.br/v2/transactions/notifications"

	CheckoutSandBoxURL     = "https://ws.sandbox.pagseguro.uol.com.br/v2/checkout"
	TransactionsSandBoxURL = "https://ws.sandbox.pagseguro.uol.com.br/v2/transactions/notifications"

	CheckoutRedirectURL        = "https://pagseguro.uol.com.br/v2/checkout/payment.html?code=%s"
	CheckoutSandBoxRedirectURL = "https://sandbox.pagseguro.uol.com.br/v2/checkout/payment.html?code=%s"

	UrlPagSeguro        = "https://pagseguro.uol.com.br"
	UrlPagSeguroSandBox = "https://sandbox.pagseguro.uol.com.br"

	TransactionTypePayment = 1

	TransactionStatusAwaitingPayment = 1
	TransactionStatusInAnalysis      = 2
	TransactionStatusPaid            = 3
	TransactionStatusAvailable       = 4
	TransactionStatusInDispute       = 5
	TransactionStatusReturned        = 6
	TransactionStatusCanceled        = 7
	TransactionStatusDebited         = 8
	TransactionStatusTemporaryHold   = 9

	PaymentMethodCreditCardVisa             = 101
	PaymentMethodCreditCardMasterCard       = 102
	PaymentMethodCreditCardAMEX             = 103
	PaymentMethodCreditCardDiners           = 104
	PaymentMethodCreditCardHipercard        = 105
	PaymentMethodCreditCardAura             = 106
	PaymentMethodCreditCardElo              = 107
	PaymentMethodCreditCardPLENOCard        = 108
	PaymentMethodCreditCardPersonalCard     = 109
	PaymentMethodCreditCardJCB              = 110
	PaymentMethodCreditCardDiscover         = 111
	PaymentMethodCreditCardBrasilCard       = 112
	PaymentMethodCreditCardFORTBRASIL       = 113
	PaymentMethodCreditCardCARDBAN          = 114
	PaymentMethodCreditCardVALECARD         = 115
	PaymentMethodCreditCardCabal            = 116
	PaymentMethodCreditCardMais             = 117
	PaymentMethodCreditCardAvista           = 118
	PaymentMethodCreditCardGRANDCARD        = 119
	PaymentMethodBoletoBradesco             = 201
	PaymentMethodBoletoSantander            = 202
	PaymentMethodDebitoOnlineBradesco       = 301
	PaymentMethodDebitoOnlineItau           = 302
	PaymentMethodDebitoOnlineUnibanco       = 303
	PaymentMethodDebitoOnlineBancoDoBrasil  = 304
	PaymentMethodDebitoOnlineBancoReal      = 305
	PaymentMethodDebitoOnlineBanrisul       = 306
	PaymentMethodDebitoOnlineHSBC           = 307
	PaymentMethodSaldoPagSeguro             = 401
	PaymentMethodOiPaggo                    = 501
	PaymentMethodDepositoContaBancoDoBrasil = 701
	PaymentMethodDepositoContaHSBC          = 702
)
