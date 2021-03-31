package jobs

import (
	"context"
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"
)

func JobRequestHandler(
	job batchv1.Job, client client.Client, log logr.Logger) (ctrl.Result, error) {
	log.Info("job handler called with " + job.Name)

	return ctrl.Result{}, nil
}

func CreateJobForPod(
	ctx context.Context,
	pod *v1.Pod,
	client client.Client,
	scheme *runtime.Scheme,
	sanity qav1.CodeSanity,
	log logr.Logger,
) error {
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	log.Info("creating testing job for pod" + pod.Name)

	job := batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{"type": "testing"},
			Annotations: map[string]string{},
			Name:        pod.Name + "-testing-job-" + timeStr,
			Namespace:   pod.Namespace,
		},
		Spec: batchv1.JobSpec{
			Template: v1.PodTemplateSpec{
				Spec: v1.PodSpec{
					Containers: []v1.Container{{
						Name:    pod.Name + "-testing-job-container" + timeStr,
						Image:   "busybox",
						Command: []string{"/bin/sh", "-ec", "sleep 5"},
					}},
					RestartPolicy: v1.RestartPolicyNever,
				},
			},
		},
	}

	if err := ctrl.SetControllerReference(&sanity, &job, scheme); err != nil {
		log.Error(err, "Failed to attach owner reference to the job")
	}

	return client.Create(ctx, &job)
}
