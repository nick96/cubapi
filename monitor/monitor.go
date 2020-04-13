package monitor

import (
	"github.com/go-chi/chi"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type component interface {
	Name() string
	Check(logger *zap.Logger) map[string]interface{}
	CheckOk(logger *zap.Logger) bool
}

type DBComponent struct {
	DB *sqlx.DB
}

func (c DBComponent) Name() string {
	return "db"
}

func (c DBComponent) Check(logger *zap.Logger) map[string]interface{} {
	return map[string]interface{}{}
}

func (c DBComponent) CheckOk(logger *zap.Logger) bool {
	var err error
	for i := 0; i < 20; i++ {
		if err = c.DB.Ping(); err == nil {
			return true
		}
		logger.Error("Failed health check ping", zap.Int("try", i), zap.String("component", c.Name()), zap.Error(err))
	}
	logger.Error(
		"Failed health check",
		zap.String("check", "CheckOK"),
		zap.String("component", c.Name()),
		zap.Error(err),
	)
	return false
}

func NewHealthRouter(logger *zap.Logger, components ...component) func(chi.Router) {
	return func(r chi.Router) {

	}
}
