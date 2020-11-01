build:
	go build -o terraform-provider-harness .

install: build
	mv terraform-provider-harness ~/.terraform.d/plugins