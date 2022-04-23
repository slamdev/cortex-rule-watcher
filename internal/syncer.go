package internal

import (
	"fmt"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"os"
	"path/filepath"
	"sigs.k8s.io/yaml"
	"syscall"
)

type Syncer interface {
	replace(rule monitoring.PrometheusRule) error
	delete(rule monitoring.PrometheusRule) error
}

type syncer struct {
	rulePath string
}

func NewSyncer(rulePath string) Syncer {
	return &syncer{rulePath: rulePath}
}

func (s *syncer) replace(rule monitoring.PrometheusRule) error {
	out, err := yaml.Marshal(rule.Spec)
	if err != nil {
		return fmt.Errorf("failed to marshal rule to yaml; %w", err)
	}

	fileName := s.fileName(rule)
	if err := os.WriteFile(fileName, out, 0644); err != nil {
		return fmt.Errorf("failed to write yaml to %s; %w", fileName, err)
	}

	return nil
}

func (s *syncer) delete(rule monitoring.PrometheusRule) error {
	fileName := s.fileName(rule)
	if err := os.Remove(fileName); err != nil {
		e, ok := err.(*os.PathError)
		if ok && e.Err == syscall.ENOENT {
			return nil
		}
		return fmt.Errorf("failed to delete %s; %w", fileName, err)
	}

	return nil
}

func (s *syncer) fileName(rule monitoring.PrometheusRule) string {
	return filepath.Join(s.rulePath, fmt.Sprintf("%s-%s.yaml", rule.Namespace, rule.Name))
}
