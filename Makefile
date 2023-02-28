PROJECT=go-boilerplate

DB_HOST=localhost
DB_PORT=5432
DB_USER=user
DB_PASSWORD=pass
DB_NAME=database

_curl_ = docker run --net=host --rm byrnedo/alpine-curl
_mockery_ = mockery --inpackage --all --dir
_go_test_ = ENV=test AWS_ACCESS_KEY_ID=test AWS_SECRET_ACCESS_KEY=test richgo test ./... -count=1 -p=1
_reflex_ = reflex -d none -s -R vendor. -r '.*\.go'
_goose_ = goose -dir ./migration postgres "host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable"

upgrade-all: upgrade install

upgrade:
	go get -u ./...

install-tools:
	cat tools/tools.go | grep "_" | awk -F '"' '{print $$2}' | xargs -L1 go get -u

clean:
	rm -rf vendor/

install: install-tools
	go mod vendor && go mod tidy

build/install:
	go install --ldflags='-w -s -extldflags "-static"' -v -a

run/api:
	go run main.go api

lint:
	go list ./... | grep -v go-boilerplate/docs/swagger | xargs -L1 staticcheck -f stylish -fail all -tests

docker/build:
	docker build -t $(DOCKER_IMAGE) .

docker-dependencies/up:
	docker-compose up --build -d
	sleep 10
	make db/migrate

docker-dependencies/down:
	docker-compose down -v

tag:
	docker tag $(DOCKER_IMAGE) $(DOCKER_IMAGE_TAG)

test: docker-dependencies/down docker-dependencies/up
	$(_go_test_)

test/coverage: docker-dependencies/down docker-dependencies/up
	$(_go_test_) -coverprofile cover.out

test/coverage/html: test/coverage
	go tool cover -html cover.out

swagger:
	swag init -generalInfo api/api.go -output ./docs/swagger

db/create-migration:
	$(_goose_) create $(MIGRATION_NAME) sql

db/migrate:
	$(_goose_) up

db/migrate-down:
	$(_goose_) down

db/migration-status:
	$(_goose_) status

db/er-diagram: docker-dependencies/down docker-dependencies/up
	docker run -v $(PWD)/migration:/share --net go-boilerplate_default schemacrawler/schemacrawler /opt/schemacrawler/schemacrawler.sh --server=postgresql --host=go-boilerplate_postgres_1 --user=user --password=pass --database=database --info-level=standard --command=schema --outputformat=png --output-file /share/database-diagram.png

mockery:
	$(_mockery_) repository && $(_mockery_) facade
