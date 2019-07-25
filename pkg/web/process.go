package web

import (
	"encoding/json"
	"fmt"
	dtypes "github.com/wish/ctl/pkg/client/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	"sort"
)

// Processes objects to create easy to print objects

type page struct {
	Title  string
	Active string
}

type cardDetails struct {
	Name      string
	Context   string
	Namespace string
	Active    int
	Suspend   bool
	LastRun   runStatus
}

//
type runStatus string

const (
	runFailed  runStatus = "failed"
	runSuccess runStatus = "success"
	runNA      runStatus = "N/A"
	runRunning runStatus = "running"
)

func toCardDetails(c *dtypes.CronJobDiscovery, r *dtypes.RunDiscovery) cardDetails {
	var runStatus runStatus
	if r == nil {
		runStatus = runNA
	} else if r.Status.Failed > 0 {
		runStatus = runFailed
	} else if r.Status.CompletionTime != nil {
		runStatus = runSuccess
	} else {
		runStatus = runRunning
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

func toCardDetailsList(lst []dtypes.CronJobDiscovery, runs []dtypes.RunDiscovery) []cardDetails {
	ret := make([]cardDetails, len(lst))
	recent := make(map[types.UID]*dtypes.RunDiscovery)

	for i := range runs {
		if len(runs[i].OwnerReferences) == 1 {
			if x, ok := recent[runs[i].OwnerReferences[0].UID]; !ok || runs[i].Status.StartTime.After(x.Status.StartTime.Time) {
				recent[runs[i].OwnerReferences[0].UID] = &runs[i]
			}
		}
	}

	for i := range lst {
		ret[i] = toCardDetails(&lst[i], recent[lst[i].UID])
	}
	return ret
}

type fullDetails struct {
	Name         string
	Context      string
	Namespace    string
	Schedule     string
	Suspend      bool
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

type byStartTime []dtypes.RunDiscovery

func (l byStartTime) Len() int      { return len(l) }
func (l byStartTime) Swap(i, j int) { l[i], l[j] = l[j], l[i] }
func (l byStartTime) Less(i, j int) bool {
	return l[i].Status.StartTime.After(l[j].Status.StartTime.Time)
}

func toFullDetails(cronjob *dtypes.CronJobDiscovery, runs []dtypes.RunDiscovery) fullDetails {
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
		Suspend:      *(cronjob.Spec.Suspend),
		Schedule:     cronjob.Spec.Schedule,
		Template:     string(template), // Change
		LastSchedule: lastSchedule,
		Runs:         toRunDetailsList(runs),
	}
}

func toRunDetails(run dtypes.RunDiscovery) runDetails {
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

func toRunDetailsList(lst []dtypes.RunDiscovery) []runDetails {
	ret := make([]runDetails, len(lst))

	for i := range lst {
		ret[i] = toRunDetails(lst[i])
	}

	return ret
}

type fullRunDetails struct {
	Name      string
	Context   string
	Namespace string
	Cronjob   string
	Start     string
	Status    string
	End       string
	Pods      []runPodDetails
}

type runPodDetails struct {
	Name string
	Logs string
}

func toFullRunDetails(path []string, run dtypes.RunDiscovery, logs map[string]rest.Result) fullRunDetails {
	pods := make([]runPodDetails, 0, len(logs))

	for p, r := range logs {
		raw, err := r.Raw()
		if err != nil {
			fmt.Println(err.Error())
			continue
		}
		pods = append(pods, runPodDetails{p, string(raw)})
	}

	details := toRunDetails(run)

	return fullRunDetails{
		Name:      run.Name,
		Cronjob:   path[3],
		Context:   path[1],
		Namespace: path[2],
		Start:     details.Start,
		Status:    details.Status,
		End:       details.End,
		Pods:      pods,
	}
}
