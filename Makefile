S3_BUCKET = linuxadept-artifacts
TEMPLATE_ENCRYPTED = template.enc.yml
TEMPLATE = template.yml
APPLICATION_NAME = hetzner-server-auction
REGION=eu-west-1


.PHONY: decrypt
decrypt:
	sops --decrypt \
		$(TEMPLATE_ENCRYPTED) > $(TEMPLATE)

.PHONY: encrypt
encrypt:
	sops $(TEMPLATE_ENCRYPTED)

.PHONY: build
build:
	GOOS=linux \
	GOARCH=amd64 \
	go build \
	-o cmd/hetzner/hetzner \
		./cmd/hetzner

.PHONY: package
package: build
	sam package \
		--template-file template.yml \
		--s3-bucket $(S3_BUCKET) \
		--s3-prefix $(APPLICATION_NAME) \
		--output-template-file output.yml \
		--region $(REGION)

.PHONY: invoke
invoke: build
	sam local \
		invoke

.PHONY: deploy
deploy: package
	sam deploy \
		--stack-name $(APPLICATION_NAME) \
		--template-file output.yml \
		--capabilities CAPABILITY_IAM \
		--tags \
				application=$(APPLICATION_NAME) \
		--region $(REGION)
