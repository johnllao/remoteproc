package security

import (
	"bufio"
	"errors"
	"io"

	"github.com/johnllao/remoteproc/pkg/hmac"
)

func ClientHandshake(reader io.Reader, writer io.Writer, token string) error {

	var err error

	// start the handshake
	var r = bufio.NewReader(reader)
	var w = bufio.NewWriter(writer)

	var tokenRequest = "TOKEN:" + token
	_, err = w.WriteString(tokenRequest + "\n")
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}
	var tokenReply string
	tokenReply, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	if tokenReply[:2] != "OK" {
		return errors.New("invalid token")
	}
	// end of the handshake

	return nil
}

func ServerHandshake(reader io.Reader, writer io.Writer, key string) error {
	var err error

	// start the handshake
	var tokenRequest string
	var tokenResponse = "INVALID_TOKEN"

	var r = bufio.NewReader(reader)
	var w = bufio.NewWriter(writer)

	var tokenErr error

	tokenRequest, err = r.ReadString('\n')
	if err != nil {
		return err
	}
	if tokenRequest[:5] == "TOKEN" {
		var token = tokenRequest[6 : len(tokenRequest)-1]

		var validToken bool
		if validToken, tokenErr = hmac.VerifyToken(key, token); validToken && tokenErr == nil {
			tokenResponse = "OK"
		}
	}

	_, err = w.WriteString(tokenResponse + "\n")
	if err != nil {
		return err
	}
	err = w.Flush()
	if err != nil {
		return err
	}

	if tokenErr != nil {
		return tokenErr
	}

	// end of the handshake

	return nil
}
