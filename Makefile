start-dev:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d

start:
	docker-compose up -d --abort-on-container-exit

test:
	docker-compose -f docker-compose.yml -f docker-compose.test.yml up --exit-code-from test
