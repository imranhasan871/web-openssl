package openssl

import "time"

type KeyType string
type CertificateFormat string
type KeyFormat string
type HashAlgorithm string
type EncryptionAlgorithm string

const (
	// Key types
	KeyTypeRSA   KeyType = "rsa"
	KeyTypeEC    KeyType = "ec"
	KeyTypeED25519 KeyType = "ed25519"
	KeyTypeDSA   KeyType = "dsa"

	// Certificate formats
	CertFormatPEM   CertificateFormat = "pem"
	CertFormatDER   CertificateFormat = "der"
	CertFormatPKCS12 CertificateFormat = "pkcs12"

	// Key formats
	KeyFormatPEM    KeyFormat = "pem"
	KeyFormatDER    KeyFormat = "der"
	KeyFormatPKCS8  KeyFormat = "pkcs8"

	// Hash algorithms
	HashMD5     HashAlgorithm = "md5"
	HashSHA1    HashAlgorithm = "sha1"
	HashSHA224  HashAlgorithm = "sha224"
	HashSHA256  HashAlgorithm = "sha256"
	HashSHA384  HashAlgorithm = "sha384"
	HashSHA512  HashAlgorithm = "sha512"
	HashBLAKE2B HashAlgorithm = "blake2b"
	HashBLAKE2S HashAlgorithm = "blake2s"

	// Encryption algorithms
	EncryptAES128 EncryptionAlgorithm = "aes-128-cbc"
	EncryptAES192 EncryptionAlgorithm = "aes-192-cbc"
	EncryptAES256 EncryptionAlgorithm = "aes-256-cbc"
	EncryptDES3   EncryptionAlgorithm = "des-ede3-cbc"
	EncryptChaCha20 EncryptionAlgorithm = "chacha20"
)

// Request/Response types for various operations

type GenerateKeyRequest struct {
	KeyType   KeyType `json:"keyType" binding:"required"`
	KeySize   int     `json:"keySize,omitempty"`
	Curve     string  `json:"curve,omitempty"`
	Password  string  `json:"password,omitempty"`
	Format    KeyFormat `json:"format,omitempty"`
}

type GenerateKeyResponse struct {
	PrivateKey string `json:"privateKey"`
	PublicKey  string `json:"publicKey"`
	Format     string `json:"format"`
}

type GenerateCertificateRequest struct {
	KeyType        KeyType   `json:"keyType" binding:"required"`
	KeySize        int       `json:"keySize,omitempty"`
	Subject        Subject   `json:"subject" binding:"required"`
	ValidDays      int       `json:"validDays" binding:"required"`
	IsCA           bool      `json:"isCA,omitempty"`
	KeyUsage       []string  `json:"keyUsage,omitempty"`
	ExtKeyUsage    []string  `json:"extKeyUsage,omitempty"`
	SANs           []string  `json:"sans,omitempty"`
	HashAlgorithm  HashAlgorithm `json:"hashAlgorithm,omitempty"`
	Format         CertificateFormat `json:"format,omitempty"`
}

type Subject struct {
	CommonName         string `json:"commonName" binding:"required"`
	Country            string `json:"country,omitempty"`
	State              string `json:"state,omitempty"`
	Locality           string `json:"locality,omitempty"`
	Organization       string `json:"organization,omitempty"`
	OrganizationalUnit string `json:"organizationalUnit,omitempty"`
	EmailAddress       string `json:"emailAddress,omitempty"`
}

type GenerateCertificateResponse struct {
	Certificate string `json:"certificate"`
	PrivateKey  string `json:"privateKey"`
	PublicKey   string `json:"publicKey"`
	Format      string `json:"format"`
}

type GenerateCSRRequest struct {
	KeyType       KeyType `json:"keyType" binding:"required"`
	KeySize       int     `json:"keySize,omitempty"`
	Subject       Subject `json:"subject" binding:"required"`
	SANs          []string `json:"sans,omitempty"`
	HashAlgorithm HashAlgorithm `json:"hashAlgorithm,omitempty"`
}

