---
title: Sync with Docker
weight: 90
---

`hostctl` allows you synchronize your hosts file with `docker` containers a simple command. 


{{<warning>}}
All `sync` actions will replace previous content of the given profile. 
{{</warning>}}


## Example

Once you have some containers up and running, you can create names for all of them to their container IP.

Command:

`hostctl sync docker -p awesome`

Output:
```
+---------+--------+------------+------------------------+
| PROFILE | STATUS |     IP     |         DOMAIN         |
+---------+--------+------------+------------------------+
| awesome | on     | 172.17.0.2 | mystifying_wescoff.loc |
| awesome | on     | 172.17.0.2 | my-awesome-web.loc     |
+---------+--------+------------+------------------------+
```

{{<info>}}
If you start your docker containers without `--name` it will have a random container name.
{{</info>}}

### Available Options

* `--domain,-d some.domain` domain name used for all containers.

* `--network networkID|networkName` filter only containers of an specific docker network.

* `--profile,-p name` profile where docker information will be set.

