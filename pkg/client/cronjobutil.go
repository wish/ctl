package client

// Sets the cron job suspend state to given value.
// Does nothing (and returns false) if already set.
func (c *Client) SetCronJobSuspend(contexts []string, namespace, name string, suspend bool) (bool, error) {
	cronjob, err := c.findCronJob(contexts, namespace, name)
	if err != nil {
		return false, err
	}
	if *cronjob.Spec.Suspend == suspend { // Already set to value
		return false, nil
	}
	*cronjob.Spec.Suspend = suspend

	cl, err := c.getContextClientset(cronjob.Context)
	if err != nil {
		panic(err.Error())
	}

	_, err = cl.BatchV1beta1().CronJobs(cronjob.Namespace).Update(&cronjob.CronJob)
	return true, err
}
