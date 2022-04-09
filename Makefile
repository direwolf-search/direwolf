SHELL = /bin/bash

.PHONY: specs

specs:
# генерируем API + документацию из protobuff allow_merge=true,merge_file_name=api.swagger:
	@protoc \
        --proto_path=protos/ \
        --go_out=internal/protos \
        --go_opt=paths=source_relative \
        --go-grpc_out=internal/protos \
        --go-grpc_opt=paths=source_relative \
        --grpc-gateway_out internal/protos \
        --grpc-gateway_opt logtostderr=true \
        --grpc-gateway_opt paths=source_relative \
        --swagger_out=./docs \
        **/*.proto