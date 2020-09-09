all: proto
	mv proto/*.pb.go .
	go test

proto:
	protoc --go_out=./ proto/pb.proto
	protoc --go_out=./ proto/hello.proto

clean:
	rm -f pb *.out *.pb.go proto/*.pb.go


.PHONY: proto
