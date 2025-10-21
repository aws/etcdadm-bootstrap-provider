package controllers

import (
	"context"

	"github.com/pkg/errors"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// MachineToBootstrapMapFunc is a handler.ToRequestsFunc to be used to enqueue
// requests for reconciliation of EtcdadmConfig.
func (r *EtcdadmConfigReconciler) MachineToBootstrapMapFunc(ctx context.Context, o client.Object) []ctrl.Request {
	var result []ctrl.Request

	m, ok := o.(*clusterv1.Machine)
	if !ok {
		r.Log.Error(errors.Errorf("expected a Machine but got a %T", o.GetObjectKind()), "failed to get EtcdadmConfigs for Machine")
		return nil
	}
	if m.Spec.Bootstrap.ConfigRef.IsDefined() && m.Spec.Bootstrap.ConfigRef.Kind == "EtcdadmConfig" {
		name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
		result = append(result, ctrl.Request{NamespacedName: name})
	}
	return result
}

// ClusterToEtcdadmConfigs is a handler.ToRequestsFunc to be used to enqeue
// requests for reconciliation of EtcdadmConfigs.
func (r *EtcdadmConfigReconciler) ClusterToEtcdadmConfigs(ctx context.Context, o client.Object) []ctrl.Request {
	var result []ctrl.Request

	c, ok := o.(*clusterv1.Cluster)
	if !ok {
		r.Log.Error(errors.Errorf("expected a Cluster but got a %T", o.GetObjectKind()), "failed to get EtcdadmConfigs for Cluster")
		return nil
	}

	selectors := []client.ListOption{
		client.InNamespace(c.Namespace),
		client.MatchingLabels{
			clusterv1.ClusterNameLabel: c.Name,
		},
	}

	machineList := &clusterv1.MachineList{}
	if err := r.Client.List(ctx, machineList, selectors...); err != nil {
		r.Log.Error(err, "failed to list Machines", "Cluster", c.Name, "Namespace", c.Namespace)
		return nil
	}

	for _, m := range machineList.Items {
		if m.Spec.Bootstrap.ConfigRef.IsDefined() &&
			m.Spec.Bootstrap.ConfigRef.Kind == "EtcdadmConfig" {
			name := client.ObjectKey{Namespace: m.Namespace, Name: m.Spec.Bootstrap.ConfigRef.Name}
			result = append(result, ctrl.Request{NamespacedName: name})
		}
	}
	return result
}
