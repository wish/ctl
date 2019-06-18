package web

import (
	"encoding/json"
	"github.com/ContextLogic/ctl/pkg/client"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sort"
)

// Processes objects to create easy to print objects

type cardDetails struct {
	Name      string
	Context   string
	Namespace string
	Active    int
	Suspend   bool
	LastRun   RunStatus
}

type RunStatus string

const (
	RunFailed  RunStatus = "failed"
	RunSuccess RunStatus = "success"
	RunNA      RunStatus = "N/A"
	RunRunning RunStatus = "running"
)

func toCardDetails(c *client.CronJobDiscovery, r *client.RunDiscovery) cardDetails {
	var runStatus RunStatus
	if r == nil {
		runStatus = RunNA
	} else if r.Status.Failed > 0 {
		runStatus = RunFailed
	} else if r.Status.CompletionTime != nil {
		runStatus = RunSuccess
	} else {
		runStatus = RunRunning
	}

	return cardDetails{
		Name:      c.Name,
		Context:   c.Context,
		Namespace: c.Namespace,
		Active:    len(c.Status.Active),
		Suspend:   *(c.Spec.Suspend),
		LastRun:   runStatus,
	}
}

func toCardDetailsList(lst []client.CronJobDiscovery, runs []client.RunDiscovery) []cardDetails {
	ret := make([]cardDetails, len(lst))
	recent := make(map[types.UID]*client.RunDiscovery)

	for i, _ := range runs {
		if len(runs[i].OwnerReferences) == 1 {
			if x, ok := recent[runs[i].OwnerReferences[0].UID]; !ok || runs[i].Status.StartTime.After(x.Status.StartTime.Time) {
				recent[runs[i].OwnerReferences[0].UID] = &runs[i]
			}
		}
	}

	for i, _ := range lst {
		ret[i] = toCardDetails(&lst[i], recent[lst[i].UID])
	}
	return ret
}

type fullDetails struct {
	Name         string
	Context      string
	Namespace    string
	Schedule     string
	Template     string
	LastSchedule string
	Runs         []runDetails
}

type runDetails struct {
	Name   string
	Start  string
	Status string
	End    string
}

type byStartTime []client.RunDiscovery

func (l byStartTime) Len() int      { return len(l) }
func (l byStartTime) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byStartTime) Less(i, j int) bool {
	return l[i].Status.StartTime.After(l[j].Status.StartTime.Time)
}

func toFullDetails(cronjob *client.CronJobDiscovery, runs []client.RunDiscovery) fullDetails {
	sort.Sort(byStartTime(runs))

	// Last schedule
	lastSchedule := "N/A"
	if cronjob.Status.LastScheduleTime != nil {
		lastSchedule = cronjob.Status.LastScheduleTime.Format("Mon Jan _2 3:04pm 2006")
	}

	// Template
	// b, err := cronjob.Spec.JobTemplate.Marshal()
	template, _ := json.MarshalIndent(cronjob.Spec.JobTemplate, "", "  ")

	return fullDetails{
		Name:         cronjob.Name,
		Context:      cronjob.Context,
		Namespace:    cronjob.Namespace,
		Schedule:     cronjob.Spec.Schedule,
		Template:     string(template), // Change
		LastSchedule: lastSchedule,
		Runs:         toRunDetailsList(runs),
	}
}

func toRunDetails(run client.RunDiscovery) runDetails {
	// Get condition:
	condition := "Running"
	for _, x := range run.Status.Conditions {
		if x.Status == corev1.ConditionTrue {
			condition = string(x.Type)
		}
	}

	end := "N/A"
	if run.Status.CompletionTime != nil {
		end = run.Status.CompletionTime.Format("Mon Jan _2 3:04pm 2006")
	}

	return runDetails{
		Name:   run.Name,
		Start:  run.Status.StartTime.Format("Mon Jan _2 3:04pm 2006"),
		Status: condition,
		End:    end,
	}
}

func toRunDetailsList(lst []client.RunDiscovery) []runDetails {
	ret := make([]runDetails, len(lst))

	for i := range lst {
		ret[i] = toRunDetails(lst[i])
	}

	return ret
}

//
// type fullRunDetails struct {
//
// }
