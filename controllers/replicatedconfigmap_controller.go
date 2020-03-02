/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"strings"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"

	rcmv1beta1 "github.com/abatilo/replicatedconfigmap/api/v1beta1"
)

// ReplicatedConfigMapReconciler reconciles a ReplicatedConfigMap object
type ReplicatedConfigMapReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=rcm.aaronbatilo.dev,resources=replicatedconfigmaps,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rcm.aaronbatilo.dev,resources=replicatedconfigmaps/status,verbs=get;update;patch

// Read only on namespaces
// +kubebuilder:rbac:groups="",resources=namespace,verbs=get;list;watch

// Read/write configmaps
// +kubebuilder:rbac:groups="",resources=configmap,verbs=get;list;watch;create;update;patch;delete

const (
	// MatchingAnnotation is the annotation to use to know which namespaces to operate on
	MatchingAnnotation = "rcm-sync"
)

// Reconcile is the loop where we integrate logic for repairing expected state
func (r *ReplicatedConfigMapReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("replicatedconfigmap", req.NamespacedName)

	// TODO Find out if there's a way to use the `req.NamespacedName` with non-namespaced resources
	replicatedConfigMaps := &rcmv1beta1.ReplicatedConfigMapList{}
	if err := r.List(ctx, replicatedConfigMaps); err != nil {
		log.Error(err, "couldn't get replicatedConfigMaps")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var allNamespaces corev1.NamespaceList
	var childNamespaces []corev1.Namespace
	if err := r.List(ctx, &allNamespaces); err != nil {
		log.Error(err, "unable to list namespaces")
		return ctrl.Result{}, err
	}

	// TODO There's MatchingFields and MatchingLabels, but it's not clear if
	// there was an easier way to filter for namespaces with a matching
	// annotation

	// https://github.com/kubernetes-sigs/controller-runtime/blob/9f8aab6551a3ecc211b6c6078245899495b7e0f5/pkg/client/options.go#L434-L436
	// https://github.com/kubernetes-sigs/controller-runtime/blob/9f8aab6551a3ecc211b6c6078245899495b7e0f5/pkg/client/options.go#L379-L380
	for _, namespace := range allNamespaces.Items {
		for k, v := range namespace.Annotations {
			if k == MatchingAnnotation && strings.EqualFold("true", strings.TrimSpace(v)) {
				childNamespaces = append(childNamespaces, namespace)
			}
		}
	}

	// For each ReplicatedConfigMap available, loop over each syncable namespace
	// and either create or overwrite ConfigMaps specified by the ReplicatedConfigMap
	for _, replicatedConfigMap := range replicatedConfigMaps.Items {
		for _, namespace := range childNamespaces {
			objectKey := client.ObjectKey{Namespace: namespace.Name, Name: replicatedConfigMap.Spec.Name}
			configMap := &corev1.ConfigMap{}
			if err := r.Get(ctx, objectKey, configMap); err != nil {
				log.Info("couldn't find configmap")
				if errors.IsNotFound(err) {
					newConfigMap := &corev1.ConfigMap{
						ObjectMeta: metav1.ObjectMeta{
							Name:      replicatedConfigMap.Spec.Name,
							Namespace: namespace.Name,
						},
						Data: replicatedConfigMap.Spec.Data,
					}

					if err := ctrl.SetControllerReference(&replicatedConfigMap, newConfigMap, r.Scheme); err != nil {
						log.Error(err, "couldn't set controller reference")
					}

					if err := r.Create(ctx, newConfigMap); err != nil {
						log.Error(err, "couldn't create new configmap")
					}
				}
			} else {
				configMap.Data = replicatedConfigMap.Spec.Data
				if err := r.Update(ctx, configMap); err != nil {
					log.Error(err, "couldn't update")
				}
			}
		}
	}

	// TODO Spend more time understanding how `ctrl.Result{}` and errors cause repeat loops.
	// This implementation assumes a lot of happy path scenarios
	return ctrl.Result{}, nil
}

// SetupWithManager registered resources that this reconciler needs to pay attention to
func (r *ReplicatedConfigMapReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// TODO Learn more about FieldIndexers as used here:
	// https://book.kubebuilder.io/cronjob-tutorial/controller-implementation.html#setup
	return ctrl.NewControllerManagedBy(mgr).
		For(&rcmv1beta1.ReplicatedConfigMap{}).
		Watches(&source.Kind{Type: &corev1.Namespace{}}, &handler.EnqueueRequestForObject{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
