package statefulset

import (
	"context"
	"strconv"
	"time"

	apiapps "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_statefulset")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new StatefulSet Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileStatefulSet{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("statefulset-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource StatefulSet
	//err = c.Watch(&source.Kind{Type: &apiapps.StatefulSet{}}, &handler.EnqueueRequestForObject{})
	//if err != nil {
	//	return err
	//}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner StatefulSet
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileStatefulSet implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileStatefulSet{}

// ReconcileStatefulSet reconciles a StatefulSet object
type ReconcileStatefulSet struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a StatefulSet object and makes changes based on the state read
// and what is in the StatefulSet.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileStatefulSet) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling StatefulSet")

	// Fetch the StatefulSet instance
	//instance := &apiapps.StatefulSet{}
	//err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	//if err != nil {
	//	if errors.IsNotFound(err) {
	//		// Request object not found, could have been deleted after reconcile request.
	//		// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
	//		// Return and don't requeue
	//		return reconcile.Result{}, nil
	//	}
	//	// Error reading the object - requeue the request.
	//	return reconcile.Result{}, err
	//}
	//reqLogger.Info(time.Now().String() + "=========----------- statefulSet:  " + instance.String())

	pod := &corev1.Pod{}
	err := r.client.Get(context.TODO(), request.NamespacedName, pod)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info(time.Now().String() + "================== not find pod:  " + request.NamespacedName.Name)
		reqLogger.Info(time.Now().String() + "++++++++++++++++++statefulSet get pod: " + err.Error())
	}
	if pod.DeletionTimestamp != nil {
		var GracePeriodSeconds int64 = 0
		reqLogger.Info(time.Now().String() + "==================statefulSet pod deletion: " + pod.DeletionTimestamp.String())
		err = r.client.Delete(context.TODO(), pod, &client.DeleteOptions{GracePeriodSeconds: &GracePeriodSeconds})
		if err != nil {
			reqLogger.Info(time.Now().String() + "==================delete pod error : " + err.Error())
		}
	}
	reqLogger.Info(time.Now().String() + "==================statefulSet pod status: " + string(pod.Status.Phase))
	reqLogger.Info(time.Now().String() + "==================statefulSet pod: " + string(pod.String()))

	//podNames := getPodNames(instance)
	//for _, podName := range podNames {
	//	pod := &corev1.Pod{}
	//	err = r.client.Get(context.TODO(), types.NamespacedName{Name: podName, Namespace: instance.Namespace}, pod)
	//	if err != nil && errors.IsNotFound(err) {
	//		reqLogger.Info(time.Now().String() + "================== not find pod:  " + podName)
	//		reqLogger.Info(time.Now().String() + "++++++++++++++++++statefulSet get pod: " + err.Error())
	//	}
	//	if pod.DeletionTimestamp != nil {
	//		var GracePeriodSeconds int64 = 0
	//		reqLogger.Info(time.Now().String() + "==================statefulSet pod deletion: " + pod.DeletionTimestamp.String())
	//		err = r.client.Delete(context.TODO(), pod, &client.DeleteOptions{GracePeriodSeconds: &GracePeriodSeconds})
	//		if err != nil {
	//			reqLogger.Info(time.Now().String() + "==================delete pod error : " + err.Error())
	//		}
	//	}
	//	reqLogger.Info(time.Now().String() + "==================statefulSet pod status: " + string(pod.Status.Phase))
	//	reqLogger.Info(time.Now().String() + "==================statefulSet pod: " + string(pod.String()))
	//}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	//if err != nil && errors.IsNotFound(err) {
	//	reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
	//	err = r.client.Create(context.TODO(), pod)
	//	if err != nil {
	//		return reconcile.Result{}, err
	//	}
	//
	//	// Pod created successfully - don't requeue
	//	return reconcile.Result{}, nil
	//} else if err != nil {
	//	return reconcile.Result{}, err
	//}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)
	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *apiapps.StatefulSet) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    "busybox",
					Image:   "busybox",
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

// getPodNames returns the pod names of the array of pods passed in
func getPodNames(cr *apiapps.StatefulSet) []string {
	var podNames []string
	var sum int = int(cr.Status.Replicas)
	for i := 0; i < sum; i++ {
		podNames = append(podNames, cr.Name + "-" + strconv.Itoa(i))
	}
	return podNames
}
