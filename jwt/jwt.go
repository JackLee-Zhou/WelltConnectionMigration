package jwt

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"github.com/multiformats/go-multibase"
	"myConnect/tlog"
	"strings"
	"time"
	"unicode/utf8"
)

// https://docs.walletconnect.com/2.0/specs/clients/core/crypto/crypto-authentication

func SignJWT(act, sub, aud string, ttl int64, publicKey ed25519.PublicKey, privateKey ed25519.PrivateKey) string {
	header := &JWTHeader{
		Alg: JWT_ALG,
		Typ: JWT_TYP,
	}
	utf8.DecodeRune(pub)
	iat := time.Now().UnixMilli()
	exp := iat + ttl
	iss := encodeIss(publicKey)
	payload := &JWTPayload{
		Iat: iat,
		Exp: exp,
		Iss: iss,
		Sub: sub,
		Aud: aud,
		Act: act,
	}
	tlog.Infof("SignJWT PayLoad is %+v ", payload)
	data := encodeData(header, payload)
	signature := ed25519.Sign(privateKey, data)
	return encodeJWT(header, payload, signature)
}
func VerifyJWT(jwtStr string) {
	sign := decodeJWT(jwtStr)
	ok := verify(sign)
	if !ok {
		tlog.Infof("verify not ok jwt is %s ", jwtStr)
	}
}
func verify(sign *JWTSigned) bool {
	if sign.Header == nil || sign.Header.Alg != JWT_ALG || sign.Header.Typ != JWT_TYP {
		tlog.Errorf("verify sign.Header err header is %+v", sign.Header)
		return false
	}
	publicKey := decodeIss(sign.Payload.Iss)
	data := encodeData(sign.Header, sign.Payload)
	return ed25519.Verify(publicKey, data, sign.Signature)
}

// iss Example: did:key:z6Mkf6PpiWkF1VyiRDWHoGRbXnoK4dAAoUUPkFTNR4YxjUne
func decodeIss(iss string) []uint8 {
	issUser := strings.Split(iss, DID_DELIMITER)
	if len(issUser) != 3 {
		tlog.Errorf("decodeIss err iss is %s ", iss)
		return nil
	}
	prefix := issUser[0]
	method := issUser[1]
	mutibaseStr := issUser[2]
	if prefix != DID_PREFIX || method != DID_METHOD {
		tlog.Errorf("Issuer must be a DID with method key iss is %s ", iss)
		return nil
	}
	base := mutibaseStr[:1]
	if base != string(MULTIBASE_BASE58BTC_PREFIX) {
		tlog.Errorf("Issuer must be a multibase with encoding base58btc, iss is %s ", iss)
		return nil
	}
	// MULTICODEC_ED25519_ENCODING
	multicodec, err := multibase.Encode(multibase.Base58BTC, []byte(mutibaseStr[1:]))
	if err != nil {
		tlog.Errorf("Issuer must be a multibase with encoding base58btc, iss is %s ", iss)
		return nil
	}
	keyType := multicodec[:2]
	if keyType != MULTICODEC_ED25519_HEADER {
		tlog.Errorf("Issuer must be a public key with type \"Ed25519\", iss is %s ", iss)
		return nil
	}
	publicKey := multicodec[2:]
	if len(publicKey) != MULTICODEC_ED25519_LENGTH {
		tlog.Errorf("Issuer must be a public key with length 32, iss is %s ", iss)
		return nil
	}
	return []uint8(publicKey)
}

func decodeJWT(jwt string) *JWTSigned {
	params := strings.Split(jwt, JWT_DELIMITER)
	header, ok := decodeJson(params[0]).(*JWTHeader)
	if !ok {
		tlog.Errorf("decodeJWTHeader false ,jwt is %s ", jwt)
		return nil
	}
	payload, ok := decodeJson(params[1]).(*JWTPayload)
	if !ok {
		tlog.Errorf("decodeJWTPayload false ,jwt is %s ", jwt)
		return nil
	}
	signature := decodeSig(params[2])
	//if !ok {
	//	tlog.Errorf("decodeJWTSignature false ,jwt is %s ", jwt)
	//	return nil
	//}
	return &JWTSigned{
		Header:    header,
		Payload:   payload,
		Signature: signature,
	}

}
func decodeJson(dst string) (res interface{}) {
	decode, err := base64.URLEncoding.DecodeString(dst)
	if err != nil {
		tlog.Errorf("decodeJson err is %s ,dst is %s ", err.Error(), dst)
		return nil
	}
	err = json.Unmarshal(decode, &res)
	if err != nil {
		tlog.Errorf("decodeJson err is %s ,dst is %s ", err.Error(), dst)
		return nil
	}
	mapRes, ok := res.(map[string]interface{})
	if !ok {
		tlog.Errorf("decodeJson err is %s ,dst is %s ", err.Error(), dst)
		return
	}
	mapLen := len(mapRes)
	switch mapLen {
	case 2:
		// header
		header := &JWTHeader{}
		err := json.Unmarshal(decode, header)
		if err != nil {
			tlog.Errorf("decodeJson err is %s ,dst is %s ", err.Error(), dst)
			return nil
		}
		res = header
	case 6:
		// payload
		payload := &JWTPayload{}
		err := json.Unmarshal(decode, payload)
		if err != nil {
			tlog.Errorf("decodeJson err is %s ,dst is %s ", err.Error(), dst)
			return nil
		}
		res = payload
	default:
		tlog.Errorf("decodeJson length err len is %d  ,dst is %s ", mapLen, dst)
		return nil
	}
	return
}

