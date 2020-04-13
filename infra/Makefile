VERSION ?= "dev"

all: build deploy

build:
	packer build -var-file=packer/secret.json -var version="$(VERSION)" packer

plan:
	terraform plan -var-file=secret.tfvars

deploy: plan
	terraform apply -var-file=secret.tfvars

destory:
	terraform destroy -var-file=secret.tfvars
