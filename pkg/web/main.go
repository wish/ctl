package web

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"html/template"
	"k8s.io/client-go/rest"
	"net/http"
	"regexp"
)

func Serve(endpoint string) {
	cl := client.GetDefaultConfigClient()

	// Main page
	templates := template.Must(template.ParseFiles("pkg/web/template/dash.html", "pkg/web/template/details.html", "pkg/web/template/run.html"))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		cronjobs, err := cl.ListCronJobsOverContexts([]string{}, "", client.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		runs, err := cl.ListRunsOverContexts(nil, "", client.ListOptions{})
		if err != nil {
			panic(err.Error())
		}

		if err := templates.ExecuteTemplate(w, "dash.html", toCardDetailsList(cronjobs, runs)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Cron job detail pages
	http.HandleFunc("/cronjob/", func(w http.ResponseWriter, r *http.Request) {
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

		if err := templates.ExecuteTemplate(w, "details.html", toFullDetails(cronjob, runs)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	// Execute cron job
	http.HandleFunc("/execute/", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/execute/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		_, err := cl.RunCronJob([]string{path[1]}, path[2], path[3])

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	})

	// Suspend cron job
	http.HandleFunc("/suspend/", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/suspend/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		success, err := cl.SetCronJobSuspend([]string{path[1]}, path[2], path[3], true)

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if !success {
			http.Error(w, "Invalid request! Cron job already suspended.", http.StatusConflict)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	})

	// Unsuspend cron job
	http.HandleFunc("/unsuspend/", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/unsuspend/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 4 {
			http.Error(w, "Invalid cron job execute path", http.StatusNotFound)
		}

		success, err := cl.SetCronJobSuspend([]string{path[1]}, path[2], path[3], false)

		if err != nil { // Could not execute
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else if !success {
			http.Error(w, "Invalid request! Cron job already unsuspended.", http.StatusConflict)
		} else {
			http.Redirect(w, r, fmt.Sprintf("/cronjob/%s/%s/%s", path[1], path[2], path[3]), http.StatusSeeOther)
		}
	})

	// Specific run page
	http.HandleFunc("/run/", func(w http.ResponseWriter, r *http.Request) {
		re := regexp.MustCompile(`/run/([^/]*)/([^/]*)/([^/]*)/([^/]*)\z`)

		path := re.FindStringSubmatch(r.URL.String())
		if !re.MatchString(r.URL.String()) || len(path) != 5 {
			http.Error(w, "Invalid cron job path", http.StatusNotFound)
		}

		runs, err := cl.ListRunsOfCronJob([]string{path[1]}, path[2], path[3], client.ListOptions{})
		if err != nil {
			http.Error(w, "Error finding runs", http.StatusInternalServerError)
		}

		var run *client.RunDiscovery
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

		logs := make(map[string]*rest.Result)

		for _, pod := range pods {
			res, err := cl.LogPod(pod.Context, pod.Namespace, pod.Name, "", client.LogOptions{})
			if err != nil {
				fmt.Println(err.Error())
				continue // hmm
			}
			logs[pod.Name] = res
		}

		if err := templates.ExecuteTemplate(w, "run.html", toFullRunDetails(path, run, logs)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Listening on", endpoint)
	fmt.Println(http.ListenAndServe(endpoint, nil))
}
