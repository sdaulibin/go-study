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
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"time"

	appv1 "github.com/sdaulibin/webserver-operator/api/v1"
	k8sappsv1 "k8s.io/api/apps/v1"
	k8scorev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// WebServerReconciler reconciles a WebServer object
type WebServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=my.domain,resources=webservers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=my.domain,resources=webservers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=my.domain,resources=webservers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebServer object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.1/pkg/reconcile
func (r *WebServerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("Webserver", req.NamespacedName)

	// TODO(user): your logic here
	instance := &appv1.WebServer{}
	err := r.Get(ctx, req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Webserver resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, err
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Webserver")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	// Check if the webserver deployment already exists, if not, create a new one
	foundDeployment := &k8sappsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name, Namespace: instance.Namespace}, foundDeployment)
	if err != nil && errors.IsNotFound(err) {
		log.Info("Deployment Not Found")
		//define a new deployment
		deployment := r.defineDeployment(instance)
		log.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Create(ctx, deployment)
		if err != nil {
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		// deployment created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	// Ensure the deployment replicas and image are the same as the spec
	var replicas int32 = int32(instance.Spec.Replicas)
	image := instance.Spec.Image
	var needUpd bool

	if *foundDeployment.Spec.Replicas != replicas {
		log.Info("Deployment spec.replicas change", "from", *foundDeployment.Spec.Replicas, "to", replicas)
		foundDeployment.Spec.Replicas = &replicas
		needUpd = true
	}

	if (*foundDeployment).Spec.Template.Spec.Containers[0].Image != image {
		log.Info("Deployment spec.template.spec.container[0].image change", "from", (*foundDeployment).Spec.Template.Spec.Containers[0].Image, "to", image)
		foundDeployment.Spec.Template.Spec.Containers[0].Image = image
		needUpd = true
	}
	if needUpd {
		err = r.Update(ctx, foundDeployment)
		if err != nil {
			log.Error(err, "Failed to update Deployment", "Deployment.Namespace", foundDeployment.Namespace, "Deployment.Name", foundDeployment.Name)
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		// Spec updated - return and requeue
		return ctrl.Result{Requeue: true}, nil
	}

	// Check if the webserver service already exists, if not, create a new one
	foundService := &k8scorev1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: instance.Name + "-service", Namespace: instance.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		//define a new service
		service := r.defineService(instance)
		log.Info("Creating a new Service", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
		err = r.Create(ctx, service)
		if err != nil {
			log.Error(err, "Failed to create new Servie", "Service.Namespace", service.Namespace, "Service.Name", service.Name)
			return ctrl.Result{RequeueAfter: time.Second * 5}, err
		}
		// Service created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Service")
		return ctrl.Result{RequeueAfter: time.Second * 5}, err
	}

	return ctrl.Result{RequeueAfter: time.Second * 10}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.WebServer{}).
		Complete(r)
}

func (r *WebServerReconciler) defineDeployment(webServer *appv1.WebServer) *k8sappsv1.Deployment {
	labels := labelForWebServer(webServer.Name)
	var replicas int32 = int32(webServer.Spec.Replicas)

	deployment := &k8sappsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: webServer.Namespace,
			Name:      webServer.Name,
		},
		Spec: k8sappsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: k8scorev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: k8scorev1.PodSpec{
					Containers: []k8scorev1.Container{
						{
							Name:            "webserver",
							Image:           webServer.Spec.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []k8scorev1.ContainerPort{
								{
									Name:          webServer.Name,
									Protocol:      k8scorev1.ProtocolSCTP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
	ctrl.SetControllerReference(webServer, deployment, r.Scheme)

	return deployment
}

func labelForWebServer(name string) map[string]string {
	return map[string]string{"app": "webserver", "webserver_cr": name}
}

func (r *WebServerReconciler) defineService(webServer *appv1.WebServer) *k8scorev1.Service {
	labels := labelForWebServer(webServer.Name)
	service := &k8scorev1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "app/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      webServer.Name + "-service",
			Namespace: webServer.Namespace,
			Labels:    labels,
		},
		Spec: k8scorev1.ServiceSpec{
			Type: k8scorev1.ServiceTypeNodePort,
			Ports: []k8scorev1.ServicePort{
				k8scorev1.ServicePort{
					Protocol: k8scorev1.ProtocolTCP,
					NodePort: 30010,
					Port:     80,
				},
			},
			Selector: map[string]string{
				"app":          "webserver",
				"webserver_cr": webServer.Name,
			},
		},
	}
	ctrl.SetControllerReference(webServer, service, r.Scheme)
	return service
}
