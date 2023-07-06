package test

import (
	"bytes"
	"crypto/ed25519"
	"myConnect/jwt"
	"myConnect/tlog"
	"testing"
)

func TestMain(m *testing.M) {
	tlog.Init()
	m.Run()
}
func TestSignJWT(t *testing.T) {
	seed := "12345678901234567890123456789012"
	publicKey, privateKey, err := ed25519.GenerateKey(bytes.NewReader([]byte(seed)))
	if err != nil {
		t.Fatalf("ed25519.GenerateKey err is %s ", err.Error())
	}
	//publicKey := privateKey.Public()
	t.Logf("privateKey is %x publicKey is %x ", privateKey, publicKey)
	res := jwt.SignJWT("act", "sub", "aud", 1000, ed25519.PublicKey(publicKey), ed25519.PrivateKey(privateKey))
	t.Logf("SignJWT res is %s ", res)
	jwt.VerifyJWT(res)
}
