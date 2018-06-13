# WARNING

> **Notice:** This document is **deprecated**, this here is kept for historical purpose.

# Retention Policy 设计

## 前提

### 现状

- 线上环境使用的是 Harbor v0.5.0 版本；
- 基于 Harbor v1.2.0 搭建了测试环境；
- 基于 `swagger.yaml` 和 `http://editor.swagger.io/` ，并参考 [View and test Harbor REST API via Swagger](https://github.com/vmware/harbor/blob/master/docs/configure_swagger.md) 进行了 RESTful API 研究和测试，确认目前 Harbor 提供的 API 中没有满足“清理最近不常用的 repos”需求的接口；
- 目前最接近需求的两个接口为
    - `GET /repositories` - Get repositories accompany with relevant project and repo name.
    - `GET /repositories/top` - Get public repositories which are accessed most.
    - 对应的 Model Schema 为
```
[
  {
    "id": 0,
    "name": "string",
    "project_id": 0,
    "description": "string",
    "pull_count": 0,
    "star_count": 0,
    "tags_count": 0,
    "creation_time": "string",
    "update_time": "string"
  }
]
```

- Harbor 官方也表示尚未支持该功能，参见：
    - [issues/1168](https://github.com/vmware/harbor/issues/1168)
    - [issues/2287](https://github.com/vmware/harbor/issues/2287)
    - [issues/737](https://github.com/vmware/harbor/issues/737)
    - [issues/3342](https://github.com/vmware/harbor/issues/3342)
    - [issues/3274](https://github.com/vmware/harbor/issues/3274)
    - [issues/1726](https://github.com/vmware/harbor/issues/1726)
    - [issues/3456](https://github.com/vmware/harbor/issues/3456)
    - [issues/3568](https://github.com/vmware/harbor/issues/3568)
- 如果要自行实现“清理最近不常用的 repos”功能，则可以考虑基于 `pull_count` + `star_count` + `tags_count` + `creation_time` + `update_time` 的组合，自定义某种规则进行处理；

### 分析

- 理论上讲，在 harbor 中针对 `projects`/`repositories`/`tags`/`images` 这些对象都可以进行清理操作（其中 images 的处理和其他有所不同）；
- 当前设计主要针对 **repositories** 级别进行清理操作；
- harbor-go-client 的实现基于 harbor v1.2.0 版本进行；



## 代码研究

在 harbor 的 `src/common/dao/repository.go` 中有

```
[#237#root@ubuntu-1604 /go/src/github.com/vmware/harbor]$vi ./src/common/dao/repository.go
...
109 //GetTopRepos returns the most popular repositories whose project ID is
110 // in projectIDs
111 func GetTopRepos(projectIDs []int64, n int) ([]*models.RepoRecord, error) {
112     repositories := []*models.RepoRecord{}
113     if len(projectIDs) == 0 {
114         return repositories, nil
115     }
116
117     _, err := GetOrmer().QueryTable(&models.RepoRecord{}).
118         Filter("project_id__in", projectIDs).
119         OrderBy("-pull_count").
120         Limit(n).
121         All(&repositories)
122
123     return repositories, err
124 }
...
```

可以看出其 top 实现为：

- 在指定的 `projectIDs` 集合中，
- 按照 `pull_count` 降序排序后，
- 排名前 `n` 的结果

另外，从源码中还得知，只有在外部触发 `docker pull` 动作时，才会导致 `pull_count` 的 +1 ，以及 `update_time` 的值更新；

```
[#237#root@ubuntu-1604 /go/src/github.com/vmware/harbor]$vi ./src/common/dao/repository.go
...
 74 // IncreasePullCount ...
 75 func IncreasePullCount(name string) (err error) {
 76     o := GetOrmer()
 77     num, err := o.QueryTable("repository").Filter("name", name).Update(
 78         orm.Params{
 79             "pull_count":  orm.ColValue(orm.ColAdd, 1),
 80             "update_time": time.Now(),
 81         })
 82     if err != nil {
 83         return err
 84     }
 85     if num == 0 {
 86         return fmt.Errorf("Failed to increase repository pull count with name: %s", name)
 87     }
 88     return nil
 89 }
 ...
```

## 策略制定

在设计 harbor-go-client 的 Retention Policy 清理策略时，希望能够基于更多维度实现更加智能一点（更复杂）的策略；

在实现中综合考虑了以下**加权项**（目前确认 harbor v1.2.0 支持）：

- `update_time` - 该值代表最后一次更新的时间（即最后一次被 pull 的时间），可准确反映用户的使用情况；
- `pull_count` - 该数值越大，代表被 pull 的次数越多，说明越受欢迎；
- `star_count` - 该数值越大，代表用户对 projects 或 repositories 的认可程度越高，说明其价值越高（官方表示目前尚不支持使用：[issues/3568](https://github.com/vmware/harbor/issues/3568)）；
- `tags_count` - 供参考使用，tags 多只能说曾被频繁使用；
- `creation_time` - 供参考使用，可以考虑最近创建的目标不做处理（当前设计中暂不使用）；

**详细规格**：

| 权重项 | 分值 | 说明 |
| -- | -- | -- |
| update_time (u) | 0.5 | (u_factor, uf) <br> Unit: day <br> (* 1.0) N <= 7  <br> (* 0.9) 7 < N <= 30 <br> (* 0.8) N > 30 |
| pull_count (p) | 0.3 | (p_factor, pf) <br> Unit: count <br> (* 1.0) N >= 50  <br> (* 0.9) 10 <= N < 50 <br> (* 0.8) N < 10 |
| star_count (s) | 0.0 | (s_factor, sf) <br> 暂时不考虑使用 |
| tags_count (t) | 0.1 | (t_factor, tf) <br> Unit: count <br> (* 1.0) N >= 10 <br> (* 0.9) 3 <= N < 10 <br> (* 0.8) N < 3 |
| creation_time (c) | 0.1 | (c_factor, cf) <br> 暂时不考虑使用 |

**计分公式**：

```
score = u * uf + p * pf + t * tf
```

> 注意：
>
> - 当前计分公式比较简单，可以按需调整；
> - 基于当前公式，若继续增加分支档位设置（如 0.7），则会出现分值计算相等的情况（代码也得修改），因此可能需要重新规划分值分布；
> - 基于当前公式，若调整每个档位的分值设置，则可能会出现分值计算相等的情况，同样可能需要重新规划分值分布；


**计分示例**：

> `prj/repo1` 在 3 天前被 pull 过，历史 pull 数量为 1000 ，tags 数量为 20

score = 0.5 * 1.0 (`update_time`) + 0.3 * 1.0 (`pull_count`) + 0.1 * 1.0 (`tags_count`) = 0.9

> `prj/repo2` 在 10 天前被 pull 过，历史 pull 数量为 1000 ，tags 数量为 20

score = 0.5 * 0.9 + 0.3 * 1.0 + 0.1 * 1.0 = 0.85

> `prj/repo3` 在 20 天前被 pull 过，历史 pull 数量为 20 ，tags 数量为 2

score = 0.5 * 0.9 + 0.3 * 0.9 + 0.1 * 0.8 = 0.80

> `prj/repo4` 在 3 天前被 pull 过，历史 pull 数量为 20 ，tags 数量为 20

score = 0.5 * 1.0 + 0.3 * 0.9 + 0.1 * 1.0 = 0.87

> `prj/repo5` 在 20 天前被 pull 过，历史 pull 数量为 1000 ，tags 数量为 2

score = 0.5 * 0.9 + 0.3 * 1.0 + 0.1 * 0.8 = 0.83



## 实际情况

线上 staging 环境中存在 209 个 repo ；

```
deployer@docker-registry-stag:~/fei$ ./harbor-go-client repos_top|grep "\"count\""|wc -l
10
deployer@docker-registry-stag:~/fei$ ./harbor-go-client repos_top -c 100|grep "\"count\""|wc -l
100
deployer@docker-registry-stag:~/fei$ ./harbor-go-client repos_top -c 1000|grep "\"count\""|wc -l
209
deployer@docker-registry-stag:~/fei$
```

实际 pull_count 的数值分布如下（top100）：

```
$ ./harbor-go-client repos_top -c 100
...
[
  {
    "name": "library/mysql",
    "count": 13431
  },
  {
    "name": "library/redis",
    "count": 13232
  },
  {
    "name": "backend/neo_staging",
    "count": 7552
  },
  {
    "name": "backend/testlb",
    "count": 6865
  },
  {
    "name": "backend/darwin",
    "count": 1839
  },
  {
    "name": "engzo-ci/ruby_mysql",
    "count": 1146
  },
  {
    "name": "backend/neo_nginx_staging",
    "count": 679
  },
  {
    "name": "library/golang",
    "count": 474
  },
  {
    "name": "platform/localstack",
    "count": 472
  },
  {
    "name": "backend/llspay",
    "count": 465
  },
  {
    "name": "backend/pronco",
    "count": 443
  },
  {
    "name": "library/k8s-fluentd-kafka",
    "count": 430
  },
  {
    "name": "backend/answerup",
    "count": 398
  },
  {
    "name": "algorithm-rls/essay_scoring_service",
    "count": 336
  },
  {
    "name": "backend/ruby",
    "count": 311
  },
  {
    "name": "backend/kensaku",
    "count": 235
  },
  {
    "name": "backend/gather",
    "count": 217
  },
  {
    "name": "algorithm-rls/dialoguesrv",
    "count": 216
  },
  {
    "name": "engzo/lingome-api-doc",
    "count": 215
  },
  {
    "name": "library/maven",
    "count": 200
  },
  {
    "name": "library/bazel",
    "count": 190
  },
  {
    "name": "backend/ielts_server",
    "count": 188
  },
  {
    "name": "frontend/lms-web",
    "count": 186
  },
  {
    "name": "algorithm-rls/chatsrv",
    "count": 176
  },
  {
    "name": "engzo-ci/neo",
    "count": 172
  },
  {
    "name": "data-ci/sqlambda-ci",
    "count": 152
  },
  {
    "name": "data-ci/python-web",
    "count": 138
  },
  {
    "name": "data-ci/sqlbuffet-base",
    "count": 135
  },
  {
    "name": "micanzhang/darwin-base",
    "count": 135
  },
  {
    "name": "data-ci/presto",
    "count": 134
  },
  {
    "name": "backend/leader_board",
    "count": 133
  },
  {
    "name": "backend/tape",
    "count": 126
  },
  {
    "name": "library/mongo",
    "count": 102
  },
  {
    "name": "backend/link_agent_staging",
    "count": 101
  },
  {
    "name": "algorithm-rls/reportsrv",
    "count": 98
  },
  {
    "name": "backend/entrance_exams",
    "count": 98
  },
  {
    "name": "algorithm-rls/suggsrv",
    "count": 93
  },
  {
    "name": "algorithm-rls/modelsrv",
    "count": 90
  },
  {
    "name": "algorithm-rls/lookupsrv",
    "count": 85
  },
  {
    "name": "backend/llspay-http",
    "count": 83
  },
  {
    "name": "algorithm-rls/the-speaking-force",
    "count": 78
  },
  {
    "name": "frontend/freetalk",
    "count": 71
  },
  {
    "name": "backend/ielts_cms",
    "count": 70
  },
  {
    "name": "algorithm-rls/chatbot_asr",
    "count": 69
  },
  {
    "name": "backend/course-center-service",
    "count": 69
  },
  {
    "name": "backend-rls/study-performance",
    "count": 67
  },
  {
    "name": "engzo-ci/neo-intest",
    "count": 65
  },
  {
    "name": "backend/llspay-outer-http",
    "count": 64
  },
  {
    "name": "backend/ccbase",
    "count": 61
  },
  {
    "name": "backend/fibre",
    "count": 61
  },
  {
    "name": "algorithm-rls/lqpt",
    "count": 60
  },
  {
    "name": "backend/anatawa",
    "count": 59
  },
  {
    "name": "platform/gitlab-runner",
    "count": 59
  },
  {
    "name": "algorithm-rls/anti-spam",
    "count": 52
  },
  {
    "name": "algorithm-rls/debugsrv",
    "count": 46
  },
  {
    "name": "algorithm-rls/regtest",
    "count": 46
  },
  {
    "name": "algorithm-rls/sesame-acceptor",
    "count": 44
  },
  {
    "name": "analysis/hive-mysql-dump-ci",
    "count": 44
  },
  {
    "name": "backend-rls/cc-performance",
    "count": 44
  },
  {
    "name": "backend/cooper",
    "count": 42
  },
  {
    "name": "suhai/golang",
    "count": 42
  },
  {
    "name": "engzo-ci/nexus-uploader",
    "count": 40
  },
  {
    "name": "algorithm-rls/keywords",
    "count": 39
  },
  {
    "name": "library/python",
    "count": 38
  },
  {
    "name": "platform/grafana-k8s",
    "count": 37
  },
  {
    "name": "analysis/airflow-ci",
    "count": 36
  },
  {
    "name": "library/ruby",
    "count": 36
  },
  {
    "name": "algorithm-rls/oracle",
    "count": 35
  },
  {
    "name": "algorithm-rls/pipeline",
    "count": 35
  },
  {
    "name": "library/buildpack-deps",
    "count": 34
  },
  {
    "name": "frontend/lms-mobile",
    "count": 33
  },
  {
    "name": "backend/coursescriptci",
    "count": 30
  },
  {
    "name": "backend/homepage_staging",
    "count": 30
  },
  {
    "name": "backend-rls/sequence",
    "count": 28
  },
  {
    "name": "backend/link_agent",
    "count": 28
  },
  {
    "name": "backend/link_agent_nginx",
    "count": 27
  },
  {
    "name": "engzo-ci/data-ci",
    "count": 27
  },
  {
    "name": "engzo/slate",
    "count": 26
  },
  {
    "name": "backend/teelog",
    "count": 25
  },
  {
    "name": "platform-rls/logstash",
    "count": 24
  },
  {
    "name": "backend-rls/pushserver",
    "count": 23
  },
  {
    "name": "analysis/service-map",
    "count": 22
  },
  {
    "name": "data-ci/data-ci",
    "count": 21
  },
  {
    "name": "backend/search_web",
    "count": 19
  },
  {
    "name": "library/debian",
    "count": 19
  },
  {
    "name": "backend/audcnv-service",
    "count": 18
  },
  {
    "name": "backend/direct_message_service",
    "count": 18
  },
  {
    "name": "platform-rls/prom-config",
    "count": 18
  },
  {
    "name": "backend-rls/telis",
    "count": 17
  },
  {
    "name": "backend/homepage",
    "count": 17
  },
  {
    "name": "library/codis",
    "count": 17
  },
  {
    "name": "platform/dva",
    "count": 17
  },
  {
    "name": "backend-rls/ielts_server",
    "count": 16
  },
  {
    "name": "backend-rls/westworld",
    "count": 16
  },
  {
    "name": "backend/agora_recorder",
    "count": 16
  },
  {
    "name": "backend/hexley",
    "count": 16
  },
  {
    "name": "frontend/telis-cms",
    "count": 16
  },
  {
    "name": "gcr-io-rls/google_containers-cluster-autoscaler",
    "count": 16
  },
  {
    "name": "backend/entrance_exams_base",
    "count": 15
  },
  {
    "name": "backend/locust",
    "count": 14
  }
]
```

由于 harbor v0.5.0 提供的可用信息太少，若仅基于 pull_count 实现 RP ，则有很大的不准确性；

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/harbor%20v0.5%20top.png)

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/harbor%20v0.5%20repositories.png)

而 harbor v1.2.0 则可以玩出很多花样了

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/harbor%20v1.2.0%20top.png)

![](https://raw.githubusercontent.com/moooofly/ImageCache/master/Pictures/harbor%20v1.2.0%20repositories.png)


## 建议

更新 harbor 版本；

----------


## Task

- [T40042](https://phab.llsapp.com/T40042)
- [T41119](https://phab.llsapp.com/T41119)


