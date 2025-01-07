/*
Copyright 2025.

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

package controller

import (
	"context"
	"fmt"
	"os"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	monitoringv1alpha1 "github.com/NoNickeD/operator-sdk-demo/api/v1alpha1"
	"github.com/go-logr/logr"
	"github.com/pingcap/errors"
	corev1 "k8s.io/api/core/v1"
)

// PodNotifRestartReconciler reconciles a PodNotifRestart object
type PodNotifRestartReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=monitoring.vodafone.com,resources=podnotifrestarts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.vodafone.com,resources=podnotifrestarts/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=monitoring.vodafone.com,resources=podnotifrestarts/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodNotifRestart object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.18.4/pkg/reconcile

func (r *PodNotifRestartReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1alpha1.PodNotifRestart{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}

func (r *PodNotifRestartReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("podnotifrestart", req.NamespacedName)

	var pnr monitoringv1alpha1.PodNotifRestart
	if err := r.Client.Get(ctx, req.NamespacedName, &pnr); err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		log.Error(err, "unable to fetch PodNotifRestart")
		return ctrl.Result{}, err
	}

	var podList corev1.PodList
	if err := r.Client.List(ctx, &podList, client.InNamespace(pnr.Namespace)); err != nil {
		log.Error(err, "unable to list pods")
		return ctrl.Result{}, err
	}

	// Initialize notifiers
	discord := &DiscordNotifier{WebhookURL: os.Getenv("DISCORD_WEBHOOK_URL")}
	teams := &TeamsNotifier{WebhookURL: os.Getenv("TEAMS_WEBHOOK_URL")}
	slack := &SlackNotifier{WebhookURL: os.Getenv("SLACK_WEBHOOK_URL")}

	for _, pod := range podList.Items {
		for _, status := range pod.Status.ContainerStatuses {
			if status.RestartCount >= pnr.Spec.MinRestarts {
				message := fmt.Sprintf("Pod %s has restarted %d times", pod.Name, status.RestartCount)
				log.Info("Sending restart notification", "pod", pod.Name, "restartCount", status.RestartCount)

				if err := sendNotification(message, discord, teams, slack); err != nil {
					log.Error(err, "failed to send notification")
					return ctrl.Result{}, err
				}
			}
		}
	}

	return ctrl.Result{RequeueAfter: time.Minute * 2}, nil
}
