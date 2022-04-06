.DEFAULT_GOAL := swagger

install_swagger:
	go get -u github.com/go-swagger/go-swagger/cmd/swagger

swagger:
	@echo Ensure you have the swagger CLI or this command will fail.
	@echo You can install the swagger CLI with: go get -u github.com/go-swagger/go-swagger/cmd/swagger
	@echo ....

	swagger generate spec -o ./swagger.yaml --scan-models

run:
	BIND_ADDRESS=3000 LOG_LEVEL=trace air -c air.toml

generate_client:
	cd sdk && swagger generate client -f ../swagger.yaml -A product-api

test:
	cd tests && go test -v