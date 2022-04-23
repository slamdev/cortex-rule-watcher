package internal

import (
	"context"
	"fmt"
	monitoring "github.com/prometheus-operator/prometheus-operator/pkg/apis/monitoring/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

type PrometheusRuleReconciler struct {
	client.Client
	Scheme *runtime.Scheme
	Syncer Syncer
}

func (r *PrometheusRuleReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx)

	var rule monitoring.PrometheusRule
	if err := r.Get(ctx, req.NamespacedName, &rule); err != nil {
		if apierrors.IsNotFound(err) {
			rule = monitoring.PrometheusRule{
				ObjectMeta: metav1.ObjectMeta{
					Name:      req.Name,
					Namespace: req.Namespace,
				},
			}
			if err := r.Syncer.delete(ctx, rule); err != nil {
				return ctrl.Result{}, fmt.Errorf("failed to sync rule deletion event; %w", err)
			}
			return ctrl.Result{}, nil
		}
		l.Error(err, "unable to fetch rule")
		return ctrl.Result{}, err
	}

	if err := r.Syncer.replace(ctx, rule); err != nil {
		return ctrl.Result{}, fmt.Errorf("failed to sync rule upsert event; %w", err)
	}

	return ctrl.Result{}, nil
}

func (r *PrometheusRuleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoring.PrometheusRule{}).
		Complete(r)
}
