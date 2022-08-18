分支合并没什么困难的，主要就是 merge 和 rebase 两种方式。
本文就用 Git 的 **rebase** 工作方式引出一个经典的算法问题：
最近公共祖先（Lowest Common Ancestor，简称 LCA）。

git pull 这个命令，我们经常会用，
它**默认是使用 merge**方式将远端别人的修改拉到本地；
如果带上上参数 git pull -r，
就会使用 rebase 的方式将远端修改拉到本地。

这二者最直观的区别就是：merge 方式合并的分支会有很多「分叉」，
而 rebase 方式合并的分支就是一条直线。
对于多人协作，merge 方式并不好
所以一般来说，实际工作中更推荐使用 rebase 方式合并代码。

rebase 是如何将两条不同的分支合并到同一条分支的呢
首先，找到这两条分支的最近公共祖先 LCA，然后从 master 节点开始，重演 LCA 到 dev 几个 commit 的修改，如果这些修改和 LCA 到 master 的 commit 有冲突，就会提示你手动解决冲突，最后的结果就是把 dev 的分支完全接到 master 上面。
