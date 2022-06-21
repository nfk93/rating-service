generate-api:
	openapi-generator generate \
		-c api/openapi-generator-config.yaml \
		-i api/api.yaml \
		-g go-server \
		--invoker-package todolist \
		--git-repo-id rating-service \
		--git-user-id nfk93 \
		--git-host github.com \
		--model-package model \
		--api-package api \
		-p servicename=Todolist \
		-p servicerPackage=service \
		-o .