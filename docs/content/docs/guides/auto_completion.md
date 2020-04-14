---
title: Auto completion
weight: 40
---


### Bash

Add this line on your `$HOME/.bashrc`

```
source <(hostctl completion bash)
```


### Zsh

Add this line on your `$HOME/.zshrc`

```
source <(hostctl completion zsh)
```

**NOTE**: If you are using `oh-my-zsh` this method won't work. Check below.


### Oh-My-Zsh

- First generate the plugin with auto completion code

```
hostctl completion zsh > $HOME/.oh-my-zsh/plugins/hostctl/_hostctl
```

- Add it to your plugin list in $HOME/.zshrc

```
plugins=(... hostctl)
```

- Check that this lines are present somewere in $HOME/.zshrc 

```
autoload -U compinit
compinit
```

