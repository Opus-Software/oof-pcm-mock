# Mock PCM Service

Esse projeto tem como objetivo servir de mock para os retornos do serviço de
PCM. Ele auxilia no processo de desenvolvimento para que não seja necessário
bater nas APIs de sandbox para obter um retorno do formato da especificação
para o cenário de POST de novos reportes.

## Como rodar o projeto

No diretório onde estiver o arquivo `main.go` executar:

`go run main.go <retornoHttpEsperado>`

### API /report-api/v1/private/report e API /report-api/v1/opendata/report

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

### API /token/

Implementado um cenário de sucesso onde um token é emitido para fins de teste.
