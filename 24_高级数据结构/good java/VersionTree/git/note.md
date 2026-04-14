- 底层是提交图，上层是分支和引用

---

我建议把 Git 理解成下面 8 层。只要这 8 层通了，Git 基本就通了。

**1. 最底层：Git 不是“版本控制命令集”，而是对象数据库**
Git 的根不是分支、提交、冲突这些用户概念，而是一个内容寻址的对象库。

对象最核心有 4 类：

1. blob
   表示文件内容。
   只管字节，不管文件名、不管路径、不管权限。

2. tree
   表示目录。
   里面记录若干条目，每个条目是：
   模式 + 名字 + 指向的对象
   这个对象可以是 blob，也可以是子 tree。

3. commit
   表示一次历史节点。
   里面至少包含：
   当前快照对应的 tree
   父提交列表
   作者、提交者
   时间
   提交说明

4. tag
   本质是给某个对象起一个稳定名字。
   轻量 tag 更像一个 ref，附注 tag 本身也是对象。

Git 一切高层行为，最后都落回这 4 类对象。

**2. 内容寻址：Git 为什么快、为什么能去重**
Git 的对象不是用“文件名”索引，而是用“内容哈希”索引。
也就是说，同样内容的 blob，会映射到同一个对象。

所以 Git 的一个根本性质是：
对象一旦写入，就不修改。
如果内容变了，不是改旧对象，而是创建新对象。

这带来两个非常强的结果：

1. 不可变
   提交历史天然稳定，已有对象不会被原地篡改。

2. 去重
   同样内容不会重复存很多份。

因此 Git 本质上更像：
一个不可变对象图
而不是
一个会不断覆盖旧状态的数据库。

**3. 提交图：Git 的历史不是链，而是 DAG**
很多人初学 Git，会把历史想成一条线。
这只在最简单情况成立。

真实情况是：
每个 commit 指向一个或多个 parent
因此整体是有向无环图，也就是 DAG。

为什么不是树？
因为 merge commit 有多个 parent。

为什么不是一般图？
因为不会出现祖先回指后代，所以无环。

所以 Git 的时间观不是“线性时间”，而是“图状时间”：

1. 一个提交可以分叉出多个未来
2. 多个未来可以再次汇合
3. 历史不是单一时间轴，而是并行演化的图

这也是为什么你做竞赛里的版本树时，Git 比喻很自然：

1. 如果只有 reset 和继续修改，就是树
2. 如果允许 merge，就从树升级为 DAG

**4. Ref：分支到底是什么**
分支不是一串提交。
分支只是一个名字，指向某个 commit。

例如：
main 指向 commit A
feature 指向 commit B

这两个名字本质上是 refs/heads 下面的引用。

所以 branch 的本质是：
可移动引用

这件事非常重要，因为它解释了很多现象：

1. 建分支很便宜
   只是多一个指针。

2. 删除分支也很轻
   只是删掉一个引用，不是马上删提交对象。

3. reset 本质是移动引用
   不是直接“删历史”。

4. merge 和 rebase 的很多区别，最后都表现在引用怎么移动。

**5. HEAD：你当前站在哪里**
HEAD 是 Git 最核心的“当前上下文”指针。

最常见情况：
HEAD 指向某个 branch ref
比如 HEAD -> refs/heads/main -> commit X

这表示：
我当前在 main 分支上工作，main 当前指向 X。

另一种情况：
HEAD 直接指向某个 commit
这叫 detached HEAD。

这意味着：
你现在站在某个提交上，但不在任何分支名字上。
如果此时继续提交，会生成新 commit，但不会自动挂到某个已有 branch 上，除非你再建分支或移动 ref。

所以：
HEAD 是当前所在位置
branch 是命名引用
commit 是不可变历史节点

这三者不能混。

**6. 三棵树模型：理解 Git 命令最关键的一层**
这是理解 add、commit、reset、restore、checkout 的钥匙。

Git 可以看成同时维护 3 份状态：

1. Commit tree
   当前 HEAD 所指向提交里的快照

2. Index
   暂存区，下一次提交准备采用的快照

3. Working tree
   工作区，磁盘上你正在编辑的文件

很多 Git 困惑，本质上都来自于没有把这三层分开。

例如：
你改了文件，但还没 add
这时：
working tree 变了
index 没变
HEAD 没变

你 add 之后：
working tree 和 index 一致了
HEAD 还没变

你 commit 之后：
index 的内容被写成新 commit
branch 被移动到新 commit

