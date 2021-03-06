[![Go Report Card](https://goreportcard.com/badge/gojp/goreportcard)](https://goreportcard.com/report/github.com/moooofly/harbor-go-client) [![Build Status](https://travis-ci.org/moooofly/harbor-go-client.svg?branch=master)](https://travis-ci.org/moooofly/harbor-go-client) [![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/moooofly/harbor-go-client/blob/master/LICENSE)

# harbor-go-client

```
 __ __   ____  ____   ____    ___   ____          ____   ___            __  _      ____    ___  ____   ______
 |  |  | /    ||    \ |    \  /   \ |    \        /    | /   \          /  ]| |    |    |  /  _]|    \ |      |
 |  |  ||  o  ||  D  )|  o  )|     ||  D  )_____ |   __||     | _____  /  / | |     |  |  /  [_ |  _  ||      |
 |  _  ||     ||    / |     ||  O  ||    /|     ||  |  ||  O  ||     |/  /  | |___  |  | |    _]|  |  ||_|  |_|
 |  |  ||  _  ||    \ |  O  ||     ||    \|_____||  |_ ||     ||_____/   \_ |     | |  | |   [_ |  |  |  |  |
 |  |  ||  |  ||  .  \|     ||     ||  .  \      |     ||     |      \     ||     | |  | |     ||  |  |  |  |
 |__|__||__|__||__|\_||_____| \___/ |__|\_|      |___,_| \___/        \____||_____||____||_____||__|__|  |__|
```

A CLI tool for the Docker Registry Harbor.

This project offer a command-line interface to the Harbor API, you can use it to manager your users, projects, repositories, etc.

## NOTE

- Currently, the implementation of `harbor-go-client` is based on harbor v1.5.0-d59c257e and swagger api version 1.4.0. If you want to use this CLI tool with Harbro v1.6.0+, you may encounter incompatible issues. See "[The issue with API version](https://github.com/moooofly/harbor-go-client/issues/27)" for more details.
- Another project named [`harborctl`](https://github.com/moooofly/harborctl) is under developement now, which is based on harbor v1.6.0-66709daa and swagger api version 1.6.0.

## Features

Current Harbor API support status:

- Login
    - [x] POST /login
- Logout
    - [x] GET /logout
- search
    - [x] GET /api/search
- projects
    - [x] GET /api/projects
    - [ ] HEAD /api/projects
    - [x] POST /api/projects
    - [x] DELETE /api/projects/{prject_id}
    - [x] GET /api/projects/{prject_id}
    - [x] PUT /api/projects/{prject_id}
    - [x] GET /api/projects/{prject_id}/logs
    - [x] GET /api/projects/{prject_id}/metadatas
    - [x] POST /api/projects/{prject_id}/metadatas
    - [x] DELETE /api/projects/{prject_id}/metadatas/{meta_name}
    - [x] GET /api/projects/{prject_id}/metadatas/{meta_name}
    - [x] PUT /api/projects/{prject_id}/metadatas/{meta_name}
    - [x] GET /api/projects/{prject_id}/members
    - [x] POST /api/projects/{prject_id}/members
    - [x] DELETE /api/projects/{prject_id}/members/{mid}
    - [x] GET /api/projects/{prject_id}/members/{mid}
    - [x] PUT /api/projects/{prject_id}/members/{mid}
- statistics
    - [x] GET /api/statistics
- users
    - [x] GET /api/users
    - [x] POST /api/users
    - [x] GET /api/users/current
    - [x] DELETE /api/users/{user_id}
    - [x] GET /api/users/{user_id}
    - [x] PUT /api/users/{user_id}
    - [x] PUT /api/users/{user_id}/password
    - [x] PUT /api/users/{user_id}/sysadmin
- repositories
    - [x] GET /api/repositories
    - [x] DELETE /api/repositories/{repo_name}
    - [x] PUT /api/repositories/{repo_name}
    - [x] GET /api/repositories/{repo_name}/labels
    - [x] POST /api/repositories/{repo_name}/labels
    - [x] DELETE /api/repositories/{repo_name}/labels/{label_id}
    - [x] DELETE /api/repositories/{repo_name}/tags/{tag}
    - [x] GET /api/repositories/{repo_name}/tags/{tag}
    - [x] GET /api/repositories/{repo_name}/tags
    - [x] GET /api/repositories/{repo_name}/tags/{tag}/labels
    - [x] POST /api/repositories/{repo_name}/tags/{tag}/labels
    - [x] DELETE /api/repositories/{repo_name}/tags/{tag}/labels/{label_id}
    - [x] GET /api/repositories/{repo_name}/tags/{tag}/manifest
    - [ ] POST /api/repositories/{repo_name}/tags/{tag}/scan
    - [ ] GET /api/repositories/{repo_name}/tags/{tag}/vulnerability/details
    - [x] GET /repositories/{repo_name}/signatures
    - [x] GET /api/repositories/top
- logs
    - [x] GET /api/logs
- jobs
    - [x] GET /api/jobs/replication
    - [x] PUT /api/jobs/replication
    - [x] DELETE /api/jobs/replication/{id}
    - [x] GET /api/jobs/replication/{id}/log
    - [x] GET /api/jobs/scan/{id}/log
- policies
    - [x] GET /api/policies/replication
    - [ ] POST /api/policies/replication
    - [x] GET /api/policies/replication/{id}
    - [ ] PUT /api/policies/replication/{id}
- labels
    - [x] GET /api/labels
    - [x] POST /api/labels
    - [x] DELETE /api/labels/{id}
    - [x] GET /api/labels/{id}
    - [x] PUT /api/labels/{id}
- replications
    - [x] POST /api/replications
- targets
    - [x] GET /api/targets
    - [x] POST /api/targets
    - [x] POST /api/targets/ping
    - [x] POST /api/targets/{id}/ping (deprecated)
    - [x] DELETE /api/targets/{id}
    - [x] GET /api/targets/{id}
    - [x] PUT /api/targets/{id}
    - [x] GET /api/targets/{id}/policies/
- internal
    - [x] POST /api/internal/syncregistry
- systeminfo
    - [x] GET /api/systeminfo
    - [x] GET /api/systeminfo/volumes
    - [x] GET /api/systeminfo/getcert
- ldap
    - [ ] POST /api/ldap/ping
    - [ ] GET /api/ldap/groups/search
    - [ ] GET /api/ldap/users/search
    - [ ] POST /api/ldap/users/import
- usergroups
    - [x] GET /api/usergroups
    - [x] POST /api/usergroups
    - [x] DELETE /api/usergroups/{group_id}
    - [x] GET /api/usergroups/{group_id}
    - [x] PUT /api/usergroups/{group_id}
- configurations
    - [x] GET /api/configurations
    - [x] PUT /api/configurations
    - [x] POST /api/configurations/reset
- email
    - [x] POST /api/email/ping


Additional features supported:

- rp_tags: Do tags deletion on repositories according to retention policy.
- rp_repos: Do soft deletion on repositories according to retention policy (prompt user performing a GC after that).

## Installation

```
go get -u github.com/moooofly/harbor-go-client
```
## Quick Start

- lint + build + test

```
make
```

- install the package into $GOPATH

```
make install
```

- remove the corresponding installed archive or binary (what 'go install' would create)

```
make clean
```

- do some tests

```
make test
```

## Documentation

See [docs](https://github.com/moooofly/harbor-go-client/tree/master/docs)

## Testing

You can run integration test with [scripts/regression_test.sh](https://github.com/moooofly/harbor-go-client/blob/master/scripts/regression_test.sh) (Assuming local Harbor installation)

## Auxiliaries Coverage

- [ ] go test
- [x] integration test (by `scripts/*.sh`)
- [x] CI (by travis-ci）
- [x] dockerization
- [x] godoc (need to optimize)
- [x] go module support
- [x] sql-like result output


## Credits

- [parnurzeal/gorequest](https://github.com/parnurzeal/gorequest) - Simplified HTTP client ( inspired by famous SuperAgent lib in Node.js )
- [jessevdk/go-flags](https://github.com/jessevdk/go-flags) - a go library for parsing command line arguments.
- [go-yaml/yaml](https://github.com/go-yaml/yaml) - YAML support for the Go language.

## License

harbor-go-client is licensed under the MIT License. See [LICENSE](https://github.com/moooofly/harbor-go-client/blob/master/LICENSE) for the full license text.

This project uses open source components which have additional licensing terms. The licensing terms for these open source components can be found at the following locations:

- GoRequest: [license](https://github.com/parnurzeal/gorequest/blob/develop/LICENSE)
- go-flags: [license](https://github.com/jessevdk/go-flags/blob/master/LICENSE)
- YAML: [license](https://github.com/go-yaml/yaml/blob/v2/LICENSE)
