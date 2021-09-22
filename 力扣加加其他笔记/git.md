1. 再简单提一些常见场景， 比如说 commit 完之后，突然发现一些错别字需要修改，又不想为改几个错别字而新开一个 commit 到 history 区，那么就可以使用下面这个命令：

```Shell
$ git commit --amend
```

这样就是把错别字的修改和之前的那个 commit 中的修改合并，作为一个 commit 提交到 history 区

2. 如何把 stage 中的修改还原到 work dir 中。

```
touch a.txt b.txt
$ git add .
$ git status
$ echo hello world >> a.txt
$ git status
```

现在，我后悔了，我认为不应该修改 a.txt，我想把它还原成 stage 中的空文件，怎么办？
答案是，使用 checkout 命令：
**git checkout a.txt**

3. 将 history 区的文件还原到 stage 区。
   比如说我用了一个 git add . 一股脑把所有修改加入 stage，但是突然想起来文件 a.txt 中的代码我还没写完，不应该把它 commit 到 history 区，所以我得把它从 stage 中撤销，等后面我写完了再提交。
   $ echo aaa >> a.txt; echo bbb >> b.txt;
   $ git add .
   $ git status
   使用
   **$ git reset a.txt**

4. 将 history 区的历史提交还原到 work dir 中。
   比如我从 GitHub 上 clone 了一个项目，然后乱改了一通代码，结果发现我写的代码根本跑不通，于是后悔了，干脆不改了，我想恢复成最初的模样，怎么办？
   依然是使用 checkout 命令，但是和之前的使用方式有一些不同：
   **$ git checkout HEAD .**
   Updated 12 paths from d480c4f
   这样，work dir 和 stage 中所有的「修改」都会被撤销，恢复成 HEAD 指向的那个 history commit
   只要找到任意一个 commit 的 HASH 值，checkout 命令可就以将文件恢复成任一个 history commit 中的样子：
   **$ git checkout 2bdf04a some_test.go**
   Updated 1 path from 2bdf04a

三、其他技巧
需求一，**合并多个 commit**。
比如说我本地从 17bd20c 到 HEAD 有多个 commit，但我希望把他们合并成一个 commit 推到远程仓库，这时候就可以使用 reset 命令：
**$ git reset 17bd20c**
$ git add .
$ git commit -m 'balabala'
reset 命令的作用，相当于把 HEAD 移到了 17bd20c 这个 commit，而且不会修改 work dir 中的数据，所以只要 add 再 commit，就相当于把中间的多个 commit 合并到一个了。

需求二，由于 HEAD 指针的回退，导致有的 commit 在 git log 命令中无法看到，**怎么得到它们的 Hash 值呢**
只要你不乱动本地的 .git 文件夹，任何修改只要提交到 commit history 中，都永远不会丢失，看不到某些 commit 只是因为它们不是我们**当前 HEAD 位置的「历史」提交**，我们可以使用如下命令查看操作记录：
**$ git reflog**
比如 reset，checkout 等等关键操作都会在这里留下记录，所有 commit 的 Hash 值都能在这里找到，所以如果你发现有哪个 commit 突然找不到了，一定都可以在这里找到。

需求三，**怎么解决冲突**
比较流行的代码编辑器或者 IDE 都会集成方便的可视化 Git 工具，至于解决冲突，可视化的表现方式不是比你在命令行里 git diff 看半天要清晰明了得多？只需要点点点就行了。
