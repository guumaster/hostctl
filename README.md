[![Tests](https://img.shields.io/github/workflow/status/guumaster/hostctl/Test)](https://github.com/guumaster/hostctl/actions?query=workflow%3ATest)
[![GitHub Release](https://img.shields.io/github/release/guumaster/hostctl.svg?logo=github&labelColor=262b30)](https://github.com/guumaster/hostctl/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/guumaster/hostctl)](https://goreportcard.com/report/github.com/guumaster/hostctl)
[![License](https://img.shields.io/github/license/guumaster/hostctl)](https://github.com/guumaster/hostctl/LICENSE)

# hostctl

> Manage your /etc/hosts like a pro!

This tool gives you more control over the use of your `hosts` file. You can have multiple profiles and enable/disable as you need.


## Why?

It is a tedious task to handle the `hosts` file by editing manually. With this tool you can automate some aspects to do it cleaner and quick. 


## Installation

Go to [release page](https://github.com/guumaster/hostctl/releases) and download the binary you need.


## Features

  * Manage groups of host names by profile.
  * Enable/disable complete profiles.
  * add/remove groups of host names.
  * add profiles directly from a `.etchosts` file that you can add to your vcs.
  
  
## Sample Usage
![sample usage](docs/hostctl.gif)


## Linux/Mac/Windows and permissions

**WARNING**: this should work on any system, but currently this is only being tested on Linux. 
If you try it on a different system please let me know [here](https://github.com/guumaster/hostctl/issues/new).

The tool recognize your system and use the right hosts file, it will use `/etc/hosts` on Linux/Mac 
and `C:\Windows\System32\Drivers\etc\hosts` on Windows.

**SUDO**: You will need permissions for any action that modify hosts file, add `sudo` to the commands below when needed.


## Usage

```
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

### List profiles

`$> hostctl list`

```
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

### Add new profile

You can store routing as a separate file and add it to the global hosts file when you need.

Say you have this routing file on any of your projects: 

```
# File stored in /path/to/some/project/.etchosts
127.0.0.1 web.my-awesome-project.local 
127.0.0.1 api.my-awesome-project.local 
```

You can add that content as a profile with this command:

`$>hostctl set -p awesome --from /path/to/some/project/.etchosts `

```
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


### Enable/Disable profile

You can enable/disable any profile, the routing will react to it state. 
Disabling a profile does not remove the content from the hosts file, this way you can re-enable it later on.

`$> hostctl disable -p awesome` 
```
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

`$> hostctl enable -p awesome` 
```
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

### Remove a profile

If you want to completely remove a profile from the hosts file you can run:

`$> hostctl remove -p awesome` 

```
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

### Backup hosts file

You can save a copy of your hosts file with this command:

`hostctl backup --path /tmp/`

It would create a file `/tmp/hosts.20200314` with the content of your hosts file.


### Restore a hosts file

You can restore a previous backup of your hosts file with this command:

`hostctl restore --from /tmp/hosts.20200314`

It would **overwrite** your hosts file with the content of your backup.


### TODO

Features that I'd like to add: 

  * [ ] `hostctl from-k8s -n namespace`
  * [ ] `hostctl from-minikube -n namespace`
  * [ ] `hostctl set ip` IP  [-p profile]
  * [ ] `hostctl set domains` dom1,dom2,etc [-p profile]
  * [ ] `hostctl add domains` dom1,dom2 [-p profile]
  * [ ] `hostctl rm domains` dom1,dom2 [-p profile]


### References

* Dependencies:
  * [spf13/cobra](https://github.com/spf13/cobra)
  * [guumaster/tablewriter](https://github.com/guumaster/tablewriter)

* Inspiration:
  * [txn2/txeh: CLI for /etc/hosts management](https://github.com/txn2/txeh)


### LICENSE
 [MIT license](LICENSE)


### Author(s)
* [guumaster](https://github.com/guumaster)
