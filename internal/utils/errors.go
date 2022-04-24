package utils

import "go.uber.org/zap"

type InternalServerError struct {
}

func (e *InternalServerError) Error() string {
	return "internal server error"
}

func LogServerError(err error, logger *zap.SugaredLogger) error {
	if err == nil {
		return nil
	}

	logger.With("error", err).Warn("error creating game")
	return &InternalServerError{}
}
