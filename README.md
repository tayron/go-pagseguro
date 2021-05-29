# GO Pagseguro

Biblioteca de integração com Pagseguro que permite:
- Gerar link com os dados do pedido para redirecionamento do checkout do Pagseguro para que o cliente efetue o pagmento.
- Criar serviço para receber notificação com os status de pagamento do pagseguro

Esta biblioteca é uma versão atualizada da biblioteca [diegomgarcia/pagsego](https://github.com/diegomgarcia/pagsego).

## Teste untiário

### Configuração
Abra o arquivo: **tests/pagseguro_test.go** e altere as variaveis com os dados de configuração da sua conta no [Sandbox do Pagseguro](https://acesso.pagseguro.uol.com.br/sandbox).
### Executando teste unitário em ambiente de desenvolvimento
```sh
go test tests/*
```
### Gerando binário do teste unitário para uso em produção
```sh
go test tests/* -c -o pagseguro_unit_test
```
### Executando teste unitário através do binário
```sh
./pagseguro_unit_test
``` 

## Exemplo de utilização da biblioteca em um projeto
Abaixo segue exemplo de utilização real da biblioteca, até agora utilizado somente em ambiente Sandbox do Pagseguro.
### Implementado método para pagamento
Exemplo de implementação do método irá redirecionar o usuário para o checkout do pagseguro e gravar os dados do pagamento em um banco de dados da aplicação.
```go
func efetuarPagamento(w http.ResponseWriter, r *http.Request) {
	loginUsuario := session.GetDadoSessao("login", w, r)

	var usuarioModel models.Usuario = models.Usuario{
		Login: loginUsuario,
	}

	usuario := usuarioModel.BuscarPorLogin()
	configuracao := models.BuscarConfiguracao()

	email := configuracao.EmailGatewayPagamento
	token := configuracao.TokenGatewayPagamento
	urlParaRedirecionamento := configuracao.UrlRedirecionamentoGatewayPagamento
	urlParaNotificacao := configuracao.UrlNotificacaoGatewayPagamento
	idTransacao := util.GerarHashAleatorio()
	quantidadeProduto := configuracao.QuantidadeProduto
	valor := configuracao.ValorProduto

	pedido := pagseguro.NewPaymentRequest(token, email, idTransacao, urlParaRedirecionamento, urlParaNotificacao)
	pedido.IsSandbox = true
	pedido.AddItem(configuracao.IdProduto, configuracao.NomeProduto, valor, configuracao.QuantidadeProduto)
	pedido.SetBuyer("Comprador Teste", "c63252160473116878087@sandbox.pagseguro.com.br")

	result := pedido.Submit()
	linkPagamento := gravarDadosPagamento(idTransacao, usuario, valor, quantidadeProduto, *pedido, *result)

	if linkPagamento != "" {
		http.Redirect(w, r, linkPagamento, 302)
	}
}

func gravarDadosPagamento(idTransacao string, usuario models.Usuario,
	valor float64, quantidadeProduto int, request pagseguro.PaymentRequest,
	retornoPagseguro pagseguro.PaymentPreSubmitResult) string {

	pagamento := models.Pagamento{
		Usuario:                usuario,
		GatewayPagamento:       "PagSeguro",
		IdGatewayPagamento:     "",
		IdTransacao:            idTransacao,
		ValorProduto:           valor,
		QuantidadeProduto:      quantidadeProduto,
		ValorBruto:             0.00,
		ValorDesconto:          0.00,
		ValorTaxa:              0.00,
		ValorLiquido:           0.00,
		StatusGatewayPagamento: "",
		StatusPagamento:        false,
		Observacao:             "",
		Log:                    retornoPagseguro.XML,
	}

	pagamento.CriarPagamento()

	if retornoPagseguro.Success {
		return retornoPagseguro.GetUrlToPagseguro(request)
	}

	return ""

```

Para a struct **Pagamento** implemente o método **CriarPagamento()** para gravar os dados de pagamento no banco de dados.

### Implementando método que irá receber as notificações do Pagseguro
```go
func ReceberNotificacaoPagamento(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	sandbox := true
	codigoNotificacao := r.PostForm.Get("notificationCode") // 414284-5CC348C348C8-5774999F9713-DB8A2E
	//tipoNotificacao := r.PostForm.Get("notificationType")   // transaction

	configuracao := models.BuscarConfiguracao()

	email := configuracao.EmailGatewayPagamento
	token := configuracao.TokenGatewayPagamento

	retornoPagseguro, err := pagseguro.FetchTransactionInfo(token, email, codigoNotificacao, sandbox)

	urlGatewayPagmento := pagseguro.GetURLPagseguro(sandbox)
	util.EnableCors(&w, urlGatewayPagmento)

	if err != nil {
		util.ExibirLogErro(err.Error())
		return
	}

	statusPagamentoPagseguro := pagseguro.GetNameStatusPaymentByCode(retornoPagseguro.Status)
	statusPagamento := false

	if retornoPagseguro.Status == pagseguro.TransactionStatusPaid {
		statusPagamento = true
	}

	if retornoPagseguro.Status == pagseguro.TransactionStatusAvailable {
		statusPagamento = true
	}

	if retornoPagseguro.Status == pagseguro.TransactionStatusDebited {
		statusPagamento = true
	}

	valorBruto, _ := strconv.ParseFloat(retornoPagseguro.GrossAmount, 64)
	valorDesconto, _ := strconv.ParseFloat(retornoPagseguro.DiscountAmount, 64)
	valorTaxa, _ := strconv.ParseFloat(retornoPagseguro.FeeAmount, 64)
	valorLiquido, _ := strconv.ParseFloat(retornoPagseguro.NetAmount, 64)

	pagamento := models.Pagamento{
		IdGatewayPagamento:     codigoNotificacao,
		IdTransacao:            retornoPagseguro.Reference,
		ValorBruto:             valorBruto,
		ValorDesconto:          valorDesconto,
		ValorTaxa:              valorTaxa,
		ValorLiquido:           valorLiquido,
		StatusGatewayPagamento: statusPagamentoPagseguro,
		StatusPagamento:        statusPagamento,
		Observacao:             "",
		Log:                    retornoPagseguro.XML,
	}

	pagamento.AtualizarPagamentoPorIdTransacao()
}
```
Para a struct **Pagamento** implemente o método **AtualizarPagamentoPorIdTransacao()** para atualizar os dados de pagamento recebido do Pagseguro.