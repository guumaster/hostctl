---
title: Usage
weight: 30
---

### Overview

```bash
 _     _  _____  _______ _______ _______ _______       
 |_____| |     | |______    |    |          |    |     
 |     | |_____| ______|    |    |_____     |    |_____

hostctl is a CLI tool to manage your hosts file with ease. 
You can have multiple profiles, enable/disable exactly what
you need each time with a simple interface.

Usage:
  hostctl [command]

Available Commands:
  add         Add content to a profile in your hosts file.
  backup      Creates a backup copy of your hosts file
  disable     Disable a profile from your hosts file.
  enable      Enable a profile on your hosts file.
  help        Help about any command
  list        Shows a detailed list of profiles on your hosts file.
  remove      Remove a profile from your hosts file.
  restore     Restore hosts file content from a backup file.
  set         Set content to a profile in your hosts file.

Flags:
  -h, --help               help for hostctl
      --host-file string   Hosts file path (default "/etc/hosts")
  -p, --profile string     Choose a profile

Use "hostctl [command] --help" for more information about a command.
```

### Available actions

* [List](#list-profiles)
* [Add from a file](#add-new-profile-from-a-file)
* [Add from args](#add-new-profile-from-args)
* [Remove](#remove-a-profile)
* [Enable](#enable-a-profile)
* [Disable](#disable-a-profile)
* [Backup hosts file](#backup-hosts-file)
* [Restore hosts file](#restore-a-hosts-file)


### List profiles

`$> hostctl list`

```bash
// Output:
+---------+--------+----------------+----------------------------+
| PROFILE | STATUS |       IP       |           DOMAIN           |
+---------+--------+----------------+----------------------------+
| default | on     | 127.0.0.1      | localhost                  |
| default | on     | 127.0.1.1      | some-existing.local        |
| default | on     | ::1            | ip6-localhost              |
+---------+--------+----------------+----------------------------+
| lite    | on     | 192.168.1.51   | jupyter.toolkit-lite.local |
+---------+--------+----------------+----------------------------+
| toolkit | on     | 192.168.99.119 | app.toolkit.local          |
| toolkit | on     | 192.168.99.119 | gitea.toolkit.local        |
| toolkit | on     | 192.168.99.119 | jupyter.toolkit.local      |
+---------+--------+----------------+----------------------------+
```


### Add new profile from a file

You can store routing as a separate file and add it to the global hosts file when you need.

Say you have this routing file on any of your projects: 

```bash
# File stored in /path/to/some/project/.etchosts
127.0.0.1 web.my-awesome-project.local 
127.0.0.1 api.my-awesome-project.local 
```

You can add that content as a profile with this command:

`$> hostctl set -p awesome --from /path/to/some/project/.etchosts `

```bash
// Output:
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
| default | on     | 127.0.0.1      | localhost                    |
| default | on     | 127.0.1.1      | some-existing.local          |
| default | on     | ::1            | ip6-localhost                |
+---------+--------+----------------+------------------------------+
| another | on     | 192.168.1.51   | jupyter.toolkit-lite.local   |
+---------+--------+----------------+------------------------------+
| awesome | on     | 127.0.0.1      | web.my-awesome-project.local |
| awesome | on     | 127.0.0.1      | api.my-awesome-project.local |
+---------+--------+----------------+------------------------------+
```


### Add new profile from args

You can add a new profile or add new domain to a specific profile directly from the cli:

You can add that content as a profile with this command:

`$> hostctl -p test add domains test.com --ip 123.123.123.123 `

```bash
// Output:
+---------+--------+-----------------+------------------------------+
| PROFILE | STATUS |       IP        |            DOMAIN            |
+---------+--------+-----------------+------------------------------+
| default | on     | 127.0.0.1       | localhost                    |
| default | on     | 127.0.1.1       | some-existing.local          |
| default | on     | ::1             | ip6-localhost                |
+---------+--------+-----------------+------------------------------+
| another | on     | 192.168.1.51    | jupyter.toolkit-lite.local   |
+---------+--------+-----------------+------------------------------+
| test    | on     | 123.123.123.123 | test.com                     |
+---------+--------+-----------------+------------------------------+
```


### Remove a profile

If you want to completely remove a profile from the hosts file you can run:

`$> hostctl remove -p awesome` 

```bash
// Output:
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
| default | on     | 127.0.0.1      | localhost                    |
| default | on     | 127.0.1.1      | some-existing.local          |
| default | on     | ::1            | ip6-localhost                |
+---------+--------+----------------+------------------------------+
| another | off    | 192.168.1.51   | jupyter.toolkit-lite.local   |
+---------+--------+----------------+------------------------------+
```


### Enable a profile

You can enable any profile, the routing will react to it state. 

`$> hostctl enable -p awesome` 
```bash
// Output:
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
...
+---------+--------+----------------+------------------------------+
| awesome | on     | 127.0.0.1      | web.my-awesome-project.local |
| awesome | on     | 127.0.0.1      | api.my-awesome-project.local |
+---------+--------+----------------+------------------------------+
```


### Disable a profile

You can disable any profile, all routing for that profile will stop working. 

Disabling a profile does not remove the content from the hosts file, this way you can re-enable it later on.

`$> hostctl disable -p awesome` 

```bash
// Output:
+---------+--------+----------------+------------------------------+
| PROFILE | STATUS |       IP       |            DOMAIN            |
+---------+--------+----------------+------------------------------+
...
+---------+--------+----------------+------------------------------+
| awesome | off    | 127.0.0.1      | web.my-awesome-project.local |
| awesome | off    | 127.0.0.1      | api.my-awesome-project.local |
+---------+--------+----------------+------------------------------+
```


### Backup hosts file

You can save a copy of your hosts file with this command:

`hostctl backup --path /tmp/`

It would create a file `/tmp/hosts.20200314` with the content of your hosts file.


### Restore a hosts file

You can restore a previous backup of your hosts file with this command:

`hostctl restore --from /tmp/hosts.20200314`

It would **overwrite** your hosts file with the content of your backup.

