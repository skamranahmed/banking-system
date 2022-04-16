package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

// PasetoMaker : implements the Maker interface
type PasetoMaker struct {
	paseto       *paseto.V2
	symmetricKey []byte
}

// NewPasetoMaker : returns a new PasetoMaker
func NewPasetoMaker(symmetricKey string) (Maker, error) {
	if len(symmetricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("invalid key size, the key size must be %d characters", chacha20poly1305.KeySize)
	}

	maker := &PasetoMaker{
		paseto:       paseto.NewV2(),
		symmetricKey: []byte(symmetricKey),
	}

	return maker, nil
}

// CreateToken : creates a new paseto token for the provided userID and duration
func (maker *PasetoMaker) CreateToken(userID uint, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(userID, duration)
	if err != nil {
		return "", nil, err
	}

	token, err := maker.paseto.Encrypt(maker.symmetricKey, payload, nil)
	if err != nil {
		return "", nil, err
	}

	return token, payload, nil
}

// VerifyToken : checks whether the paseto token is valid or not
func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symmetricKey, payload, nil)
	if err != nil {
		return nil, ErrInvalidToken
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}
