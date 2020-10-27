package client

import (
	"fmt"
	"github.com/urfave/cli/v2"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"io/ioutil"
	"os"
)

////privateKey, err := persistence.ReadPemFile("./data/certs/test.key")
////if err != nil {
////	return err
////}
////sign, err := persistence.RSA2048SignString("Hello World", pkey)
//////if err != nil {
////	return err
////}
////fmt.Println(hex.EncodeToString(sign))
////
////publicKey, err := persistence.ReadPemFile("./data/certs/test.pem")
////if err != nil {
////	return err
////}
////err = persistence.RSA2048VerifyString("Hello World", sign, publicKey)
////if err == nil {
////	fmt.Println("Yes!")
////}
//
////rel, err := persistence.SHA256Streams([]io.Reader{
////	bytes.NewReader(publicKey),
////	bytes.NewReader(privateKey),
////})
////if err != nil {
////	return err
////}
////fmt.Println(hex.EncodeToString(rel))
////_, err := judge_result.ReadJudgeResult("./result")
////if err != nil {
////	return err
////}
////fmt.Println(executor.ObjectToJSONStringFormatted(rst))
//
//pem, err := persistence.GetDigitalPEMFromFile("./data/certs/test.pem", "./data/certs/test.key")
//if err != nil {
//	return err
//}
//session, err := executor.NewSession("./data/problems/APlusB/problem.json")
//if err != nil {
//	return err
//}
//options := problems.ProblemPersisOptions{
//	DigitalSign: true,
//	DigitalPEM: *pem,
//	OutFile: "./a+b.problem",
//}
//err = problems.PackProblems(session, options)
//if err != nil {
//	return err
//}
//return nil

func Test(c *cli.Context) error {
	keyRingReader, err := os.Open("private-key.txt")
	if err != nil {
		return err
	}

	elist, err := openpgp.ReadArmoredKeyRing(keyRingReader)
	if err != nil {
		return err
	}
	pkey := elist[0].PrimaryKey

	fp, err := os.Create("/tmp/pub.key")
	if err != nil {
		return err
	}
	w, err := armor.Encode(fp, openpgp.PublicKeyType, nil)
	if err != nil {
		return err
	}
	err = pkey.Serialize(w)
	if err != nil {
		return err
	}
	w.Close()
	fp.Close()

	defer os.Remove("/tmp/pub.key")

	data, err := ioutil.ReadFile("/tmp/pub.key")
	if err != nil {
		return err
	}

	fmt.Println(string(data))
	//
	////fmt.Println(elist[0].PrimaryKey.PublicKey)
	//for _, v := range elist[0].Identities {
	//	fmt.Println(v.Name)
	//	fmt.Println(v.Signatures)
	//	//fmt.Println(v.SelfSignature)
	//	fmt.Println(v.UserId)
	//}

	//pem := persistence.DigitalSignPEM {
	//	PrivateKey: keys[0].PrivateKey.PrivateKey.(*rsa.PrivateKey),
	//	PublicKey: keys[0].PublicKey.PublicKey.(*rsa.PublicKey),
	//}
	//
	//sign, err := persistence.RSA2048SignString("Hello World", pem.PrivateKey)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(hex.EncodeToString(sign))
	//
	//err = persistence.RSA2048VerifyString("Hello World", sign, pem.PublicKey)
	//if err != nil {
	//	return err
	//}
	//fmt.Println("Yes!")
	//
	//passphrase, err := gopass.GetPasswdPrompt("please input passphrase of key>", true, os.Stdin, os.Stdout)
	//if err != nil {
	//	return err
	//}
	//fmt.Println(string(passphrase))
	return nil
}