func decodeData(data []uint8) (*JWTHeader, *JWTPayload) {
	_, res, err := multibase.Decode(string(data))
	if err != nil {
		tlog.Errorf("decodeData err is %s ,data is %s ", err.Error(), data)
		return nil, nil
	}
	params := strings.Split(string(res), JWT_DELIMITER)
	if len(params) != 2 {
		tlog.Errorf("decodeData err is %s ,data is %s ", err.Error(), data)
		return nil, nil
	}
	header, ok := decodeJson(params[0]).(*JWTHeader)
	if !ok {
		tlog.Errorf("decodeJWTHeader false ,data is %s ", data)
		return nil, nil
	}
	payload, ok := decodeJson(params[1]).(*JWTPayload)
	if !ok {
		tlog.Errorf("decodeJWTPayload false ,data is %s ", data)
		return nil, nil
	}
	return header, payload
}
func decodeSig(data string) []uint8 {
	res, err := base64.URLEncoding.DecodeString(data)
	if err != nil {
		tlog.Errorf("decodeSig err is %s ,data is %s ", err.Error(), data)
		return nil
	}
	return res
}

func encodeSig(sign []uint8) string {
	return base64.URLEncoding.EncodeToString(sign)
}

func encodeJson(val interface{}) string {
	// 默认按照 UTF-8 序列化
	res, err := json.Marshal(val)
	if err != nil {
		tlog.Errorf("encodeJson err is %s ,val is %+v ", err.Error(), val)
		return ""
	}
	data, err := multibase.Encode(multibase.Base64url, res)
	//tlog.Infof("encodeJson data is %s len is %d ", data, len(data))
	if err != nil {
		tlog.Errorf("encodeJson err is %s ,val is %+v ", err.Error(), val)
		return ""
	}
	return data[1:]
}
func encodeData(header *JWTHeader, payload *JWTPayload) []uint8 {
	headerStr := encodeJson(header)
	payloadStr := encodeJson(payload)
	temp := strings.Join([]string{headerStr, payloadStr}, JWT_DELIMITER)
	//multibase.Encode(multibase.Base64URL, []byte(temp))
	//	TODO  DATA_ENCODING 是 utf8 怎么编写
	return []uint8(temp)
}
func encodeIss(publicKey ed25519.PublicKey) string {
	//keyType, err := multibase.Encode(multibase.Base58BTC, []byte("MULTICODEC_ED25519_KEY_TYPE"))
	//if err != nil {
	//	tlog.Error("encodeIss keyType err is %s ", err.Error())
	//	return ""
	//}
	header, err := multibase.Encode(multibase.Base58BTC, []byte(MULTICODEC_ED25519_HEADER))
	if err != nil {
		tlog.Errorf("encodeIss header err is %s ", err.Error())
		return ""
	}
	header = header[1:]
	tempHeader := []byte(header)
	tempHeader = append(tempHeader, publicKey...)
	//tempJoin := bytes.Join([][]byte{[]byte(header), publicKey}, []byte{})
	temp := string(multibase.Base58BTC) + string(tempHeader)
	_, multicodec, err := multibase.Decode(temp)
	if err != nil {
		tlog.Errorf("encodeIss multibase err is %s ", err.Error())
		return ""
	}
	tlog.Infof("encodeIss multicodec is %s ", multicodec)
	//multibase.Decode(multicodec)
	multibase := string(MULTIBASE_BASE58BTC_PREFIX) + string(multicodec)
	return strings.Join([]string{DID_PREFIX, DID_METHOD, multibase}, DID_DELIMITER)
}

func encodeJWT(header *JWTHeader, payload *JWTPayload, signature []uint8) string {
	headerStr := encodeJson(header)
	payloadStr := encodeJson(payload)
	signatureStr := encodeSig(signature)
	return strings.Join([]string{headerStr, payloadStr, signatureStr}, JWT_DELIMITER)

}
