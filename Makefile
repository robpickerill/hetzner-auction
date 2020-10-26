S3_BUCKET = a-bucket
TEMPLATE_ENCRYPTED = template.enc.yml
TEMPLATE = template.yml
APPLICATION_NAME = hetzner-server-auction


.PHONY: decrypt
decrypt:
	sops --decrypt \
		$(TEMPLATE_ENCRYPTED) > $(TEMPLATE)

.PHONY: encrypt
encrypt:
	sops $(TEMPLATE_ENCRYPTED)

.PHONY: build
build:
	sam build

.PHONY: package
package: build
	sam package \
		--template-file template.yaml \
		--s3-bucket $(S3_BUCKET) \
		--s3-prefix $(APPLICATION_NAME) \
		--output-template-file output.yaml \
		--region $(REGION)

.PHONY: invoke
invoke: build
	sam local \
		invoke

.PHONY: deploy
deploy: package
	sam deploy \
		--stack-name $(APPLICATION_NAME) \
		--template-file output.yaml \
		--capabilities CAPABILITY_IAM \
		--tags \
				application=$(APPLICATION_NAME)
		--region $(REGION)
