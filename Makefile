SHELL := /bin/bash

protoc:
	protoc --go_out=. --go-grpc_out=. ./sdk/proto/*.proto
	protoc --go_out=. --go-grpc_out=. ./indexer/proto/*.proto
	protoc-go-inject-tag -input=./indexer/proto/structs.pb.go

bindata:
	go-bindata -pkg indexer -o ./indexer/state_migrations.go ./indexer/migrations

postgresql-test:
	docker run --net=host \
		-e POSTGRES_PASSWORD=password \
		-v $(HOME)/pg-data-single:/var/lib/postgresql/data \
		postgres

postgresql-test-admin:
	docker run --net=host \
		-e PGADMIN_DEFAULT_EMAIL=postgres@gmail.com \
		-e PGADMIN_DEFAULT_PASSWORD=postgres \
		dpage/pgadmin4

build-abigen:
	@echo "--> Build abigen"
	./gen --source ./indexer/artifacts/erc20/build/ERC20.json --output ./indexer/artifacts/erc20 --package erc20
	./gen --source ./indexer/artifacts/pancake/build/UniswapRouter.json --output ./indexer/artifacts/pancake --package uniswap
	./gen --source ./indexer/artifacts/pancake/build/UniswapV2Factory.json --output ./indexer/artifacts/pancake --package uniswap
	./gen --source ./indexer/artifacts/pancake/build/UniswapV2Pair.json --output ./indexer/artifacts/pancake --package uniswap