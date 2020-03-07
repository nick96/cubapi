start-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

start:
	docker-compose up -d --abort-on-container-exit

test:
	go test -short -v github.com/nick96/cubapi/...

test-docker:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --exit-code-from test
