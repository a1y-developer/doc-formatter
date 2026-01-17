package credentials

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/argon2"
)

type Argon2idHash struct {
	time    uint32
	memory  uint32
	threads uint8
	keyLen  uint32
	saltLen uint32
}

var DefaultArgon2id = Argon2idHash{
	time:    1,
	memory:  64 * 1024,
	threads: uint8(runtime.NumCPU()),
	keyLen:  32,
	saltLen: 16,
}

func NewDefaultArgon2idHash() *Argon2idHash {
	return &DefaultArgon2id
}

func NewArgon2idHash(time, memory uint32, threads uint8, keyLen, saltLen uint32) *Argon2idHash {
	return &Argon2idHash{
		time:    time,
		memory:  memory,
		threads: threads,
		keyLen:  keyLen,
		saltLen: saltLen,
	}
}

// randomSalt generates a random salt with the given length.
func (a *Argon2idHash) randomSalt(length uint32) ([]byte, error) {
	salt := make([]byte, length)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// HashPassword generates a hash of the password using Argon2id.
func (a *Argon2idHash) HashPassword(password string, salt []byte) (string, error) {
	if len(salt) == 0 {
		var err error
		salt, err = a.randomSalt(a.saltLen)
		if err != nil {
			return "", err
		}
	}
	hashPassword := argon2.IDKey([]byte(password), salt, a.time, a.memory, a.threads, a.keyLen)

	base64HashPassword := base64.RawStdEncoding.EncodeToString(hashPassword)
	base64Salt := base64.RawStdEncoding.EncodeToString(salt)

	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, a.memory, a.time, a.threads, base64Salt, base64HashPassword)
	return encodedHash, nil
}

func decodeHash(encodedHash string) (argon2idHash *Argon2idHash, hash, salt []byte, err error) {
	values := strings.Split(encodedHash, "$")
	if len(values) != 6 {
		logrus.Warnf("Invalid hash format: %s", encodedHash)
		return nil, nil, nil, ErrInvalidHashFormat
	}
	var version int
	_, err = fmt.Sscanf(values[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, ErrArgon2VersionIncompatible
	}

	a := &Argon2idHash{}
	_, err = fmt.Sscanf(values[3], "m=%d,t=%d,p=%d", &a.memory, &a.time, &a.threads)
	if err != nil {
		return nil, nil, nil, err
	}

	base64Salt := values[4]
	base64HashPassword := values[5]

	salt, err = base64.RawStdEncoding.DecodeString(base64Salt)
	if err != nil {
		return nil, nil, nil, err
	}
	a.saltLen = uint32(len(salt))

	hash, err = base64.RawStdEncoding.DecodeString(base64HashPassword)
	if err != nil {
		return nil, nil, nil, err
	}
	a.keyLen = uint32(len(hash))

	return a, hash, salt, nil
}

// Compare compares a password with its hash.
func Compare(password, encodedHash string) (bool, error) {
	argon2idHash, hash, salt, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}
	hashPassword := argon2.IDKey([]byte(password), salt, argon2idHash.time, argon2idHash.memory, argon2idHash.threads, argon2idHash.keyLen)

	if subtle.ConstantTimeCompare(hashPassword, hash) == 1 {
		return true, nil
	}
	return false, nil
}
