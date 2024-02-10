package token

import (
	"fmt"
	"time"

	"github.com/aead/chacha20poly1305"
	"github.com/o1egl/paseto"
)

type PasetoMaker struct {
	paseto *paseto.V2
	symentricKey []byte
}

func NewPasetoMaker(symentricKey string) (Maker, error) {
	if len(symentricKey) != chacha20poly1305.KeySize {
		return nil, fmt.Errorf("Invalid key size: must be exactly %d characters", chacha20poly1305.KeySize )
	}

	maker := &PasetoMaker{
		paseto: paseto.NewV2(),
		symentricKey: []byte(symentricKey),
	}

	return maker, nil
}

func (maker *PasetoMaker) CreateToken(username string, duration time.Duration) (string, *Payload, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", payload, err
	}

	token, err := maker.paseto.Encrypt(maker.symentricKey, payload, nil)
	return token, payload, err
}

func (maker *PasetoMaker) VerifyToken(token string) (*Payload, error) {
	payload := &Payload{}

	err := maker.paseto.Decrypt(token, maker.symentricKey, payload, nil)
	if err != nil {
		return nil, err
	}

	err = payload.Valid()
	if err != nil {
		return nil, err
	}

	return payload, nil
}