package client

import (
	"fmt"
	"github.com/wish/ctl/pkg/client/types"
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"time"
)

// RunCronJob creates a new job with timestamp from the specified cron job template
func (c *Client) RunCronJob(contexts []string, namespace, cronjobName string, options ListOptions) (*types.RunDiscovery, error) {
	cronjob, err := c.findCronJob(contexts, namespace, cronjobName, options)
	if err != nil {
		return nil, err
	}

	cl, err := c.getContextInterface(cronjob.Context)
	if err != nil {
		panic(err.Error())
	}

	job, err := cl.BatchV1().Jobs(cronjob.Namespace).Create(
		&batchv1.Job{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", cronjobName, time.Now().Unix()), // REVIEW: What if name is not unique??
				Namespace: cronjob.Namespace,
				OwnerReferences: []metav1.OwnerReference{
					metav1.OwnerReference{
						APIVersion: "batch/v1beta1",
						Kind:       "CronJob",
						Name:       cronjob.Name,
						UID:        cronjob.UID,
						// TODO: Set the Controller and BlockOwnerDeletion fields
					},
				},
			},
			Spec: cronjob.CronJob.Spec.JobTemplate.Spec,
		})

	if err != nil {
		return nil, err
	}

	return &types.RunDiscovery{cronjob.Context, *job}, nil
}
