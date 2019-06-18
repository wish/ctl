package web

import (
	"fmt"
	"github.com/ContextLogic/ctl/pkg/client"
	"html/template"
	"net/http"
	"regexp"
)

func Serve(endpoint string) {
	cl := client.GetDefaultConfigClient()

	// Main page
	templates := template.Must(template.ParseFiles("pkg/web/template/dash.html", "pkg/web/template/details.html"))

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
			fmt.Println("Error")
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

	fmt.Println("Listening on", endpoint)
	fmt.Println(http.ListenAndServe(endpoint, nil))
}
