1、运行简单评测

```go run main.go run --config ./data/problems/APlusB/problem.json --clean ./data/codes/APlusB/ac.c```

3、运行评测并持久化到文件（带签名和数字校验）

```go run main.go run --config ./data/problems/APlusB/problem.json --clean --p ./result --sign --public-key ./data/certs/test.pem --private-key ./data/certs/test.key ./data/codes/APlusB/ac.c```