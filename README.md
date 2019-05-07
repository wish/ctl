# wishctl
wrapper tool of kubectl for multi-clusters and easy pod/log discovering

# prerequisites
[Setting up k8s envs](https://github.com/ContextLogic/k8s/wiki/Setting-up-your-environment-for-k8s), (kubectl + kubeconfig)

# install
presumably bin/ folder is up to date. Can just move the execuables to the $PATH

## build from source
make build/wishctl.<darwin|linux>
and move bin/<darwin|linux>/wishctl to the $PATH

# usage

### get
List pods with/without namespace given.

| NAME | READY | STATUS | RESTARTS | AGE |
|------|---------|-------|-------|-------|
| Name of pod | ready containers/total containers | Curent Status | Times of Restarts | How long been running

Flag:
- --namespace, -n specify the namespace. This could largely reduce the run time of the command.

### describe [pod] 
Get a detailed description of the pods matching the name query

Flag:
- --namespace, -n specify the namespace. This could largely reduce the run time of the command.

### log [pod] [flags]
Get the logs given a container in a pod specified. If the pod has only one container, the container name is
optional. If the pod has multiple containers, choose one from them.

Flags:
- --namespace, -n specify the namespace. This could largely reduce the run time of the command.
- --container, -c specify the container name.
- --follow, -f stream pod logs (stdout).
- --tail, -t lines of recent log file to display.

### sh [pod] [flags]
Exec /bin/bash into the container of a specific pod. If the pod has only one container, the container name is
optional. If the pod has multiple containers, choose one from them.

Flag:
- --namespace, -n specify the namespace. This could largely reduce the run time of the command.
- --container, -c specify the container name.
- --shell, -s specify the shell path ((default "/bin/bash")
