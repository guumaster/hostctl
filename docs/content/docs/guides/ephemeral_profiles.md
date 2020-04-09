---
title: Ephemeral profiles
weight: 80
---

With `hostctl` you can add content to your `hosts` file temporarily with the `--wait` flag. 

The command you run will wait for the given time or until you hit `ctrl-c` signal.
 
After that, the action will be undone automatically.

## Examples

## Add temporal and wait for ctrl-c

Command:
`cat .etchosts | hostctl add -p ephemeral  --wait 0`

Output:
+-----------+--------+------------+------------------------------+
|  PROFILE  | STATUS |     IP     |            DOMAIN            |
+-----------+--------+------------+------------------------------+
| awesome   | on     | 127.0.0.1  | web.my-awesome-project.local |
| awesome   | on     | 127.0.0.1  | api.my-awesome-project.local |
+-----------+--------+------------+------------------------------+

Waiting for 10s or ctrl+c to remove from profile 'ephemeral'

Profile 'ephemeral' removed.

