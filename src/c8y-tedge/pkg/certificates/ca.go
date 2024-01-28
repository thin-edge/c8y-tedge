package certificates

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"math/big"
	"net"
	"os"
	"time"
)

type Certificate struct {
	Public *bytes.Buffer
	Key    *bytes.Buffer
	Cert   *x509.Certificate
}

func NewCertificateAuthority() (Certificate, error) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		BasicConstraintsValid: true,
	}

	cert := Certificate{
		Cert: ca,
	}

	caPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return cert, err
	}

	caBytes, err := x509.CreateCertificate(rand.Reader, ca, ca, &caPrivKey.PublicKey, caPrivKey)
	if err != nil {
		return Certificate{}, err
	}

	caPEM := new(bytes.Buffer)
	pem.Encode(caPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: caBytes,
	})

	caPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(caPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(caPrivKey),
	})
	cert.Key = caPrivKeyPEM
	cert.Public = caPEM
	return cert, nil
}

func NewLeafCertificate(parent *Certificate) (*Certificate, error) {
	// set up our server certificate
	leafCert := &x509.Certificate{
		SerialNumber: big.NewInt(2019),
		Subject: pkix.Name{
			Organization:  []string{"Company, INC."},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{"Golden Gate Bridge"},
			PostalCode:    []string{"94016"},
		},
		IPAddresses:  []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(2, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature,
	}

	certPrivKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, err
	}

	certBytes, err := x509.CreateCertificate(rand.Reader, leafCert, parent.Cert, &certPrivKey.PublicKey, parent.Key)
	if err != nil {
		return nil, err
	}

	certPEM := new(bytes.Buffer)
	pem.Encode(certPEM, &pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})

	certPrivKeyPEM := new(bytes.Buffer)
	pem.Encode(certPrivKeyPEM, &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(certPrivKey),
	})

	cert := &Certificate{}
	return cert, nil
}

func LoadCertificateFromFiles(priv string, public string) (*Certificate, error) {
	cert := &Certificate{}

	// public
	pubKeyFile, err := os.ReadFile(public)
	if err != nil {
		return nil, err
	}

	pubPEMBlock, _ := pem.Decode(pubKeyFile)
	if pubPEMBlock == nil {
		return nil, fmt.Errorf("failed to parse public key")
	}
	pubCert, err := x509.ParseCertificate(pubPEMBlock.Bytes)
	if err != nil {
		return nil, err
	}
	cert.Cert = pubCert

	// private
	privKeyFile, err := os.ReadFile(public)
	if err != nil {
		return nil, err
	}
	privPEMBlock, _ := pem.Decode(privKeyFile)
	if privPEMBlock == nil {
		return nil, fmt.Errorf("failed to parse public key")
	}
	cert.Public = bytes.NewBuffer(pubPEMBlock.Bytes)
	cert.Key = bytes.NewBuffer(privPEMBlock.Bytes)
	return cert, nil
}

func SignCertificateRequestFile(csrPath string, parentPublicKey string, parentPrivateKey string, cert io.Writer) error {
	// load client certificate request
	clientCSRFile, err := os.ReadFile(csrPath)
	if err != nil {
		return err
	}
	pemBlock, _ := pem.Decode(clientCSRFile)
	if pemBlock == nil {
		panic("pem.Decode failed")
	}
	clientCSR, err := x509.ParseCertificateRequest(pemBlock.Bytes)
	if err != nil {
		return err
	}

	ca, err := LoadCertificateFromFiles(parentPrivateKey, parentPublicKey)
	if err != nil {
		return err
	}

	return SignCSR(clientCSR, ca, cert)
}

func SignCSR(clientCSR *x509.CertificateRequest, parent *Certificate, out io.Writer) error {
	if err := clientCSR.CheckSignature(); err != nil {
		return err
	}

	// create client certificate template
	clientCRTTemplate := x509.Certificate{
		Signature:          clientCSR.Signature,
		SignatureAlgorithm: clientCSR.SignatureAlgorithm,

		PublicKeyAlgorithm: clientCSR.PublicKeyAlgorithm,
		PublicKey:          clientCSR.PublicKey,

		SerialNumber: big.NewInt(2),
		Issuer:       parent.Cert.Subject,
		Subject:      clientCSR.Subject,
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(1, 0, 0),
		KeyUsage:     x509.KeyUsageDigitalSignature,
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	}

	// create client certificate from template and CA public key
	clientCRTRaw, err := x509.CreateCertificate(rand.Reader, &clientCRTTemplate, parent.Cert, clientCSR.PublicKey, parent.Key)
	if err != nil {
		return err
	}

	// save the certificate
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: clientCRTRaw})

	// Append parent cert
	pem.Encode(out, &pem.Block{Type: "CERTIFICATE", Bytes: parent.Cert.Raw})
	return nil
}
