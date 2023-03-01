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
	"time"

	v1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	apiv1alpha1 "github.com/framsouza/kubebuilder-scaler/api/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var logger = log.Log.WithName("controller_scaler")

// ScalerReconciler reconciles a Scaler object
type ScalerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=api.framsouza.co,resources=scalers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=api.framsouza.co,resources=scalers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=api.framsouza.co,resources=scalers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Scaler object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.11.2/pkg/reconcile
func (r *ScalerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logger.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	log.Info("Reconcile was triggered")

	scaler := &apiv1alpha1.Scaler{}
	err := r.Get(ctx, req.NamespacedName, scaler)
	if err != nil {
		return ctrl.Result{}, nil
	}

	replicas := scaler.Spec.Replicas
	currentHour := time.Now().UTC().Hour()
	startTime := scaler.Spec.Start
	endTime := scaler.Spec.End

	log.Info(fmt.Sprintf("Current time %d", currentHour))

	//	if currentHour.After(startTime) && currentHour.Before(EndTime) {
	if currentHour >= startTime && currentHour <= endTime {

		for _, deploy := range scaler.Spec.Deployments {
			deployments := &v1.Deployment{}
			r.Get(ctx, types.NamespacedName{
				Namespace: deploy.Namespace,
				Name:      deploy.Name,
			},
				deployments,
			)

			if err != nil {
				return ctrl.Result{}, nil
			}

			if deployments.Spec.Replicas != &replicas {
				deployments.Spec.Replicas = &replicas
				err := r.Update(ctx, deployments)
				if err != nil {
					return ctrl.Result{}, err
				}
			}

		}
	}

	return ctrl.Result{RequeueAfter: time.Duration(30 * time.Second)}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ScalerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&apiv1alpha1.Scaler{}).
		Complete(r)
}
