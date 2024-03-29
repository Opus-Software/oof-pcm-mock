# Mock PCM Service

Esse projeto tem como objetivo servir de mock para os retornos do serviço de
PCM do Open Finance Brasil (OFB) e do Open Insurance (OPIN). Ele auxilia no
processo de desenvolvimento para que não seja necessário bater nas APIs de
sandbox para obter um retorno do formato da especificação para o cenário de
POST de novos reportes.

**Observação:** Ao rodar a ferramenta, os endpoints de OFB e OPIN estarão
disponíveis simultaneamente uma vez que as rotas são distintas entre as
especificações de PCM.

## Sobre o projeto

Este projeto é uma ferramenta auxiliar utilizada pela equipe do 
[OPUS Open Insurance](https://www.opus-software.com.br/open-insurance/) e
[OPUS Open Finance](https://www.opus-software.com.br/solucoes-cases/opus-open-finance/),
a solução de software mais completa para instituições cumprirem as regulações
do OFB e OPIN.

Nossa plataforma é uma solução completa que resolve os requisitos de segurança e negócio,
incluindo implementação das APIs obrigatórias,
gestão de consentimentos 100% transparente,
integração com a PCM, e muito mais, 
podendo ser utilizada em qualquer nuvem, on-premise, ou como SaaS.

Entre em contato e conheça mais sobre nossas soluções!

## Como rodar o projeto

No diretório onde estiver o arquivo `main.go` executar:

`go run main.go <retornoHttpEsperado>`

### APIs /report-api/v1/private/report e /report-api/v1/opendata/report de OFB e /report-api/v1/server-batch de OPIN

Os seguintes cenários de retorno de sucesso foram implementados:

- `200`: Default
- `207`: Alguns objetos do array de retorno terão o campo `message` definido
  com um exemplo.

Os seguintes cenários de retorno de erro foram implementados:

- `400`: Invalid payload format: MUST be an array
- `401`: Unauthorized
- `403`: Forbidden
- `406`: Content type not accepted
- `413`: Record limit exceeded
- `415`: Unsupported Media Type
- `429`: Unsupported Media Type
- `500`: Internal Server Error

### API /token/ de OFB

Implementado um cenário de sucesso onde um token é emitido para fins de teste.
