package web

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/wish/ctl/cmd/util/parsing"
	"github.com/wish/ctl/pkg/client"
	"github.com/wish/ctl/pkg/client/types"
	"html/template"
	"k8s.io/client-go/rest"
	"net/http"
	"regexp"
	"strings"
)

// Serve runs a webserver for kron at the specified url
func Serve(cl *client.Client, endpoint string) {
	templates := template.Must(template.New("all").Funcs(template.FuncMap{
		"Title": strings.Title,
	}).ParseGlob("pkg/web/template/*"))

	// templates = templates.Funcs(template.FuncMap{
	// 	"Title": strings.Title,
	// })

	// Main page
	http.HandleFunc("/", getDashHandleFunc(cl, templates))

	// Advanced dashboard
	http.HandleFunc("/advanced", getAdvDashHandleFunc(cl, templates))

	// Cron job detail pages
	http.HandleFunc("/cronjob/", getDetailsHandleFunc(cl, templates))

	// Execute cron job
	http.HandleFunc("/execute/", getExecuteHandleFunc(cl, templates))

	// Suspend cron job
	http.HandleFunc("/suspend/", getSuspendHandleFunc(cl, templates))

	// Unsuspend cron job
	http.HandleFunc("/unsuspend/", getUnsuspendHandleFunc(cl, templates))

	// Specific run page
	http.HandleFunc("/run/", getRunHandleFunc(cl, templates))

	fmt.Println("Listening on", endpoint)
	fmt.Println(http.ListenAndServe(endpoint, nil))
}

func getDashHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Page           page
			Contexts       map[string]bool
			Namespace      string
			Search         string
			Cronjobs       []cardDetails
			DefaultColumns map[string]string
		}{
			Page:           page{Title: "Dashboard - Kron", Active: "dashboard"},
			Contexts:       make(map[string]bool),
			DefaultColumns: make(map[string]string),
		}

		// Contexts
		for _, x := range cl.GetAllContexts() {
			data.Contexts[x] = false
		}

		// Check valid
		ctxs := r.URL.Query()["context"]
		for _, c := range ctxs {
			if _, ok := data.Contexts[c]; ok {
				data.Contexts[c] = true
			}
		}

		// Default columns
		var dcLabelsRaw []string

		for _, col := range viper.GetStringSlice("default_columns") {
			v := r.URL.Query().Get(col)
			data.DefaultColumns[col] = v
			if len(v) > 0 {
				dcLabelsRaw = append(dcLabelsRaw, col+"="+v)
			}
		}

		labelMatch, err := parsing.LabelMatchSlice(dcLabelsRaw)
		if err != nil {
			fmt.Println("error encountered when parsing default columns: ", err.Error())
		}

		namespace := r.URL.Query().Get("namespace")
		data.Namespace = namespace

		search := r.URL.Query().Get("search")
		data.Search = search

		cronjobs, err := cl.ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{LabelMatch: labelMatch})
		if err != nil {
			panic(err.Error())
		}

		// Filter searches
		var filtered []types.CronJobDiscovery
		if search == "" {
			filtered = cronjobs
		} else {
			for _, c := range cronjobs {
				if strings.Contains(c.Name, search) {
					filtered = append(filtered, c)
				}
			}
		}

		runs, err := cl.ListRunsOverContexts(ctxs, namespace, client.ListOptions{LabelMatch: labelMatch})
		if err != nil {
			panic(err.Error())
		}

		data.Cronjobs = toCardDetailsList(filtered, runs)

		if err := templates.ExecuteTemplate(w, "dash.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getAdvDashHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data := struct {
			Page           page
			Contexts       map[string]bool
			Namespace      string
			Search         string
			Cronjobs       []cardDetails
			DefaultColumns map[string]string
		}{
			Page:           page{Title: "Advanced Dashboard - Kron", Active: "advanced"},
			Contexts:       make(map[string]bool),
			DefaultColumns: make(map[string]string),
		}

		// Contexts
		for _, x := range cl.GetAllContexts() {
			data.Contexts[x] = false
		}

		// Check valid
		ctxs := r.URL.Query()["context"]
		for _, c := range ctxs {
			if _, ok := data.Contexts[c]; ok {
				data.Contexts[c] = true
			}
		}

		// Default columns
		var dcLabelsRaw []string

		for _, col := range viper.GetStringSlice("default_columns") {
			v := r.URL.Query().Get(col)
			data.DefaultColumns[col] = v
			if len(v) > 0 {
				dcLabelsRaw = append(dcLabelsRaw, col+"="+v)
			}
		}

		labelMatch, err := parsing.LabelMatchSlice(dcLabelsRaw)
		if err != nil {
			fmt.Println("error encountered when parsing default columns: ", err.Error())
		}

		namespace := r.URL.Query().Get("namespace")
		data.Namespace = namespace

		search := r.URL.Query().Get("search")
		data.Search = search

		cronjobs, err := cl.ListCronJobsOverContexts(ctxs, namespace, client.ListOptions{LabelMatch: labelMatch})
		if err != nil {
			panic(err.Error())
		}

		// Filter searches
		var filtered []types.CronJobDiscovery
		if search == "" {
			filtered = cronjobs
		} else {
			for _, c := range cronjobs {
				if strings.Contains(c.Name, search) {
					filtered = append(filtered, c)
				}
			}
		}

		runs, err := cl.ListRunsOverContexts(ctxs, namespace, client.ListOptions{LabelMatch: labelMatch})
		if err != nil {
			panic(err.Error())
		}

		data.Cronjobs = toCardDetailsList(filtered, runs)

		if err := templates.ExecuteTemplate(w, "advanced.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getDetailsHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/cronjob/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job path", http.StatusNotFound)
		}

		cronjob, err := cl.GetCronJob(path[1], path[2], path[3], client.GetOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		runs, err := cl.ListRunsOfCronJob([]string{path[1]}, path[2], path[3], client.ListOptions{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		data := struct {
			Page    page
			Details fullDetails
		}{
			Page: page{
				Title: fmt.Sprintf("%s - Cron Jobs - Kron", cronjob.Name),
			},
			Details: toFullDetails(cronjob, runs),
		}

		if err := templates.ExecuteTemplate(w, "details.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getExecuteHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/execute/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		_, err := cl.RunCronJob([]string{path[1]}, path[2], path[3], client.ListOptions{})

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	}
}

func getSuspendHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/suspend/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		success, err := cl.SetCronJobSuspend([]string{path[1]}, path[2], path[3], true, client.ListOptions{})

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if !success {
			http.Error(w, "Invalid request! Cron job already suspended.", http.StatusConflict)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	}
}

func getUnsuspendHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/unsuspend/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		success, err := cl.SetCronJobSuspend([]string{path[1]}, path[2], path[3], false, client.ListOptions{})

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if !success {
			http.Error(w, "Invalid request! Cron job already unsuspended.", http.StatusConflict)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	}
}

func getRunHandleFunc(cl *client.Client, templates *template.Template) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/run/([^/]*)/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 5 {
			http.Error(w, "Invalid cron job path", http.StatusNotFound)
		}

		runs, err := cl.ListRunsOfCronJob([]string{path[1]}, path[2], path[3], client.ListOptions{})
		if err != nil {
			http.Error(w, "Error finding runs", http.StatusInternalServerError)
		}

		var run *types.RunDiscovery
		for _, x := range runs {
			if x.Name == path[4] {
				run = &x
				break
			}
		}

		if run == nil {
			http.Error(w, "Specified run could not be found!", http.StatusNotFound)
		}

		pods, err := cl.ListPodsOfRun([]string{path[1]}, path[2], path[4], client.ListOptions{})
		if err != nil {
			http.Error(w, "Error finding pods of run", http.StatusInternalServerError)
		}

		logs := make(map[string]rest.Result)

		for _, pod := range pods {
			req, err := cl.LogPod(pod.Context, pod.Namespace, pod.Name, "", client.LogOptions{})
			if err != nil {
				fmt.Println(err.Error())
				continue // hmm
			}
			logs[pod.Name] = req.Do()
		}

		data := struct {
			Page    page
			Details fullRunDetails
		}{
			Page: page{
				Title: fmt.Sprintf("%s - %s Runs - Kron", run.Name, path[3]),
			},
			Details: toFullRunDetails(path, *run, logs),
		}

		if err := templates.ExecuteTemplate(w, "run.html", data); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