所以：
add 是 working tree -> index
commit 是 index -> commit
checkout/reset/restore 往往是 commit 或 index -> working tree / index

**7. 最常用命令，逐一还原到本质**

**7.1 git add**
本质不是“提交”，而是更新 index。

它做的事是：
把当前工作区某些路径的内容写入暂存区，准备给下一次 commit 使用。

所以 add 的真实意思是：
我要让下一次提交包含这些内容。

**7.2 git commit**
本质是：
把 index 序列化成 tree
创建 commit 对象
移动当前 branch ref 到这个新 commit

注意，它默认不直接看工作区，而是看 index。

所以 commit 提交的是暂存区，不是工作区。

**7.3 git branch**
本质是：
创建一个新的 ref，指向某个 commit。

例如新建 feature 分支：
只是让 refs/heads/feature 指向当前提交。

**7.4 git switch / git checkout**
本质是：
改变 HEAD 指向
然后根据目标状态更新 index 和 working tree

如果切到分支：
HEAD 会重新指向那个分支 ref

如果切到具体提交：
HEAD 直接指向 commit，进入 detached HEAD

所以 checkout 不是“进入另一个副本”，而是“让当前工作区呈现另一个提交快照”。

**7.5 git reset**
这是用户最容易误用但本质最清楚的命令之一。

reset 先做一件核心事情：
移动当前 branch ref

然后根据模式，决定是否同步 index 和 working tree。

1. soft
   只移动 ref

2. mixed
   移动 ref，并重置 index
   默认模式

3. hard
   移动 ref，并重置 index 和 working tree

所以 reset 的关键不是“回退文件”，而是“移动 branch 指针，并可选同步另外两层”。

**7.6 git restore**
restore 是“内容恢复”命令。
它更像在不同状态之间拷文件。

例如：
从 index 恢复工作区
从某个提交恢复工作区
从某个提交恢复到 index

所以 restore 更偏“取内容”
reset 更偏“挪引用”

**7.7 git revert**
revert 不改历史图结构，而是追加一个反向提交。

所以它的本质是：
根据某个旧提交的变化，生成一个抵消它的新提交

这和 reset 完全不同：

1. reset 改引用位置
2. revert 写一个新的纠错历史节点

因此：
公共分支用 revert 更安全
私有历史整理常用 reset 或 rebase

**7.8 git merge**
merge 的本质是：
找共同祖先
计算两边相对祖先的修改
尝试合并
必要时创建一个双亲 commit

如果一边是另一边的祖先，就只需要快进移动 ref，不需要新 commit。
这叫 fast-forward。

如果两边都各自发展过，就要生成 merge commit。
它会有两个 parent。

所以 merge 的本质是：
保留分叉结构，并在图上做汇合。

**7.9 git rebase**
rebase 的本质是：
把一串提交从旧基底上摘下来，重新应用到新基底上。

所以它不是“平移指针”，而是“重写提交”。

为什么 commit id 会变？
因为 commit 内容里包含 parent。
parent 变了，commit 哈希也就变了。

所以 rebase 的本质是：
重放补丁，重建历史

而 merge 的本质是：
保留原历史，再额外连一条汇合边

**7.10 git cherry-pick**
本质上是：
取某个提交的变化，单独重放到当前 HEAD 上

它和 rebase 很像，只是 cherry-pick 是点选提交，rebase 是批量重放一串提交。

**7.11 git stash**
stash 的本质是：
把当前未提交修改临时保存为一组可恢复状态，再把工作区清理出来

它不是魔法抽屉，底层仍然是在借 Git 对象保存快照。

**7.12 git fetch**
本质是：
下载远端对象和远端 refs 的最新信息到本地
但不改你当前 branch

所以 fetch 很像：
更新你对远端世界的认知

**7.13 git pull**
本质不是独立概念，而是：
fetch + merge
或者
fetch + rebase

**7.14 git push**
本质是：
请求远端把某个 ref 移动到某个 commit

所以 push 不是“发文件”，而是“更新远端引用”。

如果远端发现这不是快进更新，默认拒绝。
因为那说明你在改写它已有的公开历史。

force push 的本质就是：
强制远端接受一次 ref 重写

**8. Git 为什么会冲突**
冲突不是“Git 不会合并文本”这么浅。
更本质地说，冲突来自：
两个演化路径对同一逻辑区域施加了不兼容变换，而自动规则无法唯一决定结果。

所以冲突的根不在文件，而在变换组合失败。

三路合并本质上是在做：

