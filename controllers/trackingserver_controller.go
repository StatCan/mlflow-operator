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

	"github.com/go-logr/logr"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	mlv1alpha1 "github.com/statcan/mlflow-operator/api/v1alpha1"
)

// TrackingServerReconciler reconciles a TrackingServer object
type TrackingServerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=ml.mlflow.org,resources=trackingservers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=ml.mlflow.org,resources=trackingservers/status,verbs=get;update;patch

// Reconcile the desired state configuration
func (r *TrackingServerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("trackingserver", req.NamespacedName)

	// Custom logic here
	reqLogger := r.Log.WithValues("Request.Namespace", req.Namespace, "Request.Name", req.Name)
	reqLogger.Info("Reconciling TrackingServer")

	// Fetch the TrackingServer instance
	instance := &mlv1alpha1.TrackingServer{}
	err := r.Client.Get(context.TODO(), req.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// Define a new Deployment object
	var deployment = deployMLFlow(instance)
	svc := svcMLFlow(instance)
	pvc := pvcMLFlow(instance)

	// Set TrackingServer instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := controllerutil.SetControllerReference(instance, svc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}
	if err := controllerutil.SetControllerReference(instance, pvc, r.Scheme); err != nil {
		return ctrl.Result{}, err
	}

	// Check if Deployment already exists
	found := &v1.Deployment{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Deployment", "Deployment.Namespace", deployment.Namespace, "Deployment.Name", deployment.Name)
		err = r.Client.Create(context.TODO(), deployment)
		if err != nil {
			return ctrl.Result{}, err
		}
		err2 := r.Client.Create(context.TODO(), svc)
		if err2 != nil {
			return ctrl.Result{}, err2
		}
		err3 := r.Client.Create(context.TODO(), pvc)
		if err3 != nil {
			return ctrl.Result{}, err3
		}
		return ctrl.Result{}, nil
	} else if err != nil {
		return ctrl.Result{}, err
	}

	reqLogger.Info("Skip reconciliation: Deployment already exists", "Deployment.Namespace", found.Namespace, "Deployment.Name", found.Name)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller
func (r *TrackingServerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&mlv1alpha1.TrackingServer{}).
		Complete(r)
}

func svcMLFlow(cr *mlv1alpha1.TrackingServer) *corev1.Service {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       cr.Name,
				"app.kubernetes.io/managed-by": "mlflow-operator",
			},
			OwnerReferences: []metav1.OwnerReference{},
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"app.kubernetes.io/name": cr.Name,
			},
			Ports: []corev1.ServicePort{{
				Name: "http",
				Port: 5000,
			}},
		},
	}
	return service
}

func pvcMLFlow(cr *mlv1alpha1.TrackingServer) *corev1.PersistentVolumeClaim {
	var storageClassName string

	if cr.Spec.StorageClass != "" {
		storageClassName = cr.Spec.StorageClass
	} else {
		storageClassName = "default"
	}

	size := cr.Spec.Size
	persistentVolumeClaim := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels: map[string]string{
				"app.kubernetes.io/name":       cr.Name,
				"app.kubernetes.io/managed-by": "mlflow-operator",
			},
			OwnerReferences: []metav1.OwnerReference{},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.ResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse(size),
				},
			},
			StorageClassName: &storageClassName,
		},
	}
	return persistentVolumeClaim
}

func deployMLFlow(cr *mlv1alpha1.TrackingServer) *v1.Deployment {
	replicas := cr.Spec.Replicas
	labels := map[string]string{
		"app.kubernetes.io/name":       cr.Name,
		"app.kubernetes.io/managed-by": "mlflow-operator",
	}
	container := []corev1.Container{{
		Image: cr.Spec.Image,
		Name:  cr.Name,
		Env: []corev1.EnvVar{
			{
				Name:  "MLFLOW_SERVER_FILE_STORE",
				Value: "sqlite:///mnt/mlruns/mlflow.db",
			},
			{
				Name:  "MLFLOW_SERVER_DEFAULT_ARTIFACT_ROOT",
				Value: "file:///mnt/mlruns/artifacts",
			},
		},
		Ports: []corev1.ContainerPort{{
			ContainerPort: 5000,
			Name:          "trackingserver",
		}},
		VolumeMounts: []corev1.VolumeMount{{
			MountPath: "/mnt/mlruns",
			Name:      "files-mlflow",
		}},
	}}

	if len(cr.Spec.S3secretName) != 0 {
		container[0].EnvFrom = []corev1.EnvFromSource{{
			SecretRef: &corev1.SecretEnvSource{LocalObjectReference: corev1.LocalObjectReference{Name: cr.Spec.S3secretName}},
		}}
	}

	if len(cr.Spec.S3endpointURL) != 0 {
		container[0].Env = []corev1.EnvVar{{
			Name:  "MLFLOW_S3_ENDPOINT_URL",
			Value: cr.Spec.S3endpointURL,
		}}
	}

	dep := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "apps/v1",
			Kind:       "Deployment",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
		},
		Spec: v1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: container,
					Volumes: []corev1.Volume{{
						Name: "files-mlflow",
						VolumeSource: corev1.VolumeSource{
							PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
								ClaimName: cr.Name,
							},
						},
					}},
				},
			},
		},
	}

	if len(cr.Spec.ImagePullSecret) != 0 {
		dep.Spec.Template.Spec.ImagePullSecrets = []corev1.LocalObjectReference{
			{Name: cr.Spec.ImagePullSecret},
		}
	}

	return dep

}
