package client

// SetCronJobSuspend sets the cron job suspend state to given value and reports whether successful.
// Does nothing (and returns false) if the cronjob is already in the correct state.
func (c *Client) SetCronJobSuspend(contexts []string, namespace, name string, suspend bool, options ListOptions) (bool, error) {
	cronjob, err := c.findCronJob(contexts, namespace, name, options)
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

	_, err = cl.BatchV1beta1().CronJobs(cronjob.Namespace).Update(&cronjob.CronJob)
	return true, err
}
