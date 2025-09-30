package openssl

import (
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Service struct {
	opensslPath string
}

func NewService(opensslPath string) *Service {
	return &Service{
		opensslPath: opensslPath,
	}
}

func (s *Service) GenerateKey(req *GenerateKeyRequest) (*GenerateKeyResponse, error) {
	var cmd *exec.Cmd

	switch req.KeyType {
	case KeyTypeRSA:
		keySize := req.KeySize
		if keySize == 0 {
			keySize = 2048
		}
		cmd = exec.Command(s.opensslPath, "genpkey", "-algorithm", "RSA", "-pkcs8", "-out", "-", "-outform", "PEM")
		cmd.Args = append(cmd.Args, "-pkeyopt", fmt.Sprintf("rsa_keygen_bits:%d", keySize))

	case KeyTypeEC:
		curve := req.Curve
		if curve == "" {
			curve = "secp256r1"
		}
		cmd = exec.Command(s.opensslPath, "genpkey", "-algorithm", "EC", "-pkcs8", "-out", "-", "-outform", "PEM")
		cmd.Args = append(cmd.Args, "-pkeyopt", fmt.Sprintf("ec_paramgen_curve:%s", curve))

	case KeyTypeED25519:
		cmd = exec.Command(s.opensslPath, "genpkey", "-algorithm", "Ed25519", "-out", "-", "-outform", "PEM")

	default:
		return nil, fmt.Errorf("unsupported key type: %s", req.KeyType)
	}

	if req.Password != "" {
		cmd.Args = append(cmd.Args, "-aes256", "-pass", fmt.Sprintf("pass:%s", req.Password))
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("openssl error: %s", stderr.String())
	}

	privateKey := stdout.String()

	// Generate public key
	publicKeyCmd := exec.Command(s.opensslPath, "pkey", "-in", "-", "-pubout", "-outform", "PEM")
	publicKeyCmd.Stdin = strings.NewReader(privateKey)

	var pubOut, pubErr bytes.Buffer
	publicKeyCmd.Stdout = &pubOut
	publicKeyCmd.Stderr = &pubErr

	if err := publicKeyCmd.Run(); err != nil {
		return nil, fmt.Errorf("public key generation error: %s", pubErr.String())
	}

	return &GenerateKeyResponse{
		PrivateKey: privateKey,
		PublicKey:  pubOut.String(),
		Format:     "pem",
	}, nil
}

func (s *Service) GenerateCertificate(req *GenerateCertificateRequest) (*GenerateCertificateResponse, error) {
	// First generate a key
	keyReq := &GenerateKeyRequest{
		KeyType: req.KeyType,
		KeySize: req.KeySize,
		Format:  KeyFormatPEM,
	}
	keyResp, err := s.GenerateKey(keyReq)
	if err != nil {
		return nil, fmt.Errorf("key generation failed: %w", err)
	}

	// Create subject string
	subject := s.buildSubjectString(req.Subject)

	// Build OpenSSL command for certificate generation
	args := []string{
		"req", "-new", "-x509",
		"-key", "-",
		"-out", "-",
		"-subj", subject,
		"-days", strconv.Itoa(req.ValidDays),
	}

	if req.HashAlgorithm != "" {
		args = append(args, "-"+string(req.HashAlgorithm))
	} else {
		args = append(args, "-sha256")
	}

	// Add extensions if specified
	if len(req.SANs) > 0 || req.IsCA || len(req.KeyUsage) > 0 {
		config := s.buildConfigFile(req)
		args = append(args, "-config", "-")

		cmd := exec.Command(s.opensslPath, args...)
		cmd.Stdin = strings.NewReader(keyResp.PrivateKey + "\n" + config)

		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("certificate generation error: %s", stderr.String())
		}

		return &GenerateCertificateResponse{
			Certificate: stdout.String(),
			PrivateKey:  keyResp.PrivateKey,
			PublicKey:   keyResp.PublicKey,
			Format:      "pem",
		}, nil
	}

	cmd := exec.Command(s.opensslPath, args...)
	cmd.Stdin = strings.NewReader(keyResp.PrivateKey)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("certificate generation error: %s", stderr.String())
	}

	return &GenerateCertificateResponse{
		Certificate: stdout.String(),
		PrivateKey:  keyResp.PrivateKey,
		PublicKey:   keyResp.PublicKey,
		Format:      "pem",
	}, nil
}

