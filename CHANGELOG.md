<a name="0.9.6"></a>
# [0.9.6](https://github.com/moooofly/harbor-go-client/compare/v0.9.5...v0.9.6) (2018-10-17)


### Features

* **api:** add labels APIs ([603570c](https://github.com/moooofly/harbor-go-client/commit/603570c))



<a name="0.9.5"></a>
# [0.9.5](https://github.com/moooofly/harbor-go-client/compare/v0.9.4...v0.9.5) (2018-09-29)

* do some trivial ajustments


<a name="0.9.4"></a>
# [0.9.4](https://github.com/moooofly/harbor-go-client/compare/v0.9.3...v0.9.4) (2018-09-16)


### Features

* add all policy APIs ([53e2193](https://github.com/moooofly/harbor-go-client/commit/53e2193))


<a name="0.9.3"></a>
# [0.9.3](https://github.com/moooofly/harbor-go-client/compare/v0.9.2...v0.9.3) (2018-08-14)


### Features

* **project api:** add prj_member_update/prj_member_get/prj_member_del/prj_member_create/prj_members ([fea77d2](https://github.com/moooofly/harbor-go-client/commit/fea77d2))
* **project api:** add prj_metadata_update_by_name/prj_metadata_get_by_name/prj_metadata_del_by_name ([0020cec](https://github.com/moooofly/harbor-go-client/commit/0020cec))
* **project api:** add prj_metadata_add/prj_metadata_get/prj_logs_get/prj_update APIs ([5770328](https://github.com/moooofly/harbor-go-client/commit/5770328))


<a name="0.9.2"></a>
# [0.9.2](https://github.com/moooofly/harbor-go-client/compare/v0.9.1...v0.9.2) (2018-08-12)


### Features

* **repository api:** add repo_signature_get/repo_image_manifests_get APIs ([95d322d](https://github.com/moooofly/harbor-go-client/commit/95d322d))
* **repository api:** add repo_image_label_del api ([b5b152a](https://github.com/moooofly/harbor-go-client/commit/b5b152a))
* **repository api:** add repo_image_label_add api ([01f7885](https://github.com/moooofly/harbor-go-client/commit/01f7885))
* **repository api:** add repo_image_labels_get/repo_label_del APIs ([54e02be](https://github.com/moooofly/harbor-go-client/commit/54e02be))
* **repository api:** add repo_label_add/repo_labels_get/repo_desp_update APIs ([c12b012](https://github.com/moooofly/harbor-go-client/commit/c12b012))


<a name="0.9.1"></a>
# [0.9.1](https://github.com/moooofly/harbor-go-client/compare/v0.9.0...v0.9.1) (2018-08-06)


### Bug Fixes

* **regression test script:** fix wrong use of printf ([c3ffb03](https://github.com/moooofly/harbor-go-client/commit/c3ffb03))


### Features

* add fake *_test.go ([0f0c6f4](https://github.com/moooofly/harbor-go-client/commit/0f0c6f4))
* **api:** add api/doc.go for golang docs ([928c691](https://github.com/moooofly/harbor-go-client/commit/928c691))


----------


# Changelog (legacy fasion)

## v0.9.0 (2018-7-23)

### New Features

* Add all users APIs

### Improvements

### Bug Fixes

## v0.8.0 (2018-7-11)

### New Features

### Improvements

### Bug Fixes

* Fix compatibility issue derived from harbor v1.3.0. close #2

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
