package hmac

import (
	"testing"
	"time"
)

func TestSignAndVerify(t *testing.T) {
	var err error

	var key = "secret"

	var signature []byte
	signature, err = Sign(key, []byte("The quick brown fox"))
	if err != nil {
		t.Fatal(err)
	}

	var valid bool
	valid, err = Verify(key, signature, []byte("The quick brown fox"))
	if !valid || err != nil {
		t.Error(err)
	}
	valid, _ = Verify(key, signature, []byte("The quick brown fox."))
	if valid {
		t.Error("expected invalid signature since the data are changed")
	}
	valid, _ = Verify(key, []byte("01234567890123456789012345678901"), []byte("The quick brown fox"))
	if valid {
		t.Error("expected invalid signature")
	}
}

func TestVerifyToken(t *testing.T) {
	var err error
	var valid bool

	var c1 = &Claim{
		Name:   "Owner1",
		Expiry: time.Now().Add(24 * time.Hour).UnixNano(),
	}

	var t1 string
	t1, err = Token("secret1", c1)
	if err != nil {
		t.Fatal(err)
	}
	valid, err = VerifyToken("secret1", t1)
	if !valid || err != nil {
		t.Error(err)
	}
	valid, _ = VerifyToken("secret2", t1)
	if valid {
		t.Error("expecting nvalid token because the key is incorrect")
	}

	var c2 = &Claim{
		Name:   "Owner1",
		Expiry: time.Now().Add(-1 * time.Hour).UnixNano(),
	}
	var t2 string
	t2, err = Token("secret1", c2)
	if err != nil {
		t.Fatal(err)
	}
	valid, _ = VerifyToken("secret1", t2)
	if valid {
		t.Error(err.Error() + " - expecting nvalid token because token is expired")
	}
}
