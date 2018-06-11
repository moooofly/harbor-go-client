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

## Features

Current support following Harbor REST APIs:

- Login
    - POST /login
- Logout
    - GET /logout
- Search
    - GET /api/search
- Projects
    - GET /api/projects
    - DELETE /api/projects/{prject_id}
- Repositories
    - DELETE /api/repositories/{repo_name}
    - GET /api/repositories
    - GET /api/repositories/top
    - DELETE /api/repositories/{repo_name}/tags/{tag}
    - GET /api/repositories/{repo_name}/tags/{tag}
    - GET /api/repositories/{repo_name}/tags
- Uses
    - GET /api/users/current
- Targets
    - POST /api/targets
    - DELETE /api/targets/{id}
    - GET /api/targets/{id}
    - PUT /api/targets/{id}
    - GET /api/targets
    - POST /api/targets/ping
    - POST /api/targets/{id}/ping
    - GET /api/targets/{id}/policies/
- Configurations
    - PUT /api/configurations
    - GET /api/configurations
    - POST /api/configurations/reset
- Logs
    - GET /api/logs
- Statistics
    - GET /api/statistics
- Systeminfo
    - GET /api/systeminfo
    - GET /api/systeminfo/volumes
    - GET /api/systeminfo/getcert


Current support following custom feature:

- rp_analysis: Run retention policy analysis, do some soft deletion as you command, prompt user performing a GC.

## Installation

```
go get -u github.com/moooofly/harbor-go-client
```

## Quick Start

- compiles and installs the packages

```
make
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

## Testing

## TODO

- go test
- CI (gitlab-ci & travis-ci）
- dockerization
- godoc
- test.sh colorful output
- glide support
- sql-like result output


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
