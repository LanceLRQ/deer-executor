1、运行简单评测

```go run main.go run --config ./test/configs/simple.conf --clean ./test/codes/AplusB/ac.c```

3、运行评测并持久化到文件（带签名和数字校验）

```go run main.go run --config ./test/configs/simple.conf --clean --p ./result --sign --public-key ./test/certs/test.pem --private-key ./test/certs/test.key ./test/codes/AplusB/ac.c```