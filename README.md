# dm-go

Download [protoc](https://github.com/protocolbuffers/protobuf/releases/latest)

Install go plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.26
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.1
```

Update `PATH` env:
```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

```bash
git clone https://github.com/Military-Doctor/dm-go
```