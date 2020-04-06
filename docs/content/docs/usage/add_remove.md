---
title: Add & Remove
weight: 10
---


## Add 

{{<info>}}
**Add or Set**: For all this examples it's possible to replace `add` with `set` and 
it will replace all previous content of the chosen profile with the new one.
{{</info>}}

### Available Flags


## Examples


### Add new profile from a file
You can store routing as a separate file and add it to the global hosts file when you need.
For example, if you have this routing file on any of your projects: 

```bash
# Sample stored in /path/of/some/project/repo/.etchosts
127.0.0.1 web.my-awesome-project.local 
127.0.0.1 api.my-awesome-project.local 
```

Command:

`hostctl set -p awesome --from /path/to/some/project/.etchosts `

Output:
```bash
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

{{<warning>}}
If you installed from the Snap Store you **can't use** `--from file` due to Snap's confinement restriction. 
See [stdin option](#add-new-profile-from-stdin)
{{</warning>}}


### Add domains from args

You can add a new profile or add new domain to a specific profile directly from the cli:

You can add that content as a profile with this command:

`hostctl add domains test.loc another.loc  --ip 123.123.123.123 -p test`

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
| test    | on     | 123.123.123.123 | test.loc                     |
| test    | on     | 123.123.123.123 | another.loc                  |
+---------+--------+-----------------+------------------------------+
```


### Add new profile from `stdin`

Similar to the previous option, you can pipe from a previous command or redirect output to it.

{{<warning>}}
If you installed from the Snap Store this is the only way to add content from files due to Snap's confinement restrictions.
{{</warning>}}

Commands (both work similarly):

* `cat /path/to/some/project/.etchosts | hostctl set -p awesome`
* `hostctl add -p awesome < /path/to/some/project/.etchosts`

Output:
```bash
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


## Remove a profile

If you want to completely remove a profile from the hosts file you can run:

`hostctl remove -p awesome` 

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

{{<warning>}}
This action cannot be undone. If you need to enable the profile later, use `disable` instead.
{{</warning>}}
