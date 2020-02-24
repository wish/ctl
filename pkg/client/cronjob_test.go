package client

import (
	"strconv"
	"testing"

	batchv1b1 "k8s.io/api/batch/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func getRandomCronJobs(N int) []*batchv1b1.CronJob {
	cronjobs := make([]*batchv1b1.CronJob, N)
	for n := 0; n < N; n++ {
		cronjobs[n] = &batchv1b1.CronJob{
			TypeMeta: metav1.TypeMeta{
				Kind:       "cronjob",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
	}
	return cronjobs
}

func getRandomCronJobsObject(N int) []runtime.Object {
	cronjobs := make([]runtime.Object, N)
	for n := 0; n < N; n++ {
		temp := &batchv1b1.CronJob{
			TypeMeta: metav1.TypeMeta{
				Kind:       "cronjob",
				APIVersion: "v1",
			},
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(n),
				Namespace: strconv.Itoa(n),
			},
		}
		cronjobs[n] = temp.DeepCopyObject()
	}
	return cronjobs
}

func TestListCronJobsSingle(t *testing.T) {
	cronjob := batchv1b1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "cronjob",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {cronjob.DeepCopyObject()}})

	p, err := cl.ListCronJobs("hi", "", ListOptions{})
	if err != nil {
		t.Error(err.Error())
	}

	t.Log(p)

	if len(p) != 1 {
		t.Errorf("Unexpected number of cronjobs found! %d != 1; expected", len(p))
	}
}

func TestListCronJobsBadContext(t *testing.T) {
	cronjob := batchv1b1.CronJob{
		TypeMeta: metav1.TypeMeta{
			Kind:       "cronjob",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
		},
	}

	cl := GetFakeConfigClient(map[string][]runtime.Object{"hi": {cronjob.DeepCopyObject()}})

	p, err := cl.ListCronJobs("pew", "", ListOptions{})

	if len(p) == 0 && err != nil {
		t.Log("Error as expected:", err.Error())
	} else {
		t.Error("Context not found did not error")
	}
}

func TestListCronJobsMultiple(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6)})

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
		p, err := cl.ListCronJobs(q.context, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of cronjobs found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestListCronJobsOverContexts(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6), "c3": getRandomCronJobsObject(3), "c4": nil})

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
		p, err := cl.ListCronJobsOverContexts(q.contexts, q.namespace, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of cronjobs found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestGetCronJob(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6)}
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
		p, err := cl.GetCronJob(q.context, q.namespace, q.name, GetOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if p.Name != q.name {
			t.Error("Incorrect cronjob found")
		}
	}
}

func TestGetCronJobBadContext(t *testing.T) {
	vals := map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6)}
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
		c, err := cl.GetCronJob(q.context, q.namespace, q.name, GetOptions{})
		if err != nil && c == nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Cronjob get did not error!")
		}
	}
}

func TestFindCronJobs(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6), "c3": getRandomCronJobsObject(3), "c4": nil})

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
		p, err := cl.FindCronJobs(q.contexts, q.namespace, q.names, ListOptions{})
		if err != nil {
			t.Error(err.Error())
		}

		if len(p) != q.size {
			t.Errorf("Unexpected number of cronjobs found! %d != %d; expected", len(p), q.size)
		}
	}
}

func TestFindCronJobsError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6), "c3": getRandomCronJobsObject(3)})

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
		_, err := cl.FindCronJobs(q.contexts, q.namespace, q.names, ListOptions{})
		if err == nil {
			t.Error("FindCronJobs did not fail!")
		} else {
			t.Log("FindCronJobs failed as expected: ", err.Error())
		}
	}
}

func TestFindCronJob(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6), "c3": getRandomCronJobsObject(3)})
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
		c, err := cl.findCronJob(q.contexts, q.namespace, q.name, ListOptions{})

		if c == nil || err != nil {
			t.Error("Could not find pod with error:", err)
		} else if c.Name != q.name {
			t.Errorf("The found pod does not match the name requested: %s != %s", c.Name, q.name)
		}
	}
}

func TestFindCronJobError(t *testing.T) {
	cl := GetFakeConfigClient(map[string][]runtime.Object{"c1": getRandomCronJobsObject(5), "c2": getRandomCronJobsObject(6), "c3": getRandomCronJobsObject(3)})
	var queries = []struct {
		contexts  []string
		namespace string
		name      string
	}{
		{[]string{"c2", "c3"}, "", "10"},
		{[]string{"c4"}, "3", "3"},
	}

	for _, q := range queries {
		_, err := cl.findCronJob(q.contexts, q.namespace, q.name, ListOptions{})

		if err != nil {
			t.Log("Errored as expected:", err)
		} else {
			t.Error("Function did not error when finding log")
		}
	}
}
