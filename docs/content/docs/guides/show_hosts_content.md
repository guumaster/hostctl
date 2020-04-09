---
title: Show hosts content
weight: 20
---

## List profiles

You can get a list of all profiles using the `list` command. 


### List all profiles

Command: 

`hostctl list`

Output:
```
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


### List only ip and domains, without table decoration

Command:
`list  --raw -c ip,domain`

Output:
```
IP              DOMAIN                                  
127.0.0.1       localhost                               
127.0.1.1       some.existing.local                     
192.168.1.51    jupyter.toolkit-lite.local              
192.168.99.119  app.toolkit.local                       
```

