package certificate

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"time"
)

type CertificatePEM struct {
	CertPEM []byte
	KeyPEM  []byte
}

func GenerateSelfSignedCertificate(validFor time.Duration, domain string) (*CertificatePEM, error) {
	// Generate a new private key
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}

	// Set certificate properties
	notBefore := time.Now()
	notAfter := notBefore.Add(validFor) // 1-year validity

	// Create a self-signed certificate template
	serialNumber, err := rand.Int(rand.Reader, big.NewInt(1<<62)) // Generate a random serial number
	if err != nil {
		return nil, err
	}

	certTemplate := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"shiwam-cert-controller"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              []string{domain},
	}

	// Create the self-signed certificate
	certDERBytes, err := x509.CreateCertificate(rand.Reader, &certTemplate, &certTemplate, &priv.PublicKey, priv)
	if err != nil {
		return nil, err
	}

	// PEM encode the private key
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(priv)})

	// PEM encode the certificate
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certDERBytes})

	return &CertificatePEM{CertPEM: certPEM, KeyPEM: keyPEM}, nil
}