type GenerateCSRResponse struct {
	CSR        string `json:"csr"`
	PrivateKey string `json:"privateKey"`
}

type ParseCertificateRequest struct {
	Certificate string            `json:"certificate" binding:"required"`
	Format      CertificateFormat `json:"format,omitempty"`
}

type CertificateInfo struct {
	Subject          Subject           `json:"subject"`
	Issuer           Subject           `json:"issuer"`
	SerialNumber     string            `json:"serialNumber"`
	NotBefore        time.Time         `json:"notBefore"`
	NotAfter         time.Time         `json:"notAfter"`
	IsCA             bool              `json:"isCA"`
	KeyUsage         []string          `json:"keyUsage"`
	ExtKeyUsage      []string          `json:"extKeyUsage"`
	SANs             []string          `json:"sans"`
	SignatureAlgorithm string          `json:"signatureAlgorithm"`
	PublicKeyAlgorithm string          `json:"publicKeyAlgorithm"`
	PublicKeySize    int               `json:"publicKeySize"`
	Version          int               `json:"version"`
	Fingerprints     map[string]string `json:"fingerprints"`
	Extensions       map[string]string `json:"extensions"`
	IsExpired        bool              `json:"isExpired"`
	DaysUntilExpiry  int               `json:"daysUntilExpiry"`
}

type VerifyCertificateRequest struct {
	Certificate string `json:"certificate" binding:"required"`
	CAChain     string `json:"caChain,omitempty"`
	CRLFile     string `json:"crlFile,omitempty"`
}

type VerifyCertificateResponse struct {
	IsValid      bool     `json:"isValid"`
	ErrorMessage string   `json:"errorMessage,omitempty"`
	Chain        []string `json:"chain,omitempty"`
}

type EncryptRequest struct {
	Data      string              `json:"data" binding:"required"`
	Algorithm EncryptionAlgorithm `json:"algorithm" binding:"required"`
	Key       string              `json:"key,omitempty"`
	Password  string              `json:"password,omitempty"`
	PublicKey string              `json:"publicKey,omitempty"`
}

type EncryptResponse struct {
	EncryptedData string `json:"encryptedData"`
	Key           string `json:"key,omitempty"`
	IV            string `json:"iv,omitempty"`
}

type DecryptRequest struct {
	EncryptedData string              `json:"encryptedData" binding:"required"`
	Algorithm     EncryptionAlgorithm `json:"algorithm" binding:"required"`
	Key           string              `json:"key,omitempty"`
	Password      string              `json:"password,omitempty"`
	PrivateKey    string              `json:"privateKey,omitempty"`
	IV            string              `json:"iv,omitempty"`
}

type DecryptResponse struct {
	DecryptedData string `json:"decryptedData"`
}

type HashRequest struct {
	Data      string        `json:"data" binding:"required"`
	Algorithm HashAlgorithm `json:"algorithm" binding:"required"`
	Key       string        `json:"key,omitempty"`
}

type HashResponse struct {
	Hash      string `json:"hash"`
	Algorithm string `json:"algorithm"`
}

type SSLTestRequest struct {
	Hostname string `json:"hostname" binding:"required"`
	Port     int    `json:"port,omitempty"`
	Protocol string `json:"protocol,omitempty"`
	Timeout  int    `json:"timeout,omitempty"`
}

type SSLTestResponse struct {
	IsValid        bool                   `json:"isValid"`
	Certificate    *CertificateInfo       `json:"certificate"`
	Chain          []*CertificateInfo     `json:"chain"`
	Protocol       string                 `json:"protocol"`
	Cipher         string                 `json:"cipher"`
	Vulnerabilities []string              `json:"vulnerabilities"`
	Grade          string                 `json:"grade"`
	Details        map[string]interface{} `json:"details"`
}