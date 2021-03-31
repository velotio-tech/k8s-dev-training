package v1

import v1 "k8s.io/api/core/v1"

type EventType string

const (
	Create EventType = "Create"
	Update EventType = "Update"
)

// Return the event type for the CodeSanity CR
func (sanity *CodeSanity) GetEventType() EventType {
	if len(sanity.Status.UnhealthyPods) == 0 && len(sanity.Status.HealthyPods) == 0 {
		return Create
	}
	return Update
}

func (sanity *CodeSanity) PodValid(pod v1.Pod) bool {
	for _, podName := range sanity.Status.ProcessedPods {
		if pod.Name == podName {
			return false
		}
	}

	for _, podName := range sanity.Spec.PodNames {
		if pod.Name == podName {
			return true
		}
	}
	return false
}
