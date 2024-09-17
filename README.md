Dependencias:
- https://buf.build/docs/installation
- https://go.dev/doc/install (go 1.23+)

Para rodar use: `go run .`
Para regerar o proto: `buf generate`
Se tiver alguma alteração de imports no proto use: `buf dep update && buf generate`

Sugestão de app para requests: Postman, Inmsomnia, grpcurl, curl
