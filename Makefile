gen:
	export PATH="$(PATH):$(HOME)/.local/bin"	
	(which protoc > /dev/null) || (curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.15.8/protoc-3.15.8-linux-x86_64.zip && \
	unzip protoc-3.15.8-linux-x86_64.zip -d $(HOME)/.local)
	go get -u google.golang.org/protobuf/cmd/protoc-gen-go
	go get -u google.golang.org/grpc/cmd/protoc-gen-go-grpc

	(which openapi-generator-cli > /dev/null) || \
	(curl -LO https://raw.githubusercontent.com/OpenAPITools/openapi-generator/master/bin/utils/openapi-generator-cli.sh > ~/.local/bin/openapi-generator-cli && \
	chmod u+x ~/.local/bin/openapi-generator-cli && \
	apt-get update && apt-get upgrade && apt-get --fix-broken install && (which mv > /dev/null || apt-get install maven))

	go generate
