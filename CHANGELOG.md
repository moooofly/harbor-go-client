# Changelog

## v0.7.0 (2018-7-6)

### New Features

* Add Dockerfile and add 'make docker' cmd

### Improvements

* Optimize *.sh according to shellcheck
* Optimize 'make test' cmd

### Bug Fixes


## v0.6.0 (2018-6-27)

### New Features

* Add jobs related APIs
* Make version info generated automatically

### Improvements

* Integrate packing function into Makefile
* Optimize api/logs.go on values of 'operation'
* Moving scripts/*.yaml to conf/*.yaml

### Bug Fixes


## v0.5.0 (2018-6-22)

### New Features

* Add issue template
* Add `CHANGELOG.md`
* Add `version` cmd showing info

### Improvements

* Update `scripts/binPack.sh` for easy packaging
* Moving *.sh into scripts/
* Rename `rp_repos.sh` -> `rp_repos_simulation.sh`
* Rename `test.sh` -> `regression_test.sh`

### Bug Fixes


## v0.4.0 (2018-6-20)

### New Features

### Improvements

* Add doc
* Optimize log output (more sql-like)

### Bug Fixes


## v0.3.0 (2018-6-19)

### New Features

* Add `binPack.sh` for easy packing

### Improvements

* Optimize log output
* Change specific ip address to `localhost`

### Bug Fixes

* Fix wrongly used heapsort in rp_tags (from maxheap to minheap)


## v0.2.0 (2018-6-13)

### New Features

* Add 'MIT LICENSE' badge
* Add 'Build Status' badge
* Support `misspell` tool
* Support `golint` and `gometalinter` in `.travis.yml`

### Improvements

* Update README.md
* Update LICENSE
* Add some improvements for Travis CI

### Bug Fixes


## v0.1.0 (2018-6-12)

### New Features

* Support `.travis.yml`
* Support **Go Report Card**

### Improvements

* Add `.gitignore`
* Update dependencies managed by `glide`
* Remove `harbor-go-client` executable
* Change the default value of `dstip` to `localhost` in config.yaml

### Bug Fixes

* Fix typo
