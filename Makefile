
SHA=$(shell git log -n1 --format=format:"%H" 2>/dev/null)
TAG=$(shell git describe --match="v*" --tags 2>/dev/null)

docker-push: autocrat-docker-push

docker-build: autocrat-docker-build

autocrat-docker-push: autocrat-docker-build
	docker push "quay.io/nspain/bagheera-autocrat:latest"
	docker push "quay.io/nspain/bagheera-autocrat:${SHA}"

autocrat-docker-build:
	docker build -t "quay.io/nspain/bagheera-autocrat:latest" -f Dockerfile.autocrat .
	docker tag "quay.io/nspain/bagheera-autocrat:latest" "quay.io/nspain/bagheera-autocrat:${SHA}"

frontend-docker-push: frontend-docker-build
	docker push "quay.io/nspain/bagheera-frontend:latest"
	docker push "quay.io/nspain/bagheera-frontend:${SHA}"

frontend-docker-build:
	docker build -t "quay.io/nspain/bagheera-frontend:latest" -f Dockerfile.frontend frontend
	docker tag "quay.io/nspain/bagheera-frontend:latest" "quay.io/nspain/bagheera-frontend:${SHA}"

up:
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up -d


