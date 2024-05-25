build:
   export GO111MODULE=on
   export CGO_ENABLED=1
   GO_BUILD -o bin/handler/bootstrap handler/main.go
   chmod +x bin/handler/bootstrap
   zip -j bin/handler.zip bin/handler/bootstrap

deploy:
	npx sls deploy -s $(STAGE) -v