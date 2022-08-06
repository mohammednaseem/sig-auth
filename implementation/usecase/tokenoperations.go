package usecase

import (
	"errors"
	jwt "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// Verify a JWT token using an RSA public key
func VerifyJWT(token string, Certs []string, algorithm string) (bool, error) {

	// parse token // verify with all available public certificates
	state, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {

		var err error
		for _, publicCert := range Certs {

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

func IdentifyAndVerifyJWT(token string, Certs []string) (bool, error) {
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

	boolVal, err := VerifyJWT(token, Certs, signingMethod)
	if !boolVal {
		log.Error().Err(err).Msg("")
		return false, err
	}
	return true, nil
}
