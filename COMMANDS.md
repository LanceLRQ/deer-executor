1、运行Hello World测试

``` go run main.go run --tin ./test/cases/1.in --tout ./test/cases/1.out --clean  ./test/scripts/hello.c```

2、运行评测

```go run main.go run --tin ./test/cases/1.in --tout ./test/cases/1.out --clean ./test/scripts/ac.c```

3、运行评测并持久化到文件（带签名和数字校验）

```go run main.go run --tin ./test/cases/1.in --tout ./test/cases/1.out --clean --p ./result --sign --public-key ./test/certs/test.pem --private-key ./test/certs/test.key  ./test/scripts/ac.c```