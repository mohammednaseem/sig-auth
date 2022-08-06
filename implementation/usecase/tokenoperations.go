package usecase

import (
	"errors"

	jwt "github.com/golang-jwt/jwt"
	"github.com/rs/zerolog/log"
)

// Verify a JWT token using an RSA public key
func VerifyJWT(token string, Certs []string) (bool, error, string) {
	var algorithm string
	// parse token // verify with all available public certificates
	state, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodECDSA); ok {
			algorithm = "ES256"
		} else if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
			algorithm = "RS256"
		} else if algorithm == "" {
			return false, errors.New("unknown signing method")
		}
		var err error
		for _, publicCert := range Certs {
			print(publicCert)
			if algorithm == "RS256" {
				key, err := jwt.ParseRSAPublicKeyFromPEM([]byte(test))
				log.Error().Err(err).Msg("")
				if err == nil {
					return key, nil
				}
			} else { //EC
				key, err := jwt.ParseECPublicKeyFromPEM([]byte(test))
				if err == nil {
					return key, nil
				}
			}

		}
		return err, nil
	})
	if err != nil {
		return false, err, ""
	}
	if !state.Valid {
		log.Error().Err(err).Msg("invalid jwt token")
		return false, errors.New("invalid jwt token"), ""
	}
	return true, nil, algorithm
}

func IdentifyAndVerifyJWT(token string, Certs []string) (bool, error, string) {

	boolVal, err, algorithm := VerifyJWT(token, Certs)
	if !boolVal {
		log.Error().Err(err).Msg("")
		return false, err, ""
	}
	return true, nil, algorithm
}
