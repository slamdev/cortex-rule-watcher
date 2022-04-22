package internal

import (
	"context"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PrometheusRuleReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	RulePath string
}

func (r *PrometheusRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var rule monitoring.PrometheusRule
	if err := r.Get(ctx, req.NamespacedName, &rule); err != nil {
		if apierrors.IsNotFound(err) {
			l.Info("rule is deleted")
			return ctrl.Result{}, nil
		}
		l.Error(err, "unable to fetch rule")
		return ctrl.Result{}, err
	}

	l.Info("rule is found")

	return ctrl.Result{}, nil
}

func (r *PrometheusRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoring.PrometheusRule{}).
		Complete(r)
}
