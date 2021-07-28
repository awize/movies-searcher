all: bin/api
test: lint unit-test
generate-mock:
	go generate -v ./...

PLATFORM=local

bin/api:
	docker build . --target bin \
	--output bin/ \
	--platform ${PLATFORM}

unit-test:
	docker build . --target unit-test

unit-test-coverage:
	docker build . --target unit-test-coverage \
	--output coverage/
	cat coverage/cover.out

lint:
	@docker build . --target lint