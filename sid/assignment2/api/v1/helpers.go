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

// Mark the pod with podName as healthy
func (sanity *CodeSanity) AddHealthyPod(healthyPodName string) {

	// Delete Pod from unhealthy list
	podUnhealthyIndex := -1

	for index, podName := range sanity.Status.UnhealthyPods {
		if podName == healthyPodName {
			podUnhealthyIndex = index
			break
		}
	}

	if podUnhealthyIndex != -1 {
		sanity.Status.UnhealthyPods = append(
			sanity.Status.UnhealthyPods[:podUnhealthyIndex], sanity.Status.UnhealthyPods[podUnhealthyIndex+1:]...)
	}

	// Add Pod to healthy list if it's not already in that
	podHealthyIndex := -1

	for index, podName := range sanity.Status.HealthyPods {
		if podName == healthyPodName {
			podHealthyIndex = index
			break
		}
	}

	if podHealthyIndex == -1 {
		sanity.Status.HealthyPods = append(sanity.Status.HealthyPods, healthyPodName)
	}
}

// Mark the pod with podName as healthy
func (sanity *CodeSanity) AddUnhealthyPod(unhealthyPodName string) {

	// Delete Pod from healthy list
	podHealthyIndex := -1

	for index, podName := range sanity.Status.HealthyPods {
		if podName == unhealthyPodName {
			podHealthyIndex = index
			break
		}
	}

	if podHealthyIndex != -1 {
		sanity.Status.HealthyPods = append(
			sanity.Status.HealthyPods[:podHealthyIndex], sanity.Status.HealthyPods[podHealthyIndex+1:]...)
	}

	// Add Pod to unhealthy list if it's not already in that
	podUnhealthyIndex := -1

	for index, podName := range sanity.Status.UnhealthyPods {
		if podName == unhealthyPodName {
			podUnhealthyIndex = index
			break
		}
	}

	if podUnhealthyIndex == -1 {
		sanity.Status.UnhealthyPods = append(sanity.Status.UnhealthyPods, unhealthyPodName)
	}
}