1. 找共同祖先 A
2. 看 A 到分支一的变换
3. 看 A 到分支二的变换
4. 尝试把这两个变换组合到一起

如果组合失败，就冲突。

**9. Git 里哪些东西是真，哪些只是用户界面**
这点特别重要。

真正稳定的底层概念是：

1. 对象
2. 引用
3. HEAD
4. 三棵树
5. 提交图

而很多命令名只是用户界面包装。
例如：

1. switch 和 checkout 有一部分能力重叠
2. restore 和 checkout 有一部分能力重叠
3. pull 只是 fetch 加某种整合策略
4. branch -f、reset、update-ref 在某些语义上能表达相近动作

所以你不要执着于“这个命令是不是唯一正确命令”，而要回到：
它到底是在改对象、改引用，还是同步三棵树。

**10. 为什么 Git 适合作为“版本题”的比喻**
因为 Git 最核心的东西正好对应竞赛里的版本树抽象：

1. commit
   对应一个不可变版本节点

2. branch/head
   对应一个当前状态指针

3. reset
   对应跳回某个旧版本

4. commit from parent
   对应从某个旧版本分叉出新版本

5. query 挂在某个版本上
   这不是 Git 原生命令，但和“站在某个 commit 上观察状态”非常接近

所以你现在那套竞赛模板，其实抓住了 Git 最适合抽象成算法结构的那一部分：

1. 不可变版本节点
2. 可移动引用
3. 从父版本派生新版本

它没覆盖工作区、暂存区、远端、merge，这不是缺点，只是它只拿了 Git 里最适合做数据结构比喻的那一层。

**11. 进一步深入：很多人实际上没搞懂的几个点**

**11.1 为什么 commit 是快照，但又常说 rebase 是重放补丁**
因为 Git 对象层保存的是快照。
但在很多历史重写操作里，Git 会把“某个提交相对其父提交的变化”抽出来，当成补丁重新应用。
所以：
存储模型偏快照
变换模型经常用补丁思维

两者并不冲突。

**11.2 为什么删除 branch 不一定删除提交**
因为真正决定对象是否保留的，不是 branch 名字，而是：
还有没有 ref 或其他对象能到达它

一个提交只要还能从某个 ref、tag、reflog 等可达，就不会立刻被清理。

**11.3 reflog 是什么**
reflog 不是提交历史，而是“引用移动日志”。

例如 HEAD 从哪移到哪，main 从哪移到哪，这些动作会记录在 reflog 里。
所以很多“误 reset 了怎么办”，能救回来，不是因为提交历史神奇，而是因为 reflog 还记得你曾经把 ref 指到哪里。

这点非常重要：

1. commit graph 是内容历史
2. reflog 是引用运动历史

**11.4 gc 和 packfile 是什么**
Git 不是一直把所有对象都零散存着。
为了节省空间和提高性能，它会做垃圾回收和打包。

1. 不可达对象最终可能被清理
2. 很多对象会打进 packfile
3. packfile 里会做增量压缩

所以从逻辑模型看 Git 像快照数据库，
从物理存储看 Git 会做大量空间优化。

**12. 把 Git 穷尽到实践层，最有用的最终心法**
真正实战时，不要问：
这个命令我是不是会背。

而要问：
它对下面哪一层动手了。

每次分析命令，只问这 5 个问题：

1. 有没有创建新对象
2. 有没有移动 ref
3. 有没有改变 HEAD
4. 有没有改 index
5. 有没有改 working tree

你能把任何命令回答到这 5 个问题上，Git 就基本通了。

例如：

1. commit
   创建对象，移动 ref，不改 working tree

2. reset --hard
   移动 ref，改 index，改 working tree，不创建新提交

3. cherry-pick
   创建新提交，移动 ref，可能改 working tree 和 index

4. fetch
   更新远端相关 ref 和对象，不动当前工作区

这样理解，比背命令说明强得多。

**13. 如果要给 Git 一个最终定义**
Git 是一个内容寻址的不可变对象数据库。
它用 commit 组织成有向无环图，用 ref 和 HEAD 描述当前语境，再通过 working tree、index、history 三层同步机制，让人类可以编辑、暂存、提交、回退、整合历史。

---

**一层总览**
如果要把这三组问题彻底讲透，可以先用一个统一框架：

Git 里真正会变的东西只有四类。

1. 提交图
   也就是 commit DAG，本身由 commit 对象和 parent 关系组成。
2. 引用
   也就是 refs，比如分支、标签、HEAD。
3. 暂存区
   也就是 index。
