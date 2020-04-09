---
title: Sync with Docker Compose
weight: 95
---

`hostctl` allows you synchronize your hosts file with `docker-compose` project a simple command.

{{<warning>}}
All `sync` actions will replace previous content of the given profile.
{{</warning>}}


## Example

Command:

`hostctl sync docker-compose`

Output:
```
+------------------+--------+--------------+--------------+
|     PROFILE      | STATUS |      IP      |    DOMAIN    |
+------------------+--------+--------------+--------------+
| examplevotingapp | on     | 192.168.16.6 | worker_1.loc |
| examplevotingapp | on     | 192.168.16.2 | db.loc       |
| examplevotingapp | on     | 192.168.16.3 | redis.loc    |
| examplevotingapp | on     | 192.168.16.4 | result_1.loc |
| examplevotingapp | on     | 172.31.0.2   | result_1.loc |
| examplevotingapp | on     | 192.168.16.5 | vote_1.loc   |
| examplevotingapp | on     | 172.31.0.3   | vote_1.loc   |
+------------------+--------+--------------+--------------+
```

**NOTE**: This example output when `hostctl` is used with [example-voting-app](https://github.com/dockersamples/example-voting-app).


### Available Options

* `--domain,-d some.domain` domain name used for all containers.

* `--network networkID|networkName` filter only containers of an specific docker network.

* `--profile,-p name` profile where docker information will be set. Default to Docker Compose project name.

* `--compose-file /path/to/docker-compose.yml` set the path of docker-compose.yml file to use. Defaults to `$PWD/docker-compose.yml`

* `--prefix` keep the prefix used by Docker Compose based on the folder name. Defaults to `false`.
