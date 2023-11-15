package hash

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

var (
	errInvalidFormat       = errors.New("the encoded hash is not in the correct format")
	errIncompatibleVersion = errors.New("incompatible version of argon2")
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, encodedHash string) (bool, error)
}

type Argon2Hasher struct {
	parameters parameters
}

type parameters struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

func NewArgon2Hasher(memory, iterations, saltLength, keyLength uint32, parallelism uint8) *Argon2Hasher {
	return &Argon2Hasher{
		parameters: parameters{
			memory:      memory,
			iterations:  iterations,
			parallelism: parallelism,
			saltLength:  saltLength,
			keyLength:   keyLength,
		},
	}
}

func (h *Argon2Hasher) Hash(password string) (encodedHash string, err error) {
	salt, err := generateRandomBytes(h.parameters.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, h.parameters.iterations, h.parameters.memory,
		h.parameters.parallelism, h.parameters.keyLength)

	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	encodedHash = fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.parameters.memory,
		h.parameters.iterations, h.parameters.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func (h *Argon2Hasher) Verify(password, encodedHash string) (match bool, err error) {
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *parameters, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 { //nolint:gomnd
		return nil, nil, nil, errInvalidFormat
	}

	var version int
	if _, err = fmt.Sscanf(vals[2], "v=%d", &version); err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, errIncompatibleVersion
	}

	p = &parameters{}
	if _, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism); err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}

	return b, nil
}
