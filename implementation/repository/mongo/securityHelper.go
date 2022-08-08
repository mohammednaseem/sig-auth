package mongo

import (
	"crypto/x509"
	"encoding/pem"
	"errors"

	"github.com/rs/zerolog/log"
)

func verifyCert(deviceCert []byte, caCert []byte) error {
	block, _ := pem.Decode(deviceCert)
	if block == nil {
		log.Error().Msg("Cannot decode Device Cert")
		return errors.New("cannot decode device cert")
	}

	leafCert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		log.Error().Err(err).Msg("Cannot parse leaf:")
		return err
	}

	rootPool := x509.NewCertPool()
	if !rootPool.AppendCertsFromPEM(caCert) {
		log.Error().Msg("Cannot append root")
		return errors.New("cannot append root")
	}

	_, err = leafCert.Verify(x509.VerifyOptions{
		Roots:     rootPool,
		KeyUsages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
	})
	if err != nil {
		log.Error().Err(err).Msg("Cert Verify Failed")
		return err
	}
	return err
}
