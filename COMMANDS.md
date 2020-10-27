## 使用简单评测

1、运行简单评测

```go run main.go run --config ./data/problems/APlusB/problem.json --clean ./data/codes/APlusB/ac.c```

2、运行评测并持久化到文件

```go run main.go run --config ./data/problems/APlusB/problem.json --clean --persistence ./result ./data/codes/APlusB/ac.c```

3、运行评测并持久化到文件（带签名和数字校验），需要输入GPG密钥文件密码

```go run main.go run --config ./data/problems/APlusB/problem.json --clean --persistence ./result --sign --key YOUR_GPG_KEY ./data/codes/APlusB/ac.c```

## GPG 密钥生成

```准备：操作系统需要安装opengpg```

1、生成密钥：`gpg --full-generate-key`，按照提示操作，密钥类型默认即可，密钥长度大于等于2048，有效期自行决定，需要设置私钥保护密码（建议设置复杂的密码）


2、列出已生成的密钥：`gpg --list-secret-keys --keyid-format LONG`

```
/Users/xxxx/.gnupg/pubring.kbx
------------------------------------
sec   rsa2048/F896F5F1F6AFF7FA 2020-10-26 [SC]
      A0E905E7F0781682F19F3F31F896F5F1F6AFF7FA
uid                 [ 绝对 ] Deer-executor Test (Hello world) <test@wejudge.net>
ssb   rsa2048/933D155C3D66365F 2020-10-26 [E]
```

3、使用步骤2拿到的`rsa2048/F896F5F1F6AFF7FA`，导出GPG私钥

```
gpg --armor --output private-key.txt --export-secret-keys F896F5F1F6AFF7FA
```

`private-key.txt`即为`YOUR_GPG_KEY`所要用到的GPG私钥，请妥善保护好这个文件。

## 题目打包

```
 go run main.go pack --sign --key YOUR_GPG_KEY ./data/problems/APlusB/problem.json ./a+b.problem 
```

如果去掉sign参数则不签名。

运行评测时，config参数支持识别是否为题目包文件，是则自动解包校验并释放到/tmp目录下，可以通过--work-dir参数设定

```go run main.go run --config ./a+b.problem --clean ./data/codes/APlusB/ac.c```

