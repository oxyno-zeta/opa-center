<h1 align="center">opa-center</h1>

<h4 align="center"> This project is in alpha stage ! Things can change. </h4>

<p align="center">
  <a href="https://circleci.com/gh/oxyno-zeta/opa-center" rel="noopener noreferer" target="_blank"><img src="https://circleci.com/gh/oxyno-zeta/opa-center.svg?style=svg" alt="CircleCI" /></a>
  <a href="https://goreportcard.com/report/github.com/oxyno-zeta/opa-center" rel="noopener noreferer" target="_blank"><img src="https://goreportcard.com/badge/github.com/oxyno-zeta/opa-center" alt="Go Report Card" /></a>
  <a href="https://coveralls.io/github/oxyno-zeta/opa-center?branch=master" rel="noopener noreferer" target="_blank"><img src="https://coveralls.io/repos/github/oxyno-zeta/opa-center/badge.svg?branch=master" alt="Coverage Status" /></a>
</p>
<p align="center">
  <a href="https://hub.docker.com/r/oxynozeta/opa-center" rel="noopener noreferer" target="_blank"><img src="https://img.shields.io/docker/pulls/oxynozeta/opa-center.svg" alt="Docker Pulls" /></a>
  <a href="https://github.com/oxyno-zeta/opa-center/blob/master/LICENSE" rel="noopener noreferer" target="_blank"><img src="https://img.shields.io/github/license/oxyno-zeta/opa-center" alt="GitHub license" /></a>
  <a href="https://github.com/oxyno-zeta/opa-center/releases" rel="noopener noreferer" target="_blank"><img src="https://img.shields.io/github/v/release/oxyno-zeta/opa-center" alt="GitHub release (latest by date)" /></a>
</p>

## Description

This project is based on OPA ([Open Policy Agent](https://www.openpolicyagent.org/)) and offer a service to display and save decision logs and status data from OPA server instances. It won't replace the great [Styra](https://www.styra.com/) project. This is done in order to offer a simple alternative to Styra.

## Features

- Save Decision logs and Status data
- Display Decision logs and Status data
- Data retention
- OIDC authentication
- OPA authorization

## Configuration

See [configuration](./docs/configuration.md) documentation.

In [Authorizations](authorizations.md), more information are present about authorization format that can be validated from an OPA server.

## How to deploy ?

## Want to contribute ?

- Read the [CONTRIBUTING](./CONTRIBUTING.md) guide
- Read the [DEVELOPER](./DEVELOPER.md) guide

## Thanks

- My wife BH to support me doing this
- tsandall to accept this project

## Author

- Oxyno-zeta (Havrileck Alexandre)

## License

Apache 2.0 (See in LICENSE)
