package client

import (
	"strconv"
	"testing"

	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getRandomJobs(N int) []*batchv1.Job {
	jobs := make([]*batchv1.Job, N)
	for n := 0; n < N; n++ {
		jobs[n] = &batchv1.Job{
			TypeMeta: metav1.TypeMeta{
				Kind:       "job",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
	}
	return jobs
}

func getRandomJobsObject(N int) []runtime.Object {
	jobs := make([]runtime.Object, N)
	for n := 0; n < N; n++ {
		temp := &batchv1.Job{
			TypeMeta: metav1.TypeMeta{
				Kind:       "job",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
		jobs[n] = temp.DeepCopyObject()
	}
	return jobs
}

func TestListJobsSingle(t *testing.T) {
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {job.DeepCopyObject()}})

	l, err := cl.ListJobs("hi", "", ListOptions{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(l)

	if len(l) != 1 {
		t.Errorf("Unexpected number of jobs found! %d != 1; expected", len(l))
	}
}

func TestListJobsBadContext(t *testing.T) {
	job := batchv1.Job{
		TypeMeta: metav1.TypeMeta{
			Kind:       "job",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {job.DeepCopyObject()}})

	l, err := cl.ListJobs("pew", "", ListOptions{})

	if len(l) == 0 && err != nil {
		t.Log("Error as expected:", err.Error())
	} else {
		t.Error("Context not found did not error")
	}
}

func TestListJobsMultiple(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6)})

	var queries = []struct {
		context   string
		namespace string
		size      int
	}{
		{"c1", "", 5},
		{"c2", "", 6},
		{"c1", "0", 1},
		{"c2", "1", 1},
	}

	for _, q := range queries {
		l, err := cl.ListJobs(q.context, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(l) != q.size {
			t.Errorf("Unexpected number of jobs found! %d != %d; expected", len(l), q.size)
		}
	}
}

func TestListJobsOverContexts(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6), "c3": getRandomJobsObject(3), "c4": nil})

	var queries = []struct {
		contexts  []string
		namespace string
		size      int
	}{
		{[]string{"c1"}, "", 5},
		{[]string{"c2"}, "", 6},
		{[]string{"c3"}, "", 3},
		{nil, "", 14},
		{[]string{"c1", "c2", "c3"}, "", 14},
		{[]string{"c1", "c2"}, "5", 1},
		{nil, "4", 2},
		{[]string{}, "0", 3},
		{[]string{"c1", "c2"}, "", 11},
		{[]string{"c3", "c1"}, "", 8},
	}

	for _, q := range queries {
		l, err := cl.ListJobsOverContexts(q.contexts, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(l) != q.size {
			t.Errorf("Unexpected number of jobs found! %d != %d; expected", len(l), q.size)
		}
	}
}

func TestGetJob(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6)}
	cl := GetFakeConfigClient(vals)

	var queries = []struct {
		context   string
		namespace string
		name      string
	}{
		{"c1", "1", "1"},
		{"c2", "0", "0"},
		{"c1", "0", "0"},
		{"c2", "5", "5"},
	}

	for _, q := range queries {
		j, err := cl.GetJob(q.context, q.namespace, q.name, GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if j.Name != q.name {
			t.Error("Incorrect job found")
		}
	}
}

func TestGetJobBadContext(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6)}
	cl := GetFakeConfigClient(vals)

	var queries = []struct {
		context   string
		namespace string
		name      string
	}{
		{"c3", "", "1"},
		{"c3", "0", "0"},
		{"c3", "", "0"},
		{"c3", "5", "5"},
	}

	for _, q := range queries {
		j, err := cl.GetJob(q.context, q.namespace, q.name, GetOptions{})
		if err != nil && j == nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Cronjob get did not error!")
		}
	}
}

func TestFindJobs(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6), "c3": getRandomJobsObject(3), "c4": nil})

	var queries = []struct {
		contexts  []string
		namespace string
		names     []string
		size      int
	}{
		{[]string{"c1"}, "", []string{"0"}, 1},
		{nil, "", []string{"0", "3", "5"}, 6},
		{nil, "4", []string{"0", "3", "5"}, 0},
		{nil, "4", []string{"4"}, 2},
		{[]string{"c2"}, "4", []string{"4"}, 1},
		{[]string{"c1", "c3", "c2"}, "4", []string{"4"}, 2},
		{[]string{"c3"}, "", []string{"3"}, 0},
	}

	for _, q := range queries {
		j, err := cl.FindJobs(q.contexts, q.namespace, q.names, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(j) != q.size {
			t.Errorf("Unexpected number of jobs found! %d != %d; expected", len(j), q.size)
		}
	}
}

func TestFindJobsError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6), "c3": getRandomJobsObject(3)})

	var queries = []struct {
		contexts  []string
		namespace string
		names     []string
	}{
		{[]string{"c4"}, "", []string{"0"}},
		{[]string{"c4"}, "", []string{"0", "3", "5"}},
		{[]string{"c4"}, "4", []string{"0", "3", "5"}},
		{[]string{"c4"}, "4", []string{"4"}},
		{[]string{"c2", "c4"}, "4", []string{"4"}},
		{[]string{"c4", "c1", "c3", "c2"}, "4", []string{"4"}},
		{[]string{"c3", "c4"}, "", []string{"3"}},
	}

	for _, q := range queries {
		_, err := cl.FindJobs(q.contexts, q.namespace, q.names, ListOptions{})
		if err == nil {
			t.Error("FindJobs did not fail!")
		} else {
			t.Log("FindJobs failed as expected: ", err.Error())
		}
	}
}

func TestFindJob(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6), "c3": getRandomJobsObject(3)})
	var queries = []struct {
		contexts  []string
		namespace string
		name      string
	}{
		{nil, "", "1"},
		{nil, "5", "5"},
		{[]string{"c2", "c3"}, "", "2"},
		{[]string{"c2", "c3"}, "3", "3"},
		{[]string{}, "", "2"},
	}

	for _, q := range queries {
		j, err := cl.findJob(q.contexts, q.namespace, q.name, ListOptions{})

		if j == nil || err != nil {
			t.Error("Could not find pod with error:", err)
		} else if j.Name != q.name {
			t.Errorf("The found pod does not match the name requested: %s != %s", j.Name, q.name)
		}
	}
}

func TestFindJobError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomJobsObject(5), "c2": getRandomJobsObject(6), "c3": getRandomJobsObject(3)})
	var queries = []struct {
		contexts  []string
		namespace string
		name      string
	}{
		{[]string{"c2", "c3"}, "", "10"},
		{[]string{"c4"}, "3", "3"},
	}

	for _, q := range queries {
		_, err := cl.findJob(q.contexts, q.namespace, q.name, ListOptions{})

		if err != nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Function did not error when finding log")
		}
	}
}
