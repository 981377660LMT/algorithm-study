from typing import List

# 返回一个大小为 n 的数组，其中 ans[i] 表示第 i 个节点的子树中与节点 i 标签相同的节点数。

# 思路：dfs后序遍历
# 先让子子孙孙统计好，自己再动手收租(有点像发leetcoin那道题)


class Solution:
    def countSubTrees(self, n: int, edges: List[List[int]], labels: str) -> List[int]:
        adjList = [[] for _ in range(n)]
        for x, y in edges:
            adjList[x].append(y)
            adjList[y].append(x)

        # 每个结点都有counter(有点像trie了)
        counter = [[0 for _ in range(26)] for _ in range(n)]

        def dfs(root: int, parent: int) -> None:
            label = ord(labels[root]) - ord('a')
            counter[root][label] = 1
            for next in adjList[root]:
                if next == parent:
                    continue
                dfs(next, root)  # 先让子子孙孙统计好，自己再动手收租(有点像发leetcoin那道题)
                for i in range(26):
                    counter[root][i] += counter[next][i]

        dfs(0, -1)

        res = []
        for i in range(n):
            label = ord(labels[i]) - ord('a')
            res.append(counter[i][label])
        return res


print(
    Solution().countSubTrees(
        n=7, edges=[[0, 1], [0, 2], [1, 4], [1, 5], [2, 3], [2, 6]], labels="abaedcd"
    )
)
# 输出：[2,1,1,1,1,1,1]
# 解释：节点 0 的标签为 'a' ，以 'a' 为根节点的子树中，节点 2 的标签也是 'a' ，因此答案为 2 。注意树中的每个节点都是这棵子树的一部分。
# 节点 1 的标签为 'b' ，节点 1 的子树包含节点 1、4 和 5，但是节点 4、5 的标签与节点 1 不同，故而答案为 1（即，该节点本身）。

