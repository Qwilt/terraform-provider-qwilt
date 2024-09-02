package cdn

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"
)

type SelfSignedCertGenerator struct {
	Crt string
	PK  string
}

func NewSelfSignedCertGenerator() *SelfSignedCertGenerator {
	g := SelfSignedCertGenerator{}
	return &g
}

func (g *SelfSignedCertGenerator) generate(domain string) {

	// Generate a new RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		panic(err)
	}

	// Create a template for the certificate
	certTemplate := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   domain,
			Organization: []string{"My Company"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(365 * 24 * time.Hour), // 1 year validity

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},                   // Include the domain in the certificate
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")}, // Optional: Add IP addresses
	}

	// Create the certificate
	certBytes, err := x509.CreateCertificate(rand.Reader, certTemplate, certTemplate, &privateKey.PublicKey, privateKey)
	if err != nil {
		panic(err)
	}

	// Encode the private key to PEM format and store in a string
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	if privateKeyPEM == nil {
		panic("failed to encode private key to PEM format")
	}
	g.PK = string(privateKeyPEM) // Convert the byte slice to a string

	// Encode the certificate to PEM format and store in a string
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certBytes,
	})
	if certPEM == nil {
		panic("failed to encode crt to PEM format")
	}
	g.Crt = string(certPEM) // Convert the byte slice to a string
}
