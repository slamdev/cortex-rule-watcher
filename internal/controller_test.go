package internal

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes/scheme"
	"path/filepath"
	"reflect"
	"runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"testing"
)

const (
	RuleName  = "test-rule"
	Namespace = "default"
)

func Test_ShouldReconcile(t *testing.T) {
	RegisterTestingT(t)

	k8sClient, syncer, ctx := setup(t)

	rule := &monitoring.PrometheusRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      RuleName,
			Namespace: Namespace,
		},
	}

	// Verify rule creation
	syncer.EXPECT().replace(gomock.Any(), gomock.All(matchRule(rule)))
	Expect(k8sClient.Create(ctx, rule)).Should(Succeed())

	// Verify rule update
	rule.Spec = monitoring.PrometheusRuleSpec{Groups: []monitoring.RuleGroup{{
		Name: "some-group",
		Rules: []monitoring.Rule{{
			Record: "some",
			Expr:   intstr.FromInt(1),
		}},
	}}}
	syncer.EXPECT().replace(gomock.Any(), gomock.All(matchRule(rule)))
	Expect(k8sClient.Update(ctx, rule)).Should(Succeed())

	// Verify rule deletion
	rule.Spec = monitoring.PrometheusRuleSpec{}
	Expect(k8sClient.Delete(ctx, rule)).Should(Succeed())
	syncer.EXPECT().delete(gomock.Any(), gomock.All(matchRule(rule)))
}

func setup(t *testing.T) (client.Client, *MockSyncer, context.Context) {
	t.Setenv("KUBEBUILDER_ASSETS", filepath.Join("testdata", runtime.GOOS, runtime.GOARCH, "bin"))
	logf.SetLogger(zap.New(zap.UseDevMode(true)))
	ctx, cancel := context.WithCancel(context.TODO())

	testEnv := &envtest.Environment{
		CRDDirectoryPaths:     []string{"testdata"},
		ErrorIfCRDPathMissing: true,
	}

	cfg, err := testEnv.Start()
	Expect(err).NotTo(HaveOccurred())
	Expect(cfg).NotTo(BeNil())

	err = monitoring.AddToScheme(scheme.Scheme)
	Expect(err).NotTo(HaveOccurred())

	k8sClient, err := client.New(cfg, client.Options{Scheme: scheme.Scheme})
	Expect(err).NotTo(HaveOccurred())
	Expect(k8sClient).NotTo(BeNil())

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	Expect(err).ToNot(HaveOccurred())

	mockCtrl := gomock.NewController(t)
	syncer := NewMockSyncer(mockCtrl)

	err = (&PrometheusRuleReconciler{
		Client: k8sManager.GetClient(),
		Scheme: k8sManager.GetScheme(),
		Syncer: syncer,
	}).SetupWithManager(k8sManager)
	Expect(err).ToNot(HaveOccurred())

	go func() {
		err = k8sManager.Start(ctx)
		Expect(err).ToNot(HaveOccurred(), "failed to run manager")
	}()

	t.Cleanup(func() {
		cancel()
		err := testEnv.Stop()
		Expect(err).NotTo(HaveOccurred())
	})

	return k8sClient, syncer, ctx
}

type ruleMatcher struct {
	x monitoring.PrometheusRule
}

func (e ruleMatcher) Matches(x interface{}) bool {
	r := x.(monitoring.PrometheusRule)
	return e.x.Namespace == r.Namespace && e.x.Name == r.Name && reflect.DeepEqual(e.x.Spec, r.Spec)
}

func (e ruleMatcher) String() string {
	return fmt.Sprintf("is equal to %v", e.x)
}

func matchRule(rule *monitoring.PrometheusRule) ruleMatcher {
	return ruleMatcher{x: *rule}
}
