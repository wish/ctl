package client

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// SetCronJobSuspend sets the cron job suspend state to given value and reports whether successful.
// Does nothing (and returns false) if the cronjob is already in the correct state.
func (c *Client) SetCronJobSuspend(contextStr, namespace, name string, suspend bool) (bool, error) {
	cronjob, err := c.findExactCronJob(contextStr, namespace, name)
	if err != nil {
		return false, err
	}
	if *cronjob.Spec.Suspend == suspend { // Already set to value
		return false, nil
	}
	*cronjob.Spec.Suspend = suspend

	cl, err := c.getContextInterface(cronjob.Context)
	if err != nil {
		panic(err.Error())
	}

	_, err = cl.BatchV1beta1().CronJobs(cronjob.Namespace).Update(context.TODO(), &cronjob.CronJob, v1.UpdateOptions{})
	return true, err
}