4. 工作区
   也就是你磁盘上的文件。

几乎所有命令都可以精确地问这四件事：

1. 有没有创建新 commit
2. 有没有移动 ref
3. 有没有改 index
4. 有没有改 working tree

下面三部分都按这个框架讲。

**一. reset、restore、checkout、revert 的精确区分**

这四个命令最容易混，因为它们都带“回去”“恢复”的味道，但本质完全不同。

先给一句话版本：

1. reset
   核心是移动当前分支引用，再按模式决定是否同步 index 和 working tree。
2. restore
   核心是把某个来源的内容拷贝回 index 或 working tree，不动提交图。
3. checkout
   旧接口，既能切分支/切提交，也能恢复路径，职责太多。
4. revert
   不是回到过去，而是新增一个“反向提交”。

**1. reset**
reset 的真正核心动作是：移动当前分支所指向的 commit。

如果当前 HEAD 在分支上，那么 reset 改的是这个分支 ref。
如果当前是 detached HEAD，那么 reset 改的是 HEAD 自己。

它有三种常见模式：

1. soft
   只移动 ref，不动 index，不动 working tree。
   结果是：
   历史回退了，但你原来提交过的改动都还留在暂存区和工作区里。

2. mixed
   移动 ref，重置 index，不动 working tree。
   这是默认行为。
   结果是：
   历史回退了，暂存区也回到目标提交，但工作区保留你之前的文件内容，所以这些差异会变成 unstaged changes。

3. hard
   移动 ref，重置 index，重置 working tree。
   结果是：
   三层一起回到目标提交。

所以 reset 最精确的定义是：
先改引用位置，再决定把 index 和 working tree 往哪一层同步。

它改的是“你当前历史位置”，不是单纯改文件内容。

**2. restore**
restore 不改提交图，通常也不改分支引用。
它做的是“把某个来源的内容恢复到 index 或 working tree”。

最常见两类：

1. 恢复工作区
   把 index 或某个提交里的文件内容拷到 working tree

2. 恢复暂存区
   把某个提交里的文件内容拷到 index

所以 restore 的本质是：
内容搬运命令

它问的不是“历史该指到哪”，而是“这个文件现在该长什么样”。

这就是它和 reset 的根本区别：

1. reset 偏 ref movement
2. restore 偏 content restore

**3. checkout**
checkout 是 Git 早期的大一统命令，功能太杂，所以后来拆出 switch 和 restore。

它历史上同时干两类事：

1. 切换分支或切换到某个提交
   这时它会改 HEAD，并通常更新 index 和 working tree。

2. 从某个提交或 index 恢复路径
   这时它像 restore。

所以 checkout 的问题不是原理错，而是职责过多。
它把“切换位置”和“恢复文件”混在同一个命令入口里。

精确地说：

1. checkout 分支
   本质是改 HEAD 指向，并把工作现场同步到目标提交

2. checkout 提交
   本质是进入 detached HEAD，并同步工作现场

3. checkout 路径
   本质是恢复文件内容

所以现在更推荐：

1. switch 负责“站到哪里”
2. restore 负责“文件内容恢复”

**4. revert**
revert 最容易被误解成“撤销历史”，其实不是。

revert 完全不移动历史位置。
它会新建一个提交，这个提交的内容等于“把某个旧提交引入的变化反过来做一遍”。

所以：

1. reset
   改的是引用位置
2. revert
   改的是最新状态，方式是追加一个新提交

这就是为什么：

1. reset 更适合整理本地历史
2. revert 更适合公共分支纠错

因为 revert 不会改写已有公开历史，只是加一个新节点。

**把四者放在一起看**
可以用下面这个判断法：

1. 你要改历史指向
   用 reset

2. 你要恢复文件内容，但不想改历史位置
   用 restore

3. 你要切换分支或切换到某个提交
   用 switch 或 checkout

4. 你要公开地“反做某个已提交改动”
   用 revert

**图结构角度的区别**

1. reset
   通常不生成新 commit，只移动 ref
2. restore
   不生成新 commit，不移动提交图
3. checkout
   通常不生成新 commit，只改变 HEAD 所在语境
4. revert
   生成一个新 commit，图变长，但不改旧图

**二. merge、rebase、cherry-pick、reflog 在图结构上做了什么**

这一组的重点不是“文件怎么变”，而是“图怎么变”。

**1. merge**
merge 的目标是把两条历史线汇合起来。

假设有：

A -> B -> C main
\
 D -> E feature

