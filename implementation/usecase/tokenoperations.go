package usecase

import (
	"errors"
	"strings"

	"github.com/device-auth/model"
	jwt "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// Verify a JWT token using an RSA public key
func VerifyJWT(token string, mdevice model.Device, algorithm string) (bool, error) {
	var publicCerts []string
	for _, element := range mdevice.Credentials {
		if len(strings.TrimSpace(element.PublicKey.Key)) != 0 {
			publicCerts = append(publicCerts, element.PublicKey.Key)
		}
	}

	// parse token // verify with all available public certificates
	state, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		var err error
		for _, publicCert := range publicCerts {

			if algorithm == "RS256" {
				key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicCert))
				log.Error().Err(err).Msg("")
				if err == nil {
					return key, nil
				}
			} else { //EC
				key, err := jwt.ParseECPublicKeyFromPEM([]byte(publicCert))
				if err == nil {
					return key, nil
				}
			}

		}
		return err, nil
	})
	if err != nil {
		return false, err
	}
	if !state.Valid {
		log.Error().Err(err).Msg("invalid jwt token")
		return false, errors.New("invalid jwt token")
	}
	return true, nil
}

func IdentifyAndVerifyJWT(token string, mDevice model.Device) (bool, error) {
	signingMethod := ""

	// parse token
	jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			signingMethod = "ES256"
		} else if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			signingMethod = "RS256"
		} else if signingMethod == "" {
			return false, errors.New("unknown signing method")
		}
		return true, nil
	})
	boolVal, err := VerifyJWT(token, mDevice, signingMethod)
	if !boolVal {
		log.Error().Err(err).Msg("")
		return false, err
	}
	return true, nil
}
