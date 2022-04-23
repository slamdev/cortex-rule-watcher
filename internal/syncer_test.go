package internal

import (
	"fmt"
	. "github.com/onsi/gomega"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"os"
	"path/filepath"
	"testing"
)

var (
	specStruct = monitoring.PrometheusRuleSpec{Groups: []monitoring.RuleGroup{{
		Name: "example",
		Rules: []monitoring.Rule{{
			Record: "job:http_inprogress_requests:sum",
			Expr:   intstr.FromString("sum by (job) (http_inprogress_requests)"),
		}},
	}}}
	specYaml = `
groups:
  - name: default-test.rule-example
    rules:
    - record: job:http_inprogress_requests:sum
      expr: sum by (job) (http_inprogress_requests)
`
)

func Test_ShouldCreateRule(t *testing.T) {
	RegisterTestingT(t)

	rulePath := t.TempDir()
	s := NewSyncer(rulePath)

	rule := monitoring.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RuleName,
			Namespace: Namespace,
		},
		Spec: specStruct,
	}
	expected := specYaml

	file := filepath.Join(rulePath, fmt.Sprintf("%s-%s.yaml", rule.Namespace, rule.Name))

	Expect(file).ShouldNot(BeAnExistingFile())

	Expect(s.replace(rule)).Should(Succeed())

	actual, err := os.ReadFile(file)
	Expect(err).NotTo(HaveOccurred())
	Expect(actual).Should(MatchYAML(expected))
}

func Test_ShouldUpdateRule(t *testing.T) {
	RegisterTestingT(t)

	rulePath := t.TempDir()
	s := NewSyncer(rulePath)

	rule := monitoring.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RuleName,
			Namespace: Namespace,
		},
	}

	file := filepath.Join(rulePath, fmt.Sprintf("%s-%s.yaml", rule.Namespace, rule.Name))

	Expect(s.replace(rule)).Should(Succeed())

	Expect(file).Should(BeARegularFile())

	rule.Spec = specStruct
	expected := specYaml

	Expect(s.replace(rule)).Should(Succeed())

	actual, err := os.ReadFile(file)
	Expect(err).NotTo(HaveOccurred())
	Expect(actual).Should(MatchYAML(expected))
}

func Test_ShouldDeleteRule(t *testing.T) {
	RegisterTestingT(t)

	rulePath := t.TempDir()
	s := NewSyncer(rulePath)

	rule := monitoring.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RuleName,
			Namespace: Namespace,
		},
	}

	file := filepath.Join(rulePath, fmt.Sprintf("%s-%s.yaml", rule.Namespace, rule.Name))

	Expect(s.replace(rule)).Should(Succeed())

	Expect(file).Should(BeARegularFile())

	Expect(s.delete(rule)).Should(Succeed())

	Expect(file).ShouldNot(BeAnExistingFile())
}
