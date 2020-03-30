[![Tests][tests-badge]][tests-link]
[![GitHub Release][release-badge]][release-link]
[![Go Report Card][report-badge]][report-link]
[![License][license-badge]][license-link]
<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[![All Contributors](https://img.shields.io/badge/all_contributors-4-orange.svg?style=flat-square)](#contributors-)
<!-- ALL-CONTRIBUTORS-BADGE:END -->


# hostctl

> Manage your /etc/hosts like a pro!

This tool gives you more control over the use of your `hosts` file. You can have multiple profiles and enable/disable as you need.


## Why?

It is a tedious task to handle the `hosts` file by editing manually. With this tool you can automate some aspects to do it cleaner and quick. 

<details>
<summary>Table of content (click to open)</summary>

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->


- [Sample Usage](#sample-usage)
- [Installation](#installation)
  - [Pre-built binary](#pre-built-binary)
  - [Arch Linux](#arch-linux)
  - [HomeBrew](#homebrew)
  - [Snap](#snap)
  - [Scoop](#scoop)
- [Features](#features)
- [Linux/Mac/Windows and permissions](#linuxmacwindows-and-permissions)
- [Usage](#usage)
  - [List profiles](#list-profiles)
  - [Add new profile from a file](#add-new-profile-from-a-file)
  - [Add new profile from cli](#add-new-profile-from-cli)
  - [Enable/Disable profile](#enabledisable-profile)
  - [Remove a profile](#remove-a-profile)
  - [Backup hosts file](#backup-hosts-file)
  - [Restore a hosts file](#restore-a-hosts-file)
  - [TODO](#todo)
  - [References](#references)
  - [LICENSE](#license)
  - [Author(s)](#authors)
- [Contributors ✨](#contributors-)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->
</details>


## Sample Usage
![sample usage](docs/hostctl.gif)


## Installation

### Pre-built binary

Go to [release page](https://github.com/guumaster/hostctl/releases) and download the binary you need.


### Arch Linux

`hostctl` has an AUR package: <https://aur.archlinux.org/packages/hostctl/>. 
You can install it using your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.

Example:
```bash
yay -Sy hostctl
```


### HomeBrew

`brew install guumaster/tap/hostctl`


### Snap

_ Doc to be added after being published and tested _. See [Issue #14](https://github.com/guumaster/hostctl/issues/14) for helping with this.


### Scoop

```
scoop bucket add hostctl https://github.com/guumaster/hostctl.git
scoop install hostctl
```

*NOTE*: If you also installed `sudo` with Scoop, you can run the examples below with `sudo` instead of starting your terminal as administrator.



## Features

  * Manage groups of host names by profile.
  * Enable/disable complete profiles.
  * add/remove groups of host names.
  * add profiles directly from a `.etchosts` file that you can add to your vcs.
  

## Linux/Mac/Windows and permissions

The tool recognize your system and use the right hosts file, it will use `/etc/hosts` on Linux/Mac and `C:/Windows/System32/Drivers/etc/hosts` on Windows.

**SUDO/ADMIN**: You will need permissions for any action that modify hosts file, add `sudo` to the commands below when needed. If you are on windows, make sure you run it as administrator.

**WARNING**: it should work on any system. It's tested on Ubuntu and Windows 10. If you can confirm it works on other system, please let me know [here](https://github.com/guumaster/hostctl/issues/new).



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

### Add new profile from a file

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

### Add new profile from cli

You can add a new profile or add new domain to a specific profile directly from the cli:

You can add that content as a profile with this command:

`$>hostctl -p test add domains test.com --ip 123.123.123.123 `

```
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


## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/gkze"><img src="https://avatars0.githubusercontent.com/u/3131232?v=4" width="50px;" alt=""/><br /><sub><b>George Kontridze</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=gkze" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/pacodes"><img src="https://avatars2.githubusercontent.com/u/28688410?v=4" width="50px;" alt=""/><br /><sub><b>Pacodes</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=pacodes" title="Tests">⚠️</a> <a href="https://github.com/guumaster/hostctl/commits?author=pacodes" title="Code">💻</a></td>
    <td align="center"><a href="https://772424.com"><img src="https://avatars3.githubusercontent.com/u/64371?v=4" width="50px;" alt=""/><br /><sub><b>BarbUk</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=BarbUk" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/devopsbrett"><img src="https://avatars1.githubusercontent.com/u/4403441?v=4" width="50px;" alt=""/><br /><sub><b>Brett Mack</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=devopsbrett" title="Code">💻</a></td>
  </tr>
</table>

<!-- markdownlint-enable -->
<!-- prettier-ignore-end -->
<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!


<!-- JUST BADGES & LINKS -->
[tests-badge]: https://img.shields.io/github/workflow/status/guumaster/hostctl/Test
[tests-link]: https://github.com/guumaster/hostctl/actions?query=workflow%3ATest

[release-badge]: https://img.shields.io/github/release/guumaster/hostctl.svg?logo=github&labelColor=262b30
[release-link]: https://github.com/guumaster/hostctl/releases

[report-badge]: https://goreportcard.com/badge/github.com/guumaster/hostctl
[report-link]: https://goreportcard.com/report/github.com/guumaster/hostctl

[license-badge]: https://img.shields.io/github/license/guumaster/hostctl
[license-link]: https://github.com/guumaster/hostctl/LICENSE

