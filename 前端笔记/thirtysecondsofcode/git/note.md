1. merge 与 rebase 区别
   假如你想将分支 feature 合并到分支 master，那么只需执行如下两步即可：
   git checkout master
   git merge feature
   **多了一次合并的历史**
   ![merge](https://i.loli.net/2020/04/06/S5YmfCK7wW1JxTB.png)

   假如你想将分支 feature 合并到分支 master，那么只需执行如下两步即可：
   git checkout master
   git rebase feature
   而 git rebase 会将整个 master 分支移动到 feature 分支的顶端，从而有效地整合了所有 master 分支上的提交。
   ![rebase](https://i.loli.net/2020/04/06/F1qIGKTNDW6aRul.png)

   | 比较 | merge                | rebase         |
   | ---- | -------------------- | -------------- |
   | 优点 | 保留有价值的历史文档 | 删减就繁       |
   | 缺点 | 分支杂乱冗余         | 无法体现时间线 |

   如果项目庞大，需要一个简洁的线性历史树便于 leader 管理，推荐使用 git rebase 。
   如果是小型项目，需要审查历史纪录来便于编写过程报告，则推荐使用 git merge 。
