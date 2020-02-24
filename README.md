# ctl

[![Build Status](https://travis-ci.org/wish/ctl.svg?branch=master)](https://travis-ci.org/wish/ctl)
[![Code Coverage](https://codecov.io/gh/wish/ctl/branch/master/graph/badge.svg)](https://codecov.io/gh/wish/ctl/)
[![Go Report Card](https://goreportcard.com/badge/github.com/wish/ctl)](https://goreportcard.com/report/github.com/wish/ctl)

A kubectl-like helper for interacting with multiple-clusters concurrently.
___
# Table of Contents

- [Overview](#overview)
- [Getting started](#getting-started)
- [Usage](#usage)
  - [Commands](#commands)
  - [Types](#types)
  - [Multi-cluster usage](#multi-cluster-usage)
  - [Flags](#flags)
  - [Logs](#logs)
  - [Run](#run)
  - [Config](#config)
  - [Cron](#cron)
- [Setup and configuration](#setup-and-configuration)
  - [Labels](#labels)
  - [Hiding clusters](#hiding-clusters)
  - [Default columns](#default-columns)

# Overview
Ctl is a tool that helps with using multiple Kubernetes clusters simultaneously.

# Getting Started
You should have kubectl set up before using ctl. At the least, you should have a kubeconfig file. For wish, see [setting up k8s envs](https://github.com/contextlogic/k8s/wiki/Setting-up-your-environment-for-k8s).

You can compile the binary by running `make`. Move the binary created in `bin` to your `/bin` folder.

# Usage

Example:

```shell
$ ctl get pods bye --k8s_env=dev
CONTEXT               NAMESPACE  NAME                  READY  STATUS     RESTARTS  AGE      K8S_ENV  AZ
app-05-dev.k8s.local  default    bye-1565193600-trc9d  0/1    Succeeded  0         2h17m4s  dev      us-west-1a
app-05-dev.k8s.local  default    bye-1565197200-hq2z6  0/1    Succeeded  0         1h17m5s  dev      us-west-1a
app-05-dev.k8s.local  default    bye-1565200800-zg44s  0/1    Succeeded  0         17m6s    dev      us-west-1a
```

## Commands
The commands get, describe and log function and look similarly to that in kubectl. Get prints out a table of found resources. Describe prints out details about each found resource.

The main commands follow the format

```shell
$ ctl [command] [TYPE] [NAME] [flags]
```

The two main operations are currently get and describe. Get prints out a table with found resources and describe prints out details about each found resource.

## Types

Currently, ctl supports the following resource types (with shorthand in parentheses):
- pods (po)
- cronjobs
- jobs
- deployments (deploy)
- replicasets (rs)
- configmaps (cm)
- k8s_env 

*k8s_env list all possible k8s_env for the clusters

## Multi-cluster usage
Ctl operates across multiple cluster, but it may not be desired to work on all of them. For the advanced user, you may directly specify which cluster(s) to use via the context flag:

```shell
--context=app-05-dev.k8s.local,test-cluster.k8s.local
```

Howevever, for most users, we have added functionality to make it easier to specify and narrow down clusters.

When configured, new labels can be added to all clusters. For example, at Wish ctl now has three global flags: k8s_env, az and region. You may use any combination of these to narrow down which clusters to query on.

```shell
$ ctl get pods --k8s_env=dev
```

Additional columns are added to get output from labels. These columns are set by the config. See [configuration](#setup-and-configuration) for more details.

## Flags
### Namespace
Most objects can be scoped to a namespace. By default, ctl operates on all namespaces, but you may specify a single namespace with the `-n, --namespace` flag. For more expensive operations, specifying contexts and namespaces can improve the performance.

### Labels
Objects can be filtered by labels. The cluster-level labels above are shorthands for these labels filters (`--k8s_env=dev` is equivalent to `-l k8s_env=dev`). Use flag `-l, --label`. You can set any number of labels together: `-l a=b,b=c -l c!=d` is valid.

There are three types of label filters:

**Equal**  
`a=b`. The value of label `a` must equal `b`.

**Not equal**  
`a!=b`. The value of label `a` cannot be `b`.

**Set in**  
`a in (b,c,d)`. The value of label `a` must be one of `b`, `c` or `d`.

### Get Label Columns
For the get command, you can add columns in the table output that corresponds to the value of such label.

## Logs
Logging syntax is quite different from that of other commands. As logging is only for pods, the syntax is just `ctl logs [NAME]`.

Ctl by default logs from the first container of a pod instead of asking to specify a container. If you want logs from a container that is not the first, then you must use the `-c, --container` flag.

A big feature of ctl is the support for multiple logs at a time. As all of ctl uses regex search, if multiple pods are matched, then logs from all of them are printed.

Additionally, this feature supports the `-f, --follow` option which constantly streams logs.

## Run
The run command is a wrapper over `kubectl apply`. It requires cluster config to set up.

The syntax is `ctl run APPNAME [flags]`. The label flags can be used to narrow down which cluster to run on. If there are multiple clusters in which a run is specified on, ctl randomly picks one.

Ctl prints out the command before running.

## Config
Cluster level configs are cached. To update this cached data, run `ctl config fetch`. Sometimes you may need to delete the ctl config folder. This is normally at `$XDG_CONFIG_DIR/ctl` or `~/.config/ctl`.

## Cron
Most cronjob features can be accessed through the base ctl commands. The `ctl cron` command allows for diect manipulation of k8s cronjobs. You may see `ctl help cron` for more details.

# Setup and Configuration
This section refers to the optional setup on the server-side of configuration ctl for all users. To use the optional features of ctl, a ConfigMap should be added to clusters.

This ConfigMap should be located in namespace `kube-system` and have name `ctl-config`.

The data of the ConfigMap can be used to specify behaviour.

## Labels
All fields in the ConfigMap that are not prefixed with an underscore are set as cluster-level labels. These labels will be propagated down to the objects on each cluster. Setting these labels are useful for describing and filtering the clusters.

## Hiding clusters
A cluster can be hidden from users by setting `_hidden` to be true.

## Default columns
You can set default label columns to be printed with get output via `_default_columns`. Set this to be a comma separated list in string format. E.g. `"_default_columns":"col1,able"`.

