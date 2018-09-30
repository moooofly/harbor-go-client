# [Garbage collection](https://docs.docker.com/registry/garbage-collection/)

> As of v2.4.0 a **garbage collector** command is included within the registry binary. This document describes what this command does and how and why it should be used.

- 从 [v2.4.0](https://github.com/docker/distribution/tree/v2.4.0) 版本开始，registry 可执行程序提供了 garbage collector 命令；
- 本文主要说明该命令做了哪些事情，以及为何要使用该命令；

> 项目 docker/docker-registry 更名为 docker/distribution

## About garbage collection

> In the context of the Docker registry, **garbage collection** is the process of removing **blobs** from the filesystem when they are no longer referenced by a **manifest**.
> **Blobs** can include both **layers and **manifests**.

- 在 Docker registry 上下文下，garbage collection 对应从文件系统中移除 blobs 的处理过程（当且仅当 blobs 不在被 manifest 所引用时）；
- Blobs 既可以是 layers ，也可以是 manifests ；

> **Registry data** can occupy considerable amounts of disk space.
> In addition, garbage collection can be a security consideration, when it is desirable to ensure that certain layers no longer exist on the filesystem.

- Registry 数据会占用大量的磁盘空间；
- 除此之外，garbage collection 还可以作为一种安全考量，用于确保特定的 layers 在文件系统中确实不再存在；

## Garbage collection in practice

> **Filesystem layers are stored by their content address in the Registry**.
> This has many advantages, one of which is that data is stored once and referred to by manifests.
> See [here](https://docs.docker.com/registry/compatibility/#content-addressable-storage-cas) for more details.

- 文件系统 layers 是基于其 content address 保存在 Registry 中的；
- 这种方式有很多优点，其一就是 data 只需要保存一份，后续即可通过不同的 manifests 引用；

> **Layers are therefore shared amongst manifests**. **each manifest maintains a reference to the layer**.
> As long as a layer is referenced by one manifest, it cannot be garbage collected.

- Layers 是在不同 manifests 之间是共享的；
- 每个 manifest 都维护一个针对特定 layer 的引用；
- 只要一个 layer 仍旧被一个 manifest 所引用，该 layer 就不能被 GC ；

> **Manifests and layers can be `deleted` with the registry API** (refer to the API documentation [here](https://docs.docker.com/registry/spec/api/#deleting-a-layer)
> and [here](https://docs.docker.com/registry/spec/api/#deleting-an-image) for details).
> This API removes references to the target and makes them eligible for garbage collection. It also makes them unable to be read via the API.

- Manifests 和 layers 都可以通过 registry API 进行 `deleted` ；
- 该 API 能够移除针对 target 的引用，使其能够被 GC ；该 API 同样会令 target 无法通过 API 被读取；

> **If a layer is deleted**, it will be removed from the filesystem when garbage collection is run.
> **If a manifest is deleted**, the layers to which it refers will be removed from the filesystem if no other manifests refers to them.

- 如果 layer 被删除了，则可以在 GC 运行的时候，将其从文件系统移除；
- 如果 manifest 被删除了，则被其引用的 layers 在没有被其他 manifests 所引用的情况下，会被删除；


### Example

> In this example manifest A references two layers: `a` and `b`. Manifest `B` references layers `a` and `c`.
> In this state, nothing is eligible for garbage collection:

在这个例子中，manifest A 引用了 layer `a` 和 layer `b` ，manifest B 引用了 layer `a` 和 layer `c` ；在这种情况下，没有任何 blobs 能够被 GC ；

```
A -----> a <----- B
    \--> b     |
         c <--/
```

> Manifest B is deleted via the API:

通过 API 删除 Manifest B ；

```
A -----> a     B
    \--> b
         c
```

> In this state layer `c` no longer has a reference and is eligible for garbage collection.
> Layer `a` had one reference removed but will not be garbage collected as it is still referenced by manifest `A`.
> The blob representing manifest `B` will also be eligible for garbage collection.

- 在这种状态下，layer `c` 已经没有被引用，因此符合被 GC 的要求；
- 而 Layer `a` 有一个引用被移除了，但仍旧存在 manifest `A` 的引用，因此无法被 GC ；
- 另外，代表 manifest `B` 的 blob 同样符合被 GC 的要求；

> After garbage collection has been run, manifest `A` and its blobs remain.

在 GC 运行过后，manifest `A` 和其引用的 blobs 仍旧存在；

```
A -----> a
    \--> b
```


### More details about garbage collection

> Garbage collection runs in two phases.
>
> - First, in the **'mark' phase**, the process **scans all the manifests in the registry**. From these manifests, it **constructs a set of content address digests**.
>   This set is the 'mark set' and denotes the set of blobs to *not* delete.
> - Secondly, in the **'sweep' phase**, the process **scans all the blobs** and if a blob's content address digest is not in the mark set, the process will delete it.

GC 运行分为两个阶段：

- 首先，为 'mark' 阶段，GC 进程会扫描 registry 中所有 manifests ，并构建出一组相应的 content address digests ；该组实际为 'mark set' ，标识不能被删除的 blobs 组；
- 其次，为 'sweep' 阶段，GC 进程会扫描所有的 blobs ，若某个 blob 的 content address digest 不在 mark set 中，则将其删除；


>> **Note**: **You should ensure that the registry is in read-only mode or not running at all**.
>> If you were to upload an image while garbage collection is running, there is the risk that the image's layers will be mistakenly deleted, leading to a corrupted image.

**注意**：你需要确保 registry 处于 read-only 模式，或者干脆处于未运行状态；如果你在 GC 运行过程中上传了 image ，则有很大可能会导致 image 的 layers 被错误的删除，导致 image 损毁；

> This type of garbage collection is known as stop-the-world garbage collection.
> In future registry versions the intention is that garbage collection will be an automated background action and this manual process will no longer apply.

该 GC 的类型为 stop-the-world GC ，在未来的 registry 版本中，GC 将会成为一种全自动的后台行为，故当前这种手动方式将不复存在；


## Run garbage collection

> Garbage collection can be run as follows

可以按照如下命令运行 GC

```
bin/registry garbage-collect [--dry-run] /path/to/config.yml
```

> The `garbage-collect` command accepts a `--dry-run` parameter, which will print the progress of the mark and sweep phases without removing any data.
> Running with a log level of `info` will give a clear indication of what will and will not be deleted.

- GC 命令接受一个 `--dry-run` 参数，能够输出 mark 和 sweep 阶段的处理过程，而不移除任何实际数据；
- 若指定日志级别 `info` 运行上述命令，则能够清楚的看到哪些内容会被删除，哪些内容不会被删除；

> The config.yml file should be in the following format:

配置文件 config.yml 的格式如下：

```
version: 0.1
storage:
  filesystem:
    rootdirectory: /registry/data
```

> _Sample output from a dry run garbage collection with registry log level set to `info`_

```
hello-world
hello-world: marking manifest sha256:fea8895f450959fa676bcc1df0611ea93823a735a01205fd8622846041d0c7cf
hello-world: marking blob sha256:03f4658f8b782e12230c1783426bd3bacce651ce582a4ffb6fbbfa2079428ecb
hello-world: marking blob sha256:a3ed95caeb02ffe68cdd9fd84406680ae93d633cb16422d00e8a7c22955b46d4
hello-world: marking configuration sha256:690ed74de00f99a7d00a98a5ad855ac4febd66412be132438f9b8dbd300a937d
ubuntu

4 blobs marked, 5 blobs eligible for deletion
blob eligible for deletion: sha256:28e09fddaacbfc8a13f82871d9d66141a6ed9ca526cb9ed295ef545ab4559b81
blob eligible for deletion: sha256:7e15ce58ccb2181a8fced7709e9893206f0937cc9543bc0c8178ea1cf4d7e7b5
blob eligible for deletion: sha256:87192bdbe00f8f2a62527f36bb4c7c7f4eaf9307e4b87e8334fb6abec1765bcb
blob eligible for deletion: sha256:b549a9959a664038fc35c155a95742cf12297672ca0ae35735ec027d55bf4e97
blob eligible for deletion: sha256:f251d679a7c61455f06d793e43c06786d7766c88b8c24edf242b2c08e3c3f599
```
