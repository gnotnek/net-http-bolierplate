package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"sync"

	"github.com/rs/zerolog/log"
)

var (
	key       []byte
	encryptIv []byte
	once      sync.Once
	initErr   error
)

func Init(secretKey, secretKeyIv, encryptMethod string) error {
	once.Do(func() {
		if secretKey == "" || secretKeyIv == "" || encryptMethod == "" {
			initErr = errors.New("secret key, IV, and encryption method must be set")
			log.Fatal().Err(initErr).Msg("encrypt.Init failed")
			return
		}

		if encryptMethod != "aes-256-cbc" {
			initErr = errors.New("unsupported encryption method: " + encryptMethod)
			log.Fatal().Err(initErr).Msg("encrypt.Init failed")
			return
		}

		hash := sha512.Sum512([]byte(secretKey))
		key = []byte(hex.EncodeToString(hash[:]))[:32]

		hashIv := sha512.Sum512([]byte(secretKeyIv))
		encryptIv = []byte(hex.EncodeToString(hashIv[:]))[:16]

		log.Info().Msg("encryption utility initialized")
	})
	return initErr
}

func EncryptData(data string) (string, error) {
	if key == nil || encryptIv == nil {
		return "", errors.New("encrypt not initialized, call encrypt.Init() first")
	}

	log.Info().Str("func", "EncryptData").Msg("Encrypting data")

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create cipher block")
		return "", err
	}

	mode := cipher.NewCBCEncrypter(block, encryptIv)
	plaintext := pkcs7Pad([]byte(data), aes.BlockSize)
	ciphertext := make([]byte, len(plaintext))
	mode.CryptBlocks(ciphertext, plaintext)

	encoded := base64.StdEncoding.EncodeToString(ciphertext)
	log.Info().Str("func", "EncryptData").Msg("Encryption complete")
	return encoded, nil
}

func DecryptData(encryptedBase64 string) (string, error) {
	if key == nil || encryptIv == nil {
		return "", errors.New("encrypt not initialized, call encrypt.Init() first")
	}

	log.Info().Str("func", "DecryptData").Msg("Decrypting data")

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		log.Error().Err(err).Msg("Failed to decode base64 input")
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		log.Error().Err(err).Msg("Failed to create cipher block")
		return "", err
	}

	if len(ciphertext)%aes.BlockSize != 0 {
		log.Error().Msg("Invalid ciphertext block size")
		return "", errors.New("invalid ciphertext block size")
	}

	mode := cipher.NewCBCDecrypter(block, encryptIv)
	plaintext := make([]byte, len(ciphertext))
	mode.CryptBlocks(plaintext, ciphertext)

	plaintext, err = pkcs7Unpad(plaintext, aes.BlockSize)
	if err != nil {
		log.Error().Err(err).Msg("Failed to unpad decrypted data")
		return "", err
	}

	log.Info().Str("func", "DecryptData").Msg("Decryption complete")
	return string(plaintext), nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	return append(data, padText...)
}

func pkcs7Unpad(data []byte, blockSize int) ([]byte, error) {
	if len(data) == 0 || len(data)%blockSize != 0 {
		return nil, errors.New("invalid padding size")
	}

	padding := int(data[len(data)-1])
	if padding == 0 || padding > blockSize {
		return nil, errors.New("invalid padding value")
	}

	for i := len(data) - padding; i < len(data); i++ {
		if int(data[i]) != padding {
			return nil, errors.New("invalid padding byte found")
		}
	}

	return data[:len(data)-padding], nil
}
