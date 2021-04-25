package hmac

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"time"
)

const (
	SignatureSize = 32
)

type Claim struct {
	Name   string
	Expiry int64
}

func ToClaim(data []byte) (*Claim, error) {
	var err error
	var c = new(Claim)
	err = gob.NewDecoder(bytes.NewReader(data)).Decode(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (t *Claim) Bytes() ([]byte, error) {
	var err error
	var b = &bytes.Buffer{}
	err = gob.NewEncoder(b).Encode(t)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func Sign(key string, data []byte) ([]byte, error) {
	var err error
	var h = hmac.New(sha256.New, []byte(key))
	_, err = h.Write(data)
	if err != nil {
		return nil, err
	}
	return h.Sum(nil), nil
}

func Verify(key string, signature, data []byte) (bool, error) {
	var err error
	var h = hmac.New(sha256.New, []byte(key))
	_, err = h.Write(data)
	if err != nil {
		return false, err
	}
	var actualSignature = h.Sum(nil)
	var match = hmac.Equal(actualSignature, signature)
	if !match {
		return false, errors.New("invalid signature")
	}
	return match, nil
}

func Token(key string, claim *Claim) (string, error) {
	var err error
	var data []byte
	data, err = claim.Bytes()
	if err != nil {
		return "", err
	}
	var s []byte
	s, err = Sign(key, data)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(append(s, data...)), nil
}

func VerifyToken(key, token string) (bool, error) {
	var err error
	var d []byte
	d, err = base64.StdEncoding.DecodeString(token)
	if err != nil {
		return false, err
	}
	var valid bool
	valid, err = Verify(key, d[:SignatureSize], d[SignatureSize:])
	if !valid {
		return false, err
	}
	var c *Claim
	c, err = ToClaim(d[SignatureSize:])
	if err != nil {
		return false, err
	}
	if time.Now().UnixNano() > c.Expiry {
		return false, errors.New("token is expired")
	}
	return true, nil
}