如果在 main 上 merge feature，Git 会找共同祖先 B，然后尝试把 B 到 C 的变化和 B 到 E 的变化合起来。

结果有两种：

1. fast-forward
   如果 main 其实没自己往前走，只是 feature 更靠前，那么 main 直接指到 feature 的末端，不会产生 merge commit。

2. true merge
   如果两边都发展了，就生成一个新的 merge commit，比如 M，它有两个 parent：
   一个 parent 是 C
   另一个 parent 是 E

图变成：

A -> B -> C ------\
 \ M
D -> E ---/

所以 merge 的图结构本质是：
保留分叉历史，并增加一个汇合节点。

**2. rebase**
rebase 不是汇合，而是重放。

还是上面的例子：

A -> B -> C main
\
 D -> E feature

如果你对 feature 做 rebase 到 C 上，本质是：

1. 找出 feature 相对共同祖先 B 的独有提交 D、E
2. 从新基底 C 开始
3. 把 D 的变化重新应用，生成 D'
4. 把 E 的变化重新应用，生成 E'

最后变成：

A -> B -> C -> D' -> E'

原来的 D、E 不会被原地修改，它们仍存在，只是 branch 会改去指向 E'。
所以 rebase 的本质是：
重建一条新的线性历史

这就是为什么 rebase 会改写历史。
不是“提交内容变了”，而是“父节点变了”，所以 commit 身份也变了。

**3. cherry-pick**
cherry-pick 是局部重放。

如果某个提交 X 在别的分支上，你现在想把它单独拿过来，Git 会：

1. 取出 X 相对于其 parent 的变化
2. 在当前 HEAD 上重放
3. 生成一个新提交 X'

所以 cherry-pick 的图结构变化是：
当前分支新增一个新节点，但这个新节点不是原来的 X，而是一个新 commit。

所以：

1. merge 是保留整条分叉并汇合
2. rebase 是整串重放
3. cherry-pick 是单点重放

**4. reflog**
reflog 不是提交图的一部分。
它记录的是引用怎么移动过。

这一点非常关键。

提交图回答的是：
有哪些 commit，它们 parent 关系是什么

reflog 回答的是：
HEAD、main 这些 ref 以前指过哪里

比如你做了：

1. commit
2. reset --hard 到旧提交
3. 又切到别的分支

提交图只反映当前对象之间的 parent 关系。
但 reflog 会记：
HEAD 从哪到哪
main 从哪到哪

所以 reflog 更像“引用运动历史”，不是“内容历史”。

这也是很多误操作还能救回来的原因：
提交本身可能暂时没 ref 指着了，但 reflog 还记得你前一秒把某个 ref 指到过它。

**图结构上怎么理解 reflog**
它不是 DAG 的边，也不是新节点。
它更像 DAG 外部的一份日志：

1. 某时刻 ref R 指向 commit X
2. 某时刻 ref R 又指向 commit Y

所以 reflog 是对 ref movement 的时间日志，不是 commit graph 的结构部分。

**把这四个命令放一起**

1. merge
   保留分叉，必要时新增双亲节点
2. rebase
   重放原提交，重建一条新链
3. cherry-pick
   重放单个提交，生成一个新节点
4. reflog
   不改图结构，只记录 ref 怎么移动

**三. Git 内部对象、refs、packfile、gc、index 文件格式到底是什么**

这一部分是 Git 真正的底层。

**1. 对象格式**
Git 的对象底层存储形式可以理解成：

类型 + 空格 + 大小 + 空字符 + 内容

例如一个 blob 对象，逻辑上是：

blob 123\0<文件字节内容>

然后这段内容被哈希，得到对象 id。
老仓库默认是 SHA-1，新模式也可以用 SHA-256。

同理：

1. blob
   内容就是文件原始字节

2. tree
   内容是若干目录条目，每个条目包含：
   文件模式
   文件名
   对应对象 id

3. commit
   内容是文本头部加正文，大概包括：
   tree 某个 tree id
   parent 若干 parent id
   author ...
   committer ...
   空行
   提交说明

4. tag
   内容类似：
   object 某对象 id
   type 对象类型
   tag 标签名
   tagger ...
   空行
   标签说明

所以 commit 不是“神秘记录”，它本质上也是一段有固定结构的文本对象。

**2. loose object**
对象刚创建时，通常先以 loose object 存在。
就是单个对象单独压缩后存放在 .git/objects 目录里。

路径规则是：
对象 id 前两位做目录名
剩余部分做文件名

好处是创建简单。
缺点是对象多了会很碎。

