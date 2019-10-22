# Webservice Instruction


1. Run the server.
```sh
$ go run webservice.go
```
2. Run the test 
```sh
$ go test -v webservice_test.go
```
Depends on how many concurrent requests are made, The name or/and joke website may return 429 Error (too many requests).
