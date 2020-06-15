package filter

import v1 "k8s.io/api/core/v1"

// StatusMatch is used to filter pods by status
type StatusMatch struct {
	State v1.PodPhase
}
