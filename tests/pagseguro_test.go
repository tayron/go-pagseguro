package tests

import (
	"testing"

	"github.com/tayron/go-pagseguro/service/pagseguro"
)

func TestPagseguro(t *testing.T) {
	t.Parallel() // Teste em paralelo

	pagseguroEmailCustumer := "hash@sandbox.pagseguro.com.br"
	pagseguroEmailSeller := "seu_email_pagseguro@gmail.com"
	pagseguroTokenSeller := "seu_token_pagseguro_aqui"
	pagseguroUrlRedirectToYoursite := "url_do_seu_site"
	pagseguroUrlNotificationIntoYoursite := "url_do_seu_site_para_receber_notificacao"
	pagseguroAplicationId := "id_da_sua_aplicacao_no_pagseguro"

	req := pagseguro.NewPaymentRequest(pagseguroTokenSeller, pagseguroEmailSeller, pagseguroAplicationId, pagseguroUrlRedirectToYoursite, pagseguroUrlNotificationIntoYoursite)
	req.IsSandbox = true
	req.AddItem("Produto", "Descrição do produto", 23.56, 1)
	req.SetBuyer("Comprador Teste", pagseguroEmailCustumer).SetCPF("94776272083")

	result := req.Submit()

	if !result.Success {
		t.Error(result.CheckoutResponse.Code)
	}

}
