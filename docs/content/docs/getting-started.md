---
title: Getting started
weight: 20
asciinema: true
---

## Install 

Read the [Installation guide](installation.md) to get `hostctl` on your system.


## Create a profile

Add a new profile from args:

`hostctl add domains -p awesome my-awesome-ui.project.loc my-awesome-api.project.loc`

``` 
// Output:
+---------+--------+-----------+----------------------------+
| PROFILE | STATUS |    IP     |           DOMAIN           |
+---------+--------+-----------+----------------------------+
| awesome | on     | 127.0.0.1 | my-awesome-ui.project.loc  |
| awesome | on     | 127.0.0.1 | my-awesome-api.project.loc |
+---------+--------+-----------+----------------------------+
```

## Enable or Disable a profile

When you don't want to use some profile, just disable it:

`hostctl disable -p awesome`

``` 
// Output:
+---------+--------+-----------+----------------------------+
| PROFILE | STATUS |    IP     |           DOMAIN           |
+---------+--------+-----------+----------------------------+
| awesome | off    | 127.0.0.1 | my-awesome-ui.project.loc  |
| awesome | off    | 127.0.0.1 | my-awesome-api.project.loc |
+---------+--------+-----------+----------------------------+
```

You can enable it later with: 

`hostctl enable -p awesome`


That's it!


## Linux/Mac/Windows and permissions

The tool recognize your system and use the right hosts file, it will use `/etc/hosts` on Linux/Mac
 and `C:/Windows/System32/Drivers/etc/hosts` on Windows.

**SUDO/ADMIN**: You will need permissions for any action that modify hosts file, add `sudo` to the commands below when needed. 
If you are on windows, make sure you run it as administrator.


