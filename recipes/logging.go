package recipes

import (
	"context"
	"time"

	"github.com/andrewmthomas87/cookbook/models"
	"github.com/go-kit/kit/log"
)

type loggingService struct {
	logger log.Logger
	Service
}

// NewLoggingService returns a new instance of a logging Service.
func NewLoggingService(logger log.Logger, s Service) Service {
	return &loggingService{logger: logger, Service: s}
}

func (s *loggingService) GetRecipes(ctx context.Context) (rs []*models.Recipe, err error) {
	defer func(begin time.Time) {
		s.logger.Log(
			"method", "get recipes",
			"took", time.Since(begin),
			"err", err,
		)
	}(time.Now())
	return s.Service.GetRecipes(ctx)
}
