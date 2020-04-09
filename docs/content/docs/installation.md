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


### HomeBrew

```
brew install guumaster/tap/hostctl
```


### Snap

**WARNING**: Still working out a permission issue * [hostctl snap - store-requests](https://forum.snapcraft.io/t/plugs-system-files-for-hostctl-snap/16199/5) 

```
[sudo] snap install hostctl
```


### Scoop

```
scoop install hostctl
```

*NOTE*: If you also installed `sudo` with Scoop, you can run the examples below with `sudo` instead of starting your terminal as administrator.

