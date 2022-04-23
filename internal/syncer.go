package internal

import (
	"context"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
)

type Syncer interface {
	replace(ctx context.Context, rule monitoring.PrometheusRule) error
	delete(ctx context.Context, rule monitoring.PrometheusRule) error
}

type syncer struct {
	rulePath string
}

func NewSyncer(rulePath string) Syncer {
	return &syncer{rulePath: rulePath}
}

func (s *syncer) replace(ctx context.Context, rule monitoring.PrometheusRule) error {
	//TODO implement me
	panic("implement me")
}

func (s *syncer) delete(ctx context.Context, rule monitoring.PrometheusRule) error {
	//TODO implement me
	panic("implement me")
}
