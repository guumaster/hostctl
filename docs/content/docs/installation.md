---
title: Installation
weight: 10
---

### Pre-built binary

Go to the [latest release page](https://github.com/guumaster/hostctl/releases/latest) and download the binary you need.


### Ubuntu

Get the `.deb` package from the [latest release page](https://github.com/guumaster/hostctl/releases/latest) and then run:

```
sudo dpkg -i /path/to/downloaded/hostctl_<version>.deb
```


### Arch Linux

There is an [AUR package](https://aur.archlinux.org/packages/hostctl) for `hostctl`. 

You can install it using your [AUR helper](https://wiki.archlinux.org/index.php/AUR_helpers) of choice.

Example:
```
yay -Sy hostctl
```


### Nix

```
nix-env -iA nixpkgs.hostctl
```


### HomeBrew

```
brew install guumaster/tap/hostctl
```

### asdf

```
asdf plugin add hostctl https://github.com/svenluijten/asdf-hostctl.git
asdf install hostctl latest
```

### Scoop

```
scoop install hostctl
```

*NOTE*: If you also installed `sudo` with Scoop, you can run the examples below with `sudo` instead of starting your terminal as administrator.


### Snap [DEPRECATED]

**DEPRECATION NOTICE**: Last version supported on snap is `v1.0.11`. I think Snap is not for everyone. Certainly not for me. 
I've tried to maintain a snap version but is too manual, cumbersome, random and boring.
Please, get the `hostctl` binary in any other flavor.
