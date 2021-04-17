package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"time"
)

const (
	KeySize = 2048
)

func main() {
	var err error

	var caPvtKey *rsa.PrivateKey
	caPvtKey, err = rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	var caCert = &x509.Certificate{
		SerialNumber: big.NewInt(2021),
		Subject: pkix.Name{
			Country:      []string{"SG"},
			Organization: []string{"John L. Lao"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		IsCA:                  true,
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	var caCertBytes []byte
	caCertBytes, err = x509.CreateCertificate(rand.Reader, caCert, caCert, &caPvtKey.PublicKey, caPvtKey)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SavePrivateKey(caPvtKey, "./ca.key")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SaveCertificate(caCertBytes, "./ca.pem")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}

	// ------- certificate -------
	var cert1PvtKey *rsa.PrivateKey
	cert1PvtKey, err = rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	var cert1 = &x509.Certificate{
		SerialNumber: big.NewInt(2022),
		Subject: pkix.Name{
			Country:      []string{"SG"},
			Organization: []string{"John L. Lao"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	var cert1Bytes []byte
	cert1Bytes, err = x509.CreateCertificate(rand.Reader, cert1, caCert, &cert1PvtKey.PublicKey, caPvtKey)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SavePrivateKey(cert1PvtKey, "./cert1.key")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SaveCertificate(cert1Bytes, "./cert1.pem")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}

	// ------- certificate ------
	var cert2PvtKey *rsa.PrivateKey
	cert2PvtKey, err = rsa.GenerateKey(rand.Reader, KeySize)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	var cert2 = &x509.Certificate{
		SerialNumber: big.NewInt(2023),
		Subject: pkix.Name{
			Country:      []string{"SG"},
			Organization: []string{"John L. Lao"},
			CommonName:   "localhost",
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	var cert2Bytes []byte
	cert2Bytes, err = x509.CreateCertificate(rand.Reader, cert2, caCert, &cert2PvtKey.PublicKey, caPvtKey)
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SavePrivateKey(cert2PvtKey, "./cert2.key")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
	err = SaveCertificate(cert2Bytes, "./cert2.pem")
	if err != nil {
		fmt.Printf("[main] %s \n", err.Error())
		os.Exit(1)
	}
}

func SavePrivateKey(k *rsa.PrivateKey, p string) error {
	var err error

	var f *os.File
	f, err = os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(k),
	})
}

func SaveCertificate(c []byte, p string) error {
	var err error
	var f *os.File
	f, err = os.Create(p)
	if err != nil {
		return err
	}
	defer f.Close()

	return pem.Encode(f, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: c,
	})
}
