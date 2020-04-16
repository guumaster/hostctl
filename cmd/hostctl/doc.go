// contains docs about the command
/*
		    __                    __           __     __
		   / /_   ____    _____  / /_  _____  / /_   / /
		  / __ \ / __ \  / ___/ / __/ / ___/ / __/  / /
		 / / / // /_/ / (__  ) / /_  / /__  / /_   / /
		/_/ /_/ \____/ /____/  \__/  \___/  \__/  /_/


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
  replace     Replace content to a profile in your hosts file.
  restore     Restore hosts file content from a backup file.
  status      Shows a list of profile names and statuses on your hosts file.
  sync        Sync some system IPs with a profile.
  toggle      Change status of a profile on your hosts file.

Flags:
  -c, --column strings     Columns to show on lists
  -h, --help               help for hostctl
      --host-file string   Hosts file path (default "/etc/hosts")
  -q, --quiet              Run command without output
      --raw                Output without table borders
  -v, --version            version for hostctl

Use "hostctl [command] --help" for more information about a command.

*/
package main
