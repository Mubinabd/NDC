

DBURL="postgres://postgres:root@localhost:5432/delivery_product_db?sslmode=disable"
exp:
	export DBURL="postgres://postgres:root@localhost:5432/delivery_product_db?sslmode=disable"

mig-up:
	migrate -path muallimah-submodule/migrations -database ${DBURL} -verbose up

mig-down:
	migrate -path muallimah-submodule/migrations -database ${DBURL} -verbose down


mig-create:
	migrate create -ext sql -dir muallimah-submodule/migrations -seq create_table

mig-insert:
	migrate create -ext sql -dir muallimah-submodule/migrations -seq insert_table

CURRENT_DIR=$(shell pwd)

proto-gen:
	protoc --go_out=./ --go-grpc_out=./ internal/pkg/scripts/submodule/*.proto
	
swag-gen:
	~/go/bin/swag init -g /internal/http/router.go -o docs force 1

git:
	git add .
	git commit -m "update"
	git push -u origin main

run:
	clear
	go run main.go