package client

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/LanceLRQ/deer-executor/executor"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

//生成RSA私钥和公钥，保存到文件中
func generateRSAKey(bits int){
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err!=nil{
		panic(err)
	}
	X509PrivateKey := x509.MarshalPKCS1PrivateKey(privateKey)
	privateFile, err := os.Create("private.pem")
	if err!=nil {
		panic(err)
	}
	defer privateFile.Close()
	privateBlock:= pem.Block{Type: "RSA Private Key",Bytes:X509PrivateKey}
	_ = pem.Encode(privateFile,&privateBlock)

	publicKey := privateKey.PublicKey
	X509PublicKey,err := x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil{
		panic(err)
	}
	publicFile, err := os.Create("public.pem")
	if err!=nil{
		panic(err)
	}
	defer publicFile.Close()
	publicBlock:= pem.Block{Type: "RSA Public Key",Bytes:X509PublicKey}
	//保存到文件
	_ = pem.Encode(publicFile,&publicBlock)
}

func MakeConfigFile(c *cli.Context) error {
	config := executor.NewSession()
	config.TestCases = []executor.TestCase {
		{
			Id: "1",
			TestCaseIn: "",
			TestCaseOut: "",
		},
	}
	output := c.String("output")
	if output != "" {
		_, err := os.Stat(output)
		if os.IsExist(err) {
			log.Fatal("output file exists")
			return nil
		}
		fp, err := os.OpenFile(output, os.O_WRONLY | os.O_CREATE, 0644)
		if err != nil {
			log.Fatalf("open output file error: %s\n", err.Error())
			return nil
		}
		defer fp.Close()
		_, err = fp.WriteString(executor.ObjectToJSONStringFormatted(config))
		if err != nil {
			return err
		}
	} else {
		fmt.Println(executor.ObjectToJSONStringFormatted(config))
	}
	return nil
}

func GenerateRSA(c *cli.Context) error {
	bit := c.Int("bit")
	if bit < 2048 || bit % 2 != 0 {
		return fmt.Errorf("RSA bit must larger than 2048 (or equal)")
	}
	generateRSAKey(bit)
	return nil
}