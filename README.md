# ctl
wrapper tool of kubectl for multi-clusters and easy pod/log discovering

# prerequisites
[Setting up k8s envs](https://github.com/wish/k8s/wiki/Setting-up-your-environment-for-k8s), (kubectl + kubeconfig)

## build from source
make build/ctl.<darwin|linux>
and move bin/<darwin|linux>/ctl to the $PATH

# usage

### get
List pods with/without namespace given and over contexts.

| CONTEXT | NAMESPACE | NAME | READY | STATUS | RESTARTS | AGE |
|------|---------|-------|-------|-------|-------|-------
| cluster name | namespace of pod | name of pod | ready containers/total containers | current status |number of restarts | time since starting

Flags:
- `--context, -c` specify the contexts. Can reduce the run time.
- `--namespace, -n` specify the namespace.

### describe pods
Get a detailed description of the pods matching the name query and context/namespace.

Flag:
- `--context, -c` specify the contexts. Can reduce the run time.
- `--namespace, -n` specify the namespace.

### logs pod [flags]
Get the logs given a container in a pod specified. If the pod has only one container, the container name is optional. If the pod has multiple containers, choose one from them. If there are multiple pods that match, the command only operates on the first one found.

Flags:
- `--context, -c` specify the contexts. Can reduce the run time.
- `--namespace, -n` specify the namespace. This could largely reduce the run time of the command.
- `--container, -c` specify the container name. If no container specified, the command will fail.

### sh pod [flags]
Exec /bin/sh into the container of a specific pod. If the pod has only one container, the container name is optional. If the pod has multiple containers, choose one from them. If there are multiple pods that match, the command only operates on the first one found.


Flag:
- `--context, -c` specify the contexts. Can reduce the run time.
- `--namespace, -n` specify the namespace. This could largely reduce the run time of the command.
- `--container, -c` specify the container name.
- `--shell, -s` specify the shell path (default "/bin/sh")

## kron
Kron is a subcommand of ctl for operations involving cron jobs. Generally, the commands follow the same format but do vary due to different requirements.

### get

Like in ctl, `get` retrieves a list of all cron jobs with namespace and context flags.

| NAME | SCHEDULE | SUSPEND | ACTIVE | LAST SCHEDULE | NEXT RUN | AGE | CONTEXT |
|------|------|------|------|------|------|------|------|
| name of cron job | cron schedule | suspended boolean | active pods | time since last run | time until next run | time since adding | cluster name |
Flags:
- `--by-last-run, -l` sorts with latest run cron jobs first
- `--by-last-run-reverse, -L` sorts with oldest run cron jobs first
- `--by-next-run, -e` sorts with the closest next run cron job first
- `--by-next-run-reverse, -E` sorts with farthest next run first

### describe
Print information about the specified cron jobs like in ctl.

Currently displays the following fields: context, name, namespace, schedule, active, last schedule, next run, creation timestamp.

### favorite/unfavorite
A utility function to create a list of cron jobs for easy use with describe and get.
Able to set context and namespace for filtering.

### select
Another utility function to operate on a single function.

These utility commands persist the state throughout runs in `~/.kron/config.yaml`

---
Other commands still need to be done and are a WIP.