func (s *Service) ParseCertificate(req *ParseCertificateRequest) (*CertificateInfo, error) {
	cmd := exec.Command(s.opensslPath, "x509", "-in", "-", "-noout", "-text")
	cmd.Stdin = strings.NewReader(req.Certificate)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("certificate parsing error: %s", stderr.String())
	}

	return s.parseCertificateOutput(stdout.String(), req.Certificate)
}

func (s *Service) VerifyCertificate(req *VerifyCertificateRequest) (*VerifyCertificateResponse, error) {
	args := []string{"verify"}

	if req.CAChain != "" {
		args = append(args, "-CAfile", "-")
	}

	args = append(args, "-")

	cmd := exec.Command(s.opensslPath, args...)

	var input strings.Builder
	if req.CAChain != "" {
		input.WriteString(req.CAChain)
		input.WriteString("\n")
	}
	input.WriteString(req.Certificate)

	cmd.Stdin = strings.NewReader(input.String())

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stdout.String()

	if err != nil {
		return &VerifyCertificateResponse{
			IsValid:      false,
			ErrorMessage: stderr.String(),
		}, nil
	}

	isValid := strings.Contains(output, "OK")

	return &VerifyCertificateResponse{
		IsValid: isValid,
	}, nil
}

func (s *Service) SymmetricEncrypt(req *EncryptRequest) (*EncryptResponse, error) {
	// Generate random key if not provided
	key := req.Key
	if key == "" && req.Password == "" {
		keyBytes := make([]byte, 32) // 256-bit key
		if _, err := rand.Read(keyBytes); err != nil {
			return nil, fmt.Errorf("key generation failed: %w", err)
		}
		key = base64.StdEncoding.EncodeToString(keyBytes)
	}

	args := []string{"enc", "-" + string(req.Algorithm), "-base64"}

	if key != "" {
		args = append(args, "-K", key, "-iv", "0")
	} else {
		args = append(args, "-pass", fmt.Sprintf("pass:%s", req.Password))
	}

	cmd := exec.Command(s.opensslPath, args...)
	cmd.Stdin = strings.NewReader(req.Data)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("encryption error: %s", stderr.String())
	}

	return &EncryptResponse{
		EncryptedData: strings.TrimSpace(stdout.String()),
		Key:           key,
	}, nil
}

func (s *Service) Decrypt(req *DecryptRequest) (*DecryptResponse, error) {
	args := []string{"enc", "-d", "-" + string(req.Algorithm), "-base64"}

	if req.Key != "" {
		args = append(args, "-K", req.Key)
		if req.IV != "" {
			args = append(args, "-iv", req.IV)
		} else {
			args = append(args, "-iv", "0")
		}
	} else if req.Password != "" {
		args = append(args, "-pass", fmt.Sprintf("pass:%s", req.Password))
	}

	cmd := exec.Command(s.opensslPath, args...)
	cmd.Stdin = strings.NewReader(req.EncryptedData)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("decryption error: %s", stderr.String())
	}

	return &DecryptResponse{
		DecryptedData: stdout.String(),
	}, nil
}

func (s *Service) GenerateHash(req *HashRequest) (*HashResponse, error) {
	var args []string

	if req.Key != "" {
		// HMAC
		args = []string{"dgst", "-" + string(req.Algorithm), "-hmac", req.Key}
	} else {
		// Regular hash
		args = []string{"dgst", "-" + string(req.Algorithm)}
	}

	cmd := exec.Command(s.opensslPath, args...)
	cmd.Stdin = strings.NewReader(req.Data)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("hash generation error: %s", stderr.String())
	}

	// Extract hash from output (format: "algorithm(stdin)= hash")
	output := strings.TrimSpace(stdout.String())
	parts := strings.Split(output, "= ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("unexpected hash output format")
	}

	return &HashResponse{
		Hash:      parts[1],
		Algorithm: string(req.Algorithm),
	}, nil
}

// Helper functions

