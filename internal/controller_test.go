package internal

import (
	"context"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ = Describe("PrometheusRule controller", func() {

	const (
		RuleName  = "test-rule"
		Namespace = "default"

		timeout  = time.Second * 10
		interval = time.Millisecond * 250
	)

	It("Should create PrometheusRule", func() {
		ctx := context.Background()
		rule := &monitoring.PrometheusRule{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "monitoring.coreos.com/v1",
				Kind:       "PrometheusRule",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      RuleName,
				Namespace: Namespace,
			},
			Spec: monitoring.PrometheusRuleSpec{},
		}
		Expect(k8sClient.Create(ctx, rule)).Should(Succeed())

		ruleLookupKey := types.NamespacedName{Name: RuleName, Namespace: Namespace}
		createdRule := &monitoring.PrometheusRule{}

		Eventually(func() bool {
			err := k8sClient.Get(ctx, ruleLookupKey, createdRule)
			if err != nil {
				return false
			}
			return true
		}, timeout, interval).Should(BeTrue())
	})
})
