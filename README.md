# missionaries-and-cannibals

## building proto
```
go get -v -u github.com/golang/protobuf/protoc-gen-go
protoc -I=./errors --go_out=./errors ./errors/errors.proto
```

## running
```
go run main.go
```