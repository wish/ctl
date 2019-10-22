#
# Regular cron jobs for the ctl package
#
0 4	* * *	root	[ -x /usr/bin/ctl_maintenance ] && /usr/bin/ctl_maintenance
