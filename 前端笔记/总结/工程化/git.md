1.  撤销修改
    1.  只在工作区修改了，没提交 add 到暂存区
        git checkout -- index.html 撤销工作区修改(其实 git checkout -- file 就是**用暂存区的版本来代替工作区的版本**)
    2.  修改提交 add 到暂存区之后 又修改了
        git reset HEAD index.html 将暂存区的修改撤销掉，重新放回工作区然后重复 1）
    3.  已经提交到本地的版本库分支 master 上了（**前提是没推送到远程版本库**）
        版本回退，这时候需要使用命令 **git reset --hard 编号**（编号就是上节说的 commit id（版本号） 不需要写出完整的版本号 一般前 8 位以上就 ok 了）
        **git reset HEAD 文件** 命令用于取消已缓存的内容。
2.  标签
    `git tag 标签名` 默认是为我们最后一次 commit 创建标签
    如果想为每次 commit 创建标签怎么办呢
    很简单 拿到 commit id（版本号）就行
    `$ git tag 标签名 版本号`
