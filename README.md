[![Tests][tests-badge]][tests-link]
[![GitHub Release][release-badge]][release-link]
[![Go Report Card][report-badge]][report-link]
[![License][license-badge]][license-link]
[![All Contributors][all-contributors-badge]][all-contributors-link]
[![Coverage][coverage-badge]][coverage-link]
[![Mentioned in Awesome Go][awesome-go-badge]][awesome-go-link]

# hostctl

[![Your dev tool to manage /etc/hosts like a pro!](docs/hostctl.png)](http://guumaster.github.io/hostctl/)

This tool gives you more control over the use of your hosts file. 
You can have multiple profiles and switch them on/off as you need.


## Why?

It is a tedious task to handle the `hosts` file by editing manually.
With this tool you can automate some aspects to do it cleaner and quick. 


## Features

  * Manage groups of host names by profile.
  * Enable/disable complete profiles.
  * Add/remove groups of host names.
  * Add profiles directly from a text file that you can add to your git repo or any VCS.
  * Sync a profile with Docker or Docker Compose containers.
  

## Documentation

You can read about how to use this tool here:

* [Getting started guide](http://guumaster.github.io/hostctl/docs/getting-started/) 
* [Usage Guides](http://guumaster.github.io/hostctl/docs/guides/)
* [CLI commands, options and flags](http://guumaster.github.io/hostctl/docs/cli-usage/)

Or you can start directly on the [Docs homepage](http://guumaster.github.io/hostctl/).


## Sample Usage
![sample usage](docs/hostctl.gif)


## Installation

Go to [release page](https://github.com/guumaster/hostctl/releases) and download the binary you prefer.

Or read about other installation methods on different systems [here](http://guumaster.github.io/hostctl/docs/installation/)


## Linux/Mac/Windows and permissions

The tool recognize your system and use the right hosts file, it will use `/etc/hosts` on Linux/Mac and `C:/Windows/System32/Drivers/etc/hosts` on Windows.

**SUDO/ADMIN**: You will need permissions for any action that modify hosts file, add `sudo` to the commands below when needed. If you are on windows, make sure you run it as administrator.

**WARNING**: it should work on any system. It's tested on Ubuntu and Windows 10. If you can confirm it works on other system, please let me know [here](https://github.com/guumaster/hostctl/issues/new).

## Contributing

Be sure to read [CONTRIBUTING](CONTRIBUTING.md) and [CODE_OF_CONDUCT](CODE_OF_CONDUCT.md).


## Contributors ✨

Thanks goes to these wonderful people ([emoji key](https://allcontributors.org/docs/en/emoji-key)):

<!-- ALL-CONTRIBUTORS-LIST:START - Do not remove or modify this section -->
<!-- prettier-ignore-start -->
<!-- markdownlint-disable -->
<table>
  <tr>
    <td align="center"><a href="https://github.com/gkze"><img src="https://avatars0.githubusercontent.com/u/3131232?v=4?s=50" width="50px;" alt=""/><br /><sub><b>George Kontridze</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=gkze" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/pacodes"><img src="https://avatars2.githubusercontent.com/u/28688410?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Pacodes</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=pacodes" title="Tests">⚠️</a> <a href="https://github.com/guumaster/hostctl/commits?author=pacodes" title="Code">💻</a></td>
    <td align="center"><a href="https://772424.com"><img src="https://avatars3.githubusercontent.com/u/64371?v=4?s=50" width="50px;" alt=""/><br /><sub><b>BarbUk</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=BarbUk" title="Code">💻</a> <a href="#ideas-BarbUk" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/devopsbrett"><img src="https://avatars1.githubusercontent.com/u/4403441?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Brett Mack</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=devopsbrett" title="Code">💻</a></td>
    <td align="center"><a href="https://peterthaleikis.com"><img src="https://avatars0.githubusercontent.com/u/8433587?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Peter Thaleikis</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=spekulatius" title="Code">💻</a> <a href="https://github.com/guumaster/hostctl/commits?author=spekulatius" title="Documentation">📖</a></td>
    <td align="center"><a href="https://github.com/smhmd"><img src="https://avatars0.githubusercontent.com/u/46059092?v=4?s=50" width="50px;" alt=""/><br /><sub><b>simohamed</b></sub></a><br /><a href="#ideas-smhmd" title="Ideas, Planning, & Feedback">🤔</a> <a href="https://github.com/guumaster/hostctl/commits?author=smhmd" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/fabiang"><img src="https://avatars2.githubusercontent.com/u/348344?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Fabian Grutschus</b></sub></a><br /><a href="#ideas-fabiang" title="Ideas, Planning, & Feedback">🤔</a></td>
    <td align="center"><a href="https://github.com/kadern0"><img src="https://avatars.githubusercontent.com/u/13137996?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Pablo Caderno</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/commits?author=kadern0" title="Code">💻</a></td>
    <td align="center"><a href="https://github.com/ddl-jbrown"><img src="https://avatars.githubusercontent.com/u/67287351?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Josh Brown</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/issues?q=author%3Addl-jbrown" title="Bug reports">🐛</a></td>
    <td align="center"><a href="https://github.com/lukasmrtvy"><img src="https://avatars.githubusercontent.com/u/768572?v=4?s=50" width="50px;" alt=""/><br /><sub><b>Lukas M</b></sub></a><br /><a href="https://github.com/guumaster/hostctl/issues?q=author%3Alukasmrtvy" title="Bug reports">🐛</a> <a href="#ideas-lukasmrtvy" title="Ideas, Planning, & Feedback">🤔</a></td>
  </tr>
</table>

<!-- markdownlint-restore -->
<!-- prettier-ignore-end -->

<!-- ALL-CONTRIBUTORS-LIST:END -->

This project follows the [all-contributors](https://github.com/all-contributors/all-contributors) specification. Contributions of any kind welcome!


## Author(s)

* [guumaster](https://github.com/guumaster)

If you like this tool and think is worth it, you can invite me a coffee: 

[![paypal.me/guumaster][paypal-badge]][paypal-link]
[![ko-fi.com/guumaster][kofi-badge]][kofi-link]


## LICENSE

 [MIT license](LICENSE)




<!-- JUST BADGES & LINKS -->
[tests-badge]: https://img.shields.io/github/actions/workflow/status/guumaster/hostctl/tests.yml
[tests-link]: https://github.com/guumaster/hostctl/actions/workflows/tests.yml

[release-badge]: https://img.shields.io/github/release/guumaster/hostctl.svg?logo=github&labelColor=262b30
[release-link]: https://github.com/guumaster/hostctl/releases

[report-badge]: https://goreportcard.com/badge/github.com/guumaster/hostctl
[report-link]: https://goreportcard.com/report/github.com/guumaster/hostctl

[license-badge]: https://img.shields.io/github/license/guumaster/hostctl
[license-link]: https://github.com/guumaster/hostctl/blob/master/LICENSE

[coverage-badge]: https://sonarcloud.io/api/project_badges/measure?project=guumaster_hostctl&metric=coverage
[coverage-link]: https://sonarcloud.io/dashboard?id=guumaster_hostctl

[awesome-go-badge]: https://awesome.re/mentioned-badge.svg
[awesome-go-link]: https://github.com/avelino/awesome-go

<!-- ALL-CONTRIBUTORS-BADGE:START - Do not remove or modify this section -->
[all-contributors-badge]: https://img.shields.io/badge/contributors-10-blue.svg
<!-- ALL-CONTRIBUTORS-BADGE:END -->
[all-contributors-link]: #contributors-

[paypal-link]: https://www.paypal.me/guumaster
[paypal-badge]: https://img.shields.io/static/v1?label=donate&message=PayPal&color=00457C&logo=paypal

[kofi-link]: https://ko-fi.com/guumaster
[kofi-badge]: https://img.shields.io/static/v1?label=donate&message=Ko-fi&color=F16061&logo=ko-fi
