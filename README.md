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


## Usage

`hostctl list`
+---------+--------+----------------+----------------------------+
| PROFILE | STATUS |       IP       |           DOMAIN           |
+---------+--------+----------------+----------------------------+
| lite    | on     | 192.168.1.51   | jupyter.toolkit-lite.local |
+---------+--------+----------------+----------------------------+
| toolkit | on     | 192.168.99.119 | minio.toolkit.local        |
| toolkit | on     | 192.168.99.119 | app.toolkit.local          |
| toolkit | on     | 192.168.99.119 | gitea.toolkit.local        |
| toolkit | on     | 192.168.99.119 | jupyter.toolkit.local      |
| toolkit | on     | 192.168.99.119 | drone.toolkit.local        |
| toolkit | on     | 192.168.99.119 | code.toolkit.local         |
+---------+--------+----------------+----------------------------+


`hostctl from-file`  -f /path/to/.etchosts [-p profile]
`hostctl remove ` [-p profile] 
`hostctl enable`  [-p profile]
`hostctl disable` [-p profile]
`hostctl backup -f path_to_file`
`hostctl restore -f path_to_file`


### TODO

`hostctl from-k8s -n namespace`
`hostctl from-minikube -n namespace`


`hostctl set ip` IP  [-p profile]
`hostctl set domains` dom1,dom2,etc [-p profile]
`hostctl add domains` dom1,dom2 [-p profile]
`hostctl rm domains` dom1,dom2 [-p profile]


#### References

* Dependencies:
  * [spf13/cobra](https://github.com/spf13/cobra)
  * [spf13/viper](https://github.com/spf13/viper)
  * [mitchellh/go-homedir](https://github.com/mitchellh/go-homedir)
  * [guumaster/tablewriter](https://github.com/guumaster/tablewriter)

* Inspiration:
  * [txn2/txeh: CLI for /etc/hosts management](https://github.com/txn2/txeh)


### LICENSE

 [MIT license](LICENSE)


### Author(s)

* [guumaster](https://github.com/guumaster)
