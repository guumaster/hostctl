---
title: Backup & Restore
weight: 80
---


## Backup your hosts file

You can save a copy of your hosts file with this command:

Command:

`hostctl backup --path $HOME/hostctl/`

It will create a file `$HOME/hostctl/hosts.20200314` with the full content of your hosts file.



## Restore a hosts file

You can restore a previous backup of your hosts file with this command:

Command: 

`hostctl restore --from $HOME/hostctl/hosts.20200314`


{{<danger>}}
This action will **overwrite** your hosts file with the content of your backup. It cannot be undone.
{{</danger>}}