**3. refs**
refs 本质上是“名字 -> 对象 id”的映射。

常见分类：

1. refs/heads/\*
   本地分支

2. refs/tags/\*
   标签

3. refs/remotes/\*
   远端跟踪分支

HEAD 自身通常不是普通 ref 文件，而是一个特殊引用，里面常写：
ref: refs/heads/main

表示 HEAD 目前附着在 main 上。

如果是 detached HEAD，HEAD 里可能直接写 commit id。

所以 ref 的本质不是“分支结构”，而是“可命名指针”。

**4. packed-refs**
当 refs 很多时，不会总是一个 ref 一个小文件。
Git 会把一批 refs 合并进一个 packed-refs 文件。

逻辑上仍然是“名字 -> 对象 id”，只是存储更紧凑。

**5. packfile**
这是 Git 的性能核心之一。

如果所有对象都散着存，仓库会非常碎，读写慢，空间也浪费。
所以 Git 会把很多对象打包成 packfile。

packfile 做两件大事：

1. 把大量对象装进少量大文件
2. 使用增量压缩

例如两个相近版本的大文件，不会各存一整份。
packfile 可能保存一个基对象，再保存另一个对象相对于它的 delta。

所以从逻辑模型看：
Git 像快照系统

从物理存储看：
Git 会做大量基于差异的空间优化

这两个层次不要混。

pack 通常伴随一个 idx 文件，用来快速定位包内对象。

**6. gc**
gc 就是垃圾回收和仓库整理。

它主要干几件事：

1. 找不可达对象
   如果某些对象不再能从任何 ref、tag、reflog 等入口到达，它们未来可能被清理。

2. repack
   把零散对象重新打包成 packfile

3. prune
   删除 不可达对象

4. 压缩和整理元数据
   提升访问效率

所以 gc 不是“修仓库”的附属功能，而是 Git 持久维持性能和空间效率的重要机制。

**7. index 文件**
index 是 Git 最容易被低估、但最关键的文件之一。
它不是“简单的待提交文件列表”，而是一个二进制缓存结构。

它大致记录：

1. 路径名
2. 对应 blob 的对象 id
3. 文件模式
4. 文件 stat 信息
   例如 mtime、ctime、inode、size
5. 各种标志位
   例如阶段信息、扩展标记等

index 的意义有三层：

1. 表示“下一次提交准备长什么样”
2. 作为 working tree 和 commit tree 之间的中间层
3. 作为性能缓存，避免每次都全量重扫所有文件内容

所以 index 既是语义层的“暂存区”，也是工程层的“高性能缓存”。

**8. index 里的 stage**
在普通情况下，一个路径在 index 里只有一个条目。
但在 merge 冲突时，一个路径可以同时有多个 stage：

1. stage 1
   共同祖先版本

2. stage 2
   当前分支版本

3. stage 3
   对方分支版本

所以冲突时 index 其实不是简单“待提交列表”，而是一个能装多版本候选的结构。

这能解释为什么 merge 冲突时 Git 还能知道三方信息。

**9. 可达性**
Git 里很多“为什么还没删除”或“为什么找不到了”的问题，都和可达性有关。

一个对象只要还能从以下入口追到，就算可达：

1. branch ref
2. tag
3. remote ref
4. reflog 中仍保留的引用历史
5. 某些临时保留结构

所以对象是否存在、是否会被 gc 清掉，不取决于“你眼前看不看得到”，而取决于“是否仍然可达”。

**10. 最终把这些拼起来**
如果把 Git 从底到顶串成一条线，就是：

1. working tree
   你正在编辑的文件

2. index
   下一次提交的候选快照 + 状态缓存

3. tree
   一次目录快照

4. commit
   一次历史节点，指向 tree 和 parent

5. refs / HEAD
   告诉你当前名字和当前位置

6. packfile / gc
   负责把底层对象高效地存起来

所以 Git 不是“命令系统”，而是：

1. 一个不可变对象库
2. 一个引用系统
3. 一个三层状态同步系统
4. 一套空间优化与回收机制

**把三部分合成一句终极理解**

1. reset、restore、checkout、revert
   区别在于它们分别操作的是 ref、内容、当前位置还是新提交

2. merge、rebase、cherry-pick、reflog
   区别在于它们分别是在保留图结构、重写图结构、局部重放历史，还是记录 ref 运动

3. objects、refs、packfile、gc、index
   构成了 Git 真正的底层机器：对象负责历史内容，refs 负责命名与定位，index 负责中间态与缓存，pack/gc 负责物理存储效率
