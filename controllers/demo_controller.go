/*
Copyright 2022.

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
	"strconv"

	demov1 "github.com/YiVang/operator-demo/api/v1"
	v1 "github.com/YiVang/operator-demo/api/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// DemoReconciler reconciles a Demo object
type DemoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=demo.github.com,resources=demoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=demo.github.com,resources=demoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=demo.github.com,resources=demoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Demo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *DemoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	// TODO(user): your logic here
	log.Info("Receive change notify...")
	inst := &v1.Demo{}
	if err := r.Get(ctx, req.NamespacedName, inst); err != nil && !errors.IsNotFound(err) {
		log.Error(err, "get demo inst failed")
		return ctrl.Result{}, err
	}
	expectReplicas := inst.Spec.Replicas
	realReplicas := 0
	if inst.Spec.RealReplicas != nil {
		realReplicas = *(inst.Spec.RealReplicas)
	}

	log.Info("Replicas ", "expect", strconv.Itoa(expectReplicas), "current", strconv.Itoa(realReplicas))

	if expectReplicas == realReplicas {
		log.Info("don't need change")
		return ctrl.Result{}, nil
	}
	if expectReplicas > realReplicas {
		log.Info("start scale out...")
	} else {
		log.Info("start scale in...")
	}
	inst.Spec.RealReplicas = &expectReplicas
	// ??????update?????????????????????????????????Reconcile????????????
	if err := r.Update(ctx, inst); err != nil {
		log.Error(err, "update inst info failed")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DemoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&demov1.Demo{}).
		Complete(r)
}
