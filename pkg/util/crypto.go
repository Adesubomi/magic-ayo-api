package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"io"
)

type AesCrypt struct {
	AppKey string
}

// EncryptAES https://www.melvinvivas.com/how-to-encrypt-and-decrypt-data-using-aes
func (aesCrypto AesCrypt) EncryptAES(text string) (string, error) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(aesCrypto.AppKey)
	plaintext := []byte(text)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext), nil
}

func (aesCrypto AesCrypt) DecryptAES(encryptedString string) (string, error) {
	key, _ := hex.DecodeString(aesCrypto.AppKey)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s", plaintext), nil
}

func BcryptHash(text string) (string, error) {
	result, err := bcrypt.GenerateFromPassword([]byte(text), bcrypt.DefaultCost)

	if err != nil {
		return "", err
	}

	return string(result), nil
}

func BcryptCompare(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (aesCrypto AesCrypt) GenerateCheckSum(jsonStr []byte) string {
	mac := hmac.New(sha256.New, []byte(aesCrypto.AppKey))
	mac.Write(jsonStr)
	return hex.EncodeToString(mac.Sum(nil))
}

func (aesCrypto AesCrypt) ValidateCheckSum(jsonStr, checksum string) bool {
	generatedChecksum := aesCrypto.GenerateCheckSum([]byte(jsonStr))
	return hmac.Equal([]byte(generatedChecksum), []byte(checksum))
}
