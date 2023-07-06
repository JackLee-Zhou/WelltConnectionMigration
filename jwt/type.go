package jwt

type JWTHeader struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
}

//type  struct {
//	Iss string `json:"iss"`
//	Sub string `json:"sub"`
//}

type JWTSigned struct {
	Header    *JWTHeader  `json:"header"`
	Payload   *JWTPayload `json:"payload"`
	Signature []uint8     `json:"signature"`
}

type JWTPayload struct {
	Iat int64  `json:"iat"`
	Exp int64  `json:"exp"`
	Iss string `json:"iss"`
	Sub string `json:"sub"`
	Aud string `json:"aud"`
	Act string `json:"act"`
}

//func (header JWTHeader) MarshalJSON() ([]byte, error) {
//	return json.Marshal(map[string]interface{}{
//		"alg": header.alg,
//		"typ": header.typ,
//	})
//}
//
//func (payload JWTPayload) MarshalJSON() ([]byte, error) {
//	return json.Marshal(map[string]interface{}{
//		"iat": payload.iat,
//		"exp": payload.exp,
//		"iss": payload.iss,
//		"sub": payload.sub,
//		"aud": payload.aud,
//		"act": payload.act,
//	})
//}
//
//func (sign JWTSigned) MarshalJSON() ([]byte, error) {
//	return json.Marshal(map[string]interface{}{
//		"header":    sign.Header,
//		"payload":   sign.Payload,
//		"signature": sign.Signature,
//	})
//}