func (s *Service) buildSubjectString(subject Subject) string {
	var parts []string

	if subject.Country != "" {
		parts = append(parts, fmt.Sprintf("C=%s", subject.Country))
	}
	if subject.State != "" {
		parts = append(parts, fmt.Sprintf("ST=%s", subject.State))
	}
	if subject.Locality != "" {
		parts = append(parts, fmt.Sprintf("L=%s", subject.Locality))
	}
	if subject.Organization != "" {
		parts = append(parts, fmt.Sprintf("O=%s", subject.Organization))
	}
	if subject.OrganizationalUnit != "" {
		parts = append(parts, fmt.Sprintf("OU=%s", subject.OrganizationalUnit))
	}
	if subject.CommonName != "" {
		parts = append(parts, fmt.Sprintf("CN=%s", subject.CommonName))
	}
	if subject.EmailAddress != "" {
		parts = append(parts, fmt.Sprintf("emailAddress=%s", subject.EmailAddress))
	}

	return "/" + strings.Join(parts, "/")
}

func (s *Service) buildConfigFile(req *GenerateCertificateRequest) string {
	var config strings.Builder

	config.WriteString("[req]\n")
	config.WriteString("distinguished_name = req_distinguished_name\n")
	config.WriteString("req_extensions = v3_req\n")
	config.WriteString("[req_distinguished_name]\n")
	config.WriteString("[v3_req]\n")

	if req.IsCA {
		config.WriteString("basicConstraints = CA:TRUE\n")
	}

	if len(req.KeyUsage) > 0 {
		config.WriteString("keyUsage = " + strings.Join(req.KeyUsage, ",") + "\n")
	}

	if len(req.ExtKeyUsage) > 0 {
		config.WriteString("extendedKeyUsage = " + strings.Join(req.ExtKeyUsage, ",") + "\n")
	}

	if len(req.SANs) > 0 {
		config.WriteString("subjectAltName = @alt_names\n")
		config.WriteString("[alt_names]\n")
		for i, san := range req.SANs {
			if strings.Contains(san, "@") {
				config.WriteString(fmt.Sprintf("email.%d = %s\n", i+1, san))
			} else {
				config.WriteString(fmt.Sprintf("DNS.%d = %s\n", i+1, san))
			}
		}
	}

	return config.String()
}

func (s *Service) parseCertificateOutput(output, certPEM string) (*CertificateInfo, error) {
	info := &CertificateInfo{
		Fingerprints: make(map[string]string),
		Extensions:   make(map[string]string),
	}

	// Parse the certificate using Go's x509 package for accurate parsing
	block, _ := pem.Decode([]byte(certPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse certificate: %w", err)
	}

	// Fill basic info
	info.Subject = Subject{
		CommonName:         cert.Subject.CommonName,
		Country:            strings.Join(cert.Subject.Country, ","),
		State:              strings.Join(cert.Subject.Province, ","),
		Locality:           strings.Join(cert.Subject.Locality, ","),
		Organization:       strings.Join(cert.Subject.Organization, ","),
		OrganizationalUnit: strings.Join(cert.Subject.OrganizationalUnit, ","),
	}

	info.Issuer = Subject{
		CommonName:         cert.Issuer.CommonName,
		Country:            strings.Join(cert.Issuer.Country, ","),
		State:              strings.Join(cert.Issuer.Province, ","),
		Locality:           strings.Join(cert.Issuer.Locality, ","),
		Organization:       strings.Join(cert.Issuer.Organization, ","),
		OrganizationalUnit: strings.Join(cert.Issuer.OrganizationalUnit, ","),
	}

	info.SerialNumber = cert.SerialNumber.String()
	info.NotBefore = cert.NotBefore
	info.NotAfter = cert.NotAfter
	info.Version = cert.Version
	info.SignatureAlgorithm = cert.SignatureAlgorithm.String()

	info.IsExpired = time.Now().After(cert.NotAfter)
	if !info.IsExpired {
		info.DaysUntilExpiry = int(time.Until(cert.NotAfter).Hours() / 24)
	}

	// Parse SANs
	for _, san := range cert.DNSNames {
		info.SANs = append(info.SANs, san)
	}
	for _, san := range cert.EmailAddresses {
		info.SANs = append(info.SANs, san)
	}

	return info, nil
}