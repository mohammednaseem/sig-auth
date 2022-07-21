package helper

import (
	"net/http"
	"strings"

	"github.com/device-auth/model"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
	"github.com/sirupsen/logrus"
)

func GetStatusCode(err error) int {
	if err == nil {
		return http.StatusOK
	}

	logrus.Error(err)
	switch err {
	case model.ErrInternalServerError:
		return http.StatusInternalServerError
	case model.ErrNotFound:
		return http.StatusNotFound
	case model.ErrConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

func GetStatusMessage(err error) string {
	if err == nil {
		return "OK"
	}

	logrus.Error(err)
	switch err {
	case model.ErrInternalServerError:
		return "Internal Server Error"
	case model.ErrNotFound:
		return "Resource Not Found"
	case model.ErrConflict:
		return "Conflict"
	default:
		return "Internal Server Error"
	}
}

//Database error handling - move to helper
func CheckDatabaseError(err error) (dberr error) {
	if pqErr, ok := err.(*pq.Error); ok {

		log.Error().Err(err).Msg("Database Error")
		if pqErr.Code.Name() == "unique_violation" {
			dberr = model.ErrConflict
			log.Error().Err(dberr).Msg("unique_violation")
			return
		} else if pqErr.Code.Name() == "foreign_key_violation" {
			dberr = model.ErrBadParamInput
			log.Error().Err(dberr).Msg("foreign_key_violation")
			return
		} else if pqErr.Error() == "pq: the database system is starting up" {
			dberr = model.ErrServiceUnavailable
			log.Error().Err(dberr).Msg("pq: the database system is starting up")
			return
		} else if strings.Contains(pqErr.Error(), "violates not-null constraint") {
			dberr = model.ErrBadParamInput
			log.Error().Err(dberr).Msg("item does not exits")
			return
		}
	} else if strings.Contains(err.Error(), "No connection could be made because the target machine actively refused it") {
		dberr = model.ErrServiceUnavailable
		log.Error().Err(dberr).Msg("No connection could be made because the target machine actively refused it")
		return model.ErrServiceUnavailable
	} else if err.Error() == "sql: no rows in result set" {
		dberr = model.ErrBadParamInput
		log.Error().Err(dberr).Msg("sql: no rows in result set")
		return
	}
	log.Error().Err(err).Msg("Unknown Database Error")
	return err
}
