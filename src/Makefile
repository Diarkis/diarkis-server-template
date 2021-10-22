.PHONY: init
init: # download diarkis-cli and mars binary
	./script/install_tools.sh

.PHONY: build-local
build-local:
	./diarkis-cli/bin/diarkis-cli build -c local-build.yaml --host builder.diarkis.io

.PHONY: build-container
build-container: 
	./diarkis-cli/bin/diarkis-cli build -c container-build.yaml --host builder.diarkis.io
	docker build --platform=linux/amd64 -f docker/http/Dockerfile remote_bin -t AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker build --platform=linux/amd64 -f docker/udp/Dockerfile remote_bin -t AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker build --platform=linux/amd64 -f docker/tcp/Dockerfile remote_bin -t AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker build --platform=linux/amd64 -f docker/mars/Dockerfile remote_bin -t AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/mars

.DEFAULT_GOAL := build-container

.PHONY: push-container
push-container:
	docker push AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/http
	docker push AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/udp
	docker push AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/tcp
	docker push AWS_ACCOUNT_NUM.dkr.ecr.ap-northeast-1.amazonaws.com/mars
