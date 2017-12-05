build_protos:
	protoc -I product/ product/product.proto --go_out=plugins=grpc:product