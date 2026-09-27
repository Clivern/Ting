<p align="center">
    <h3 align="center">Ting</h3>
    <p align="center">A reverse proxy for the OpenRouter API</p>
    <p align="center">
        <a href="https://github.com/clivern/ting/actions/workflows/ci.yml">
            <img alt="CI" src="https://github.com/clivern/ting/actions/workflows/ci.yml/badge.svg">
        </a>
        <a href="https://github.com/clivern/ting/releases">
            <img src="https://img.shields.io/badge/Version-v0.1.0-red.svg">
        </a>
        <a href="https://github.com/clivern/ting/blob/main/LICENSE">
            <img src="https://img.shields.io/badge/LICENSE-Apache%202.0-blue.svg">
        </a>
    </p>
</p>

`Ting` sits in front of [OpenRouter](https://openrouter.ai) and injects the API key on the way through.


### Install

Download a release from the [releases page](https://github.com/clivern/ting/releases), or build from source:

```bash
go build -o ting .
```


### Usage

Copy `config.dist.yml` if you want a local override, export the API key, and start the proxy:

```bash
export OPENROUTER_API_KEY=sk-or-...
export ZIEE_MGMT_URL=http://127.0.0.1:9090
ting server -c config.dist.yml
```

Print build information:

```bash
ting version
```


### Versioning

For transparency into our release cycle and in striving to maintain backward compatibility, ting is maintained under the [Semantic Versioning guidelines](https://semver.org/) and release process is predictable and business-friendly.

See the [Releases section of our GitHub project](https://github.com/clivern/ting/releases) for changelogs for each release version of ting. It contains summaries of the most noteworthy changes made in each release. Also see the [Milestones section](https://github.com/clivern/ting/milestones) for the future roadmap.


### Bug tracker

If you have any suggestions, bug reports, or annoyances please report them to our issue tracker at https://github.com/clivern/ting/issues


### Security Issues

If you discover a security vulnerability within ting, please send an email to [hello@clivern.com](mailto:hello@clivern.com)


### Contributing

We are an open source, community-driven project so please feel free to join us. see the [contributing guidelines](CONTRIBUTING.md) for more details.


### License

© 2026 Clivern. Released under the [Apache License 2.0](https://www.apache.org/licenses/LICENSE-2.0).

**Ting** is authored and maintained by [@Clivern](http://github.com/Clivern).
