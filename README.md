## how to run

```
docker-compose up -d
```

```
go install github.com/ftamhar/nrpc/protoc-gen-nrpc@v0.1.0
```

```
protoc --go_out=. --go_opt=paths=source_relative \
    --nrpc_out=. --nrpc_opt=paths=source_relative \
    proto/*/*.proto
```

```
go run cmd/server/main.go
```


in another terminal 

```
go run cmd/client/main.go
```

### next step
go to http://localhost:8080 to get my name

go to http://localhost:8080/upload to upload image
