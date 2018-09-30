# Harbor 镜像删除问题

**原始需求**：使用过程中，一般情况都只是向 harbor 上传 image ，但是几乎没有人会主动进行删除，导致 image 积累越来越多；所以希望可以按照某种策略把最近不常用的 image 删除掉；


----------

[官方文档](https://github.com/vmware/harbor/blob/master/docs/user_guide.md#deleting-repositories)给出的删除说明如下：

> Repository deletion runs in two steps.
>
> - **First, delete a repository in Harbor's UI**. This is `soft deletion`. You can delete the entire repository or just a tag of it. After the soft deletion, the repository is no longer managed in Harbor, however, the files of the repository still remain in Harbor's storage.
>
>> **CAUTION**: If both tag A and tag B refer to the same image, after deleting tag A, B will also get deleted.
>
> - **Next, delete the actual files of the repository using the registry's garbage collection(GC)**. Make sure that no one is pushing images or Harbor is not running at all before you perform a GC. If someone were pushing an image while GC is running, there is a risk that the image's layers will be mistakenly deleted which results in a corrupted image. So before running GC, a preferred approach is to stop Harbor first.
>
> Run the below commands on the host which Harbor is deployed on to **preview** what files/images will be affected:
>
> ```
> $ docker-compose stop
> $ docker run -it --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect --dry-run /etc/registry/config.yml
> ```
>
> NOTE: The above option "--dry-run" will print the progress without removing any data.
>
> Verify the result of the above test, then use the below commands to **perform garbage collection** and restart Harbor.
>
> ```
> $ docker run -it --name gc --rm --volumes-from registry vmware/registry:2.6.2-photon garbage-collect  /etc/registry/config.yml
> $ docker-compose start
> ```
>
> For more information about GC, please see GC.

简单来说，Repository 删除分两步：

- 首先，通过 Harbor UI 或者 RESTful API 删除目标 repository 或者目标 repository 的某些 tags ；此时的删除被称作“soft deletion”；在 soft deletion 之后，repository 不再直接被 harbor 所管理，但 repository 对应的文件内容仍旧存在于 harbor 所使用的存储设备中；
- 其次，通过 repository 提供的 GC 命令删除 repository 相关的文件系统文件；在执行 GC 命令时，要求 harbor 处于只读状态，或者干脆不要运行；否则，可能发生 image 的某些 layers 被错误的删除的情况，进而导致 image 不可用（相关 issues ：[goharbor/harbor/issues#4214](https://github.com/goharbor/harbor/issues/4214) [goharbor/harbor/issues#5167](https://github.com/goharbor/harbor/issues/5167)）；


这里有一个问题需要弄清楚：**即删除镜像的过程中需要注意哪些问题**？

详见《[Garbage Collection](https://github.com/moooofly/harbor-go-client/blob/master/docs/Garbage%20Collection.md)》

相关链接：

- [Deleting repositories](https://github.com/vmware/harbor/blob/master/docs/user_guide.md#deleting-repositories)
- [issues/1168](https://github.com/vmware/harbor/issues/1168)
- [issues/2287](https://github.com/vmware/harbor/issues/2287)
- [issues/737](https://github.com/vmware/harbor/issues/737)
- [issues/3456](https://github.com/vmware/harbor/issues/3456)
