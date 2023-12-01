/*
Copyright 2023.

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
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	examplev1alpha1 "github.com/yitsushi/controller-ownref-experiment/api/v1alpha1"
)

// MyResReconciler reconciles a MyRes object
type MyResReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=example.k8s.experiments.efertone.me,resources=myres,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=example.k8s.experiments.efertone.me,resources=myres/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=example.k8s.experiments.efertone.me,resources=myres/finalizers,verbs=update
func (r *MyResReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx).WithName("MyResController")

	var myres examplev1alpha1.MyRes
	if err := r.Get(ctx, req.NamespacedName, &myres); err != nil {
		logger.Error(err, "Hit an error", "namespacedName", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Being deleted.
	if !myres.ObjectMeta.DeletionTimestamp.IsZero() {
		if err := r.finalizer(&myres); err != nil {
			logger.Error(err, "finalizer is not happy")
		}

		myres.ObjectMeta.Finalizers = removeString(myres.ObjectMeta.Finalizers, examplev1alpha1.MyResFinalizer)

		return reconcile.Result{}, r.Update(context.Background(), &myres)
	}

	// Add Finalizer.
	if err := r.addFinalizer(&myres); err != nil {
		logger.Error(err, "unable to add finalizer")
		return ctrl.Result{}, err
	}

	// Create secret.
	fancySecret, err := r.createSecret(ctx, &myres)
	if err != nil {
		logger.Error(err, "unable to create secret")
		return ctrl.Result{}, err
	}

	logger.Info("secret is ready", "name", fancySecret.GetName())

	return ctrl.Result{}, nil
}

func (r *MyResReconciler) addFinalizer(myres *examplev1alpha1.MyRes) error {
	// We didn't set Finalizer yet.
	if !containsString(myres.ObjectMeta.Finalizers, examplev1alpha1.MyResFinalizer) {
		myres.ObjectMeta.Finalizers = append(myres.ObjectMeta.Finalizers, examplev1alpha1.MyResFinalizer)
		if err := r.Update(context.Background(), myres); err != nil {
			return err
		}
	}

	return nil
}

func (r *MyResReconciler) finalizer(myres *examplev1alpha1.MyRes) error {
	if !containsString(myres.ObjectMeta.Finalizers, examplev1alpha1.MyResFinalizer) {
		return nil
	}

	// We do not delete the Secret here.
	return nil
}

func (r *MyResReconciler) createSecret(ctx context.Context, myres *examplev1alpha1.MyRes) (v1.Secret, error) {
	secretName := fmt.Sprintf("%s-fancy-secret", myres.GetName())
	namespacedName := types.NamespacedName{
		Namespace: myres.GetNamespace(),
		Name:      secretName,
	}

	var secret v1.Secret
	if err := r.Get(ctx, namespacedName, &secret); err == nil {
		// Secret already exists.
		return secret, nil
	}

	fancySecret := v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretName,
			Namespace: myres.GetNamespace(),
			OwnerReferences: []metav1.OwnerReference{
				{
					APIVersion: examplev1alpha1.GroupVersion.Group + "/" + examplev1alpha1.GroupVersion.Version,
					Kind:       examplev1alpha1.MyResKind,
					Name:       myres.GetName(),
					UID:        myres.GetUID(),
				},
			},
		},
		Type: v1.SecretTypeOpaque,
		Data: map[string][]byte{},
	}

	if err := r.Client.Create(ctx, &fancySecret); err != nil {
		return secret, err
	}

	return secret, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyResReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&examplev1alpha1.MyRes{}).
		Complete(r)
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func removeString(slice []string, s string) (result []string) {
	for _, item := range slice {
		if item == s {
			continue
		}
		result = append(result, item)
	}
	return
}
