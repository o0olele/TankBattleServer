package common

import (
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"io"
	"io/ioutil"
	"os"

	"github.com/golang/glog"
)

const salt string = "cbh8a932gnvf"

func GetToken(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return string(h.Sum([]byte(salt)))
}

func CheckToken(s, token string) bool {
	tok := GetToken(s)
	if tok == token {
		return true
	}
	return false
}

func ReadAll(filePth string) ([]byte, error) {
	f, err := os.Open(filePth)
	if nil != err {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

// GetToken 获取token
func GetTokenSSL(bytes []byte) (string, error) {
	publickey, err := ReadAll("./rsa_public_key.pem")
	if nil != err {
		return "", err
	}

	token, err := RsaEncrypt(bytes, publickey)
	if nil != err {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(token), nil
}

// DecryptToken 对加密的token进行解码
func DecryptTokenSSL(token string) (*Token, error) {
	bytes, err := base64.StdEncoding.DecodeString(token)
	if nil != err {
		glog.Error("[common] base64 decode fail ", err)
		return nil, err
	}

	privatekey, err := ReadAll("./rsa_private_key.pem")
	if nil != err {
		return nil, err
	}

	result, err := RsaDecrypt(bytes, []byte(privatekey))
	if nil != err {
		glog.Error("[common] rsa decode fail ", err)
		return nil, err
	}

	var info Token
	err = json.Unmarshal(result, &info)
	if nil != err {
		glog.Error("[common] json decode fail ", err)
		return nil, err
	}
	return &info, nil
}

// 私钥生成
// openssl genrsa -out rsa_private_key.pem 1024
// 公钥生成
// openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem

// 加密
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if nil == block {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if nil != err {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 解密
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if nil == block {
		return nil, errors.New("private key error!")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if nil != err {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}
