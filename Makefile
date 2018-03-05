all: client server realserver

%: %.go p/svc.pb.go Gopkg.lock
	go build -o $@ $<

p/svc.pb.go: p/svc.proto
	protoc p/svc.proto --go_out=plugins=grpc:.
