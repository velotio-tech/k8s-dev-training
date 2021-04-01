package jobs

import (
	"context"
	qav1 "github.com/farkaskid/k8s-dev-training/assignment2/api/v1"
	"github.com/go-logr/logr"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/api/batch/v1beta1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"math/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"strconv"
	"time"
)

type JobObject interface {
	metav1.Object
	runtime.Object
}

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
	sanity *qav1.CodeSanity,
	log logr.Logger,
) error {
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)

	log.Info("creating testing job for pod: " + pod.Name)

	var job JobObject

	if sanity.Spec.Resource == qav1.CronJob {
		job = createCronJob(timeStr, pod, "")
	} else {
		job = createJob(timeStr, pod)
	}

	if err := ctrl.SetControllerReference(sanity, job, scheme); err != nil {
		log.Error(err, "Failed to attach owner reference to the job")
	}

	return client.Create(ctx, job)
}

func createJob(timeStr string, pod *v1.Pod) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{"type": "testing", "pod": pod.Name},
			Annotations: map[string]string{},
			Name:        pod.Name + "-testing-job-" + timeStr,
			Namespace:   pod.Namespace,
		},
		Spec: createJobSpec(pod, timeStr),
	}
}

func createCronJob(timeStr string, pod *v1.Pod, schedule string) *v1beta1.CronJob {
	var jobSchedule string

	if len(schedule) == 0 {
		jobSchedule = "* * * * *"
	} else {
		jobSchedule = schedule
	}

	return &v1beta1.CronJob{
		ObjectMeta: metav1.ObjectMeta{
			Labels:      map[string]string{"type": "testing", "pod": pod.Name},
			Annotations: map[string]string{},
			Name:        pod.Name + "-testing-job-" + timeStr,
			Namespace:   pod.Namespace,
		},
		Spec: v1beta1.CronJobSpec{
			Schedule:                   jobSchedule,
			FailedJobsHistoryLimit:     new(int32),
			SuccessfulJobsHistoryLimit: new(int32),
			JobTemplate: v1beta1.JobTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{Labels: map[string]string{"pod": pod.Name}},
				Spec:       createJobSpec(pod, timeStr),
			},
		},
	}
}

func createJobSpec(pod *v1.Pod, timeStr string) batchv1.JobSpec {
	return batchv1.JobSpec{
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
	}
}

func IsJobFinished(job *batchv1.Job) (bool, batchv1.JobConditionType) {
	for _, condition := range job.Status.Conditions {
		if (condition.Type == batchv1.JobComplete || condition.Type == batchv1.JobFailed) && condition.Status == v1.ConditionTrue {
			// A Random number decides if the job failed or not to emulate failures
			if rand.Intn(2) == 1 {
				return true, batchv1.JobComplete
			}

			return true, batchv1.JobFailed
		}
	}

	return false, ""
}
