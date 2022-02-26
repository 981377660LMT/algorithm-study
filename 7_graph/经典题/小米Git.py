#
# 代码中的类名、方法名、参数名已经指定，请勿修改，直接返回方法规定的值即可
#
#
# @param matrix string字符串一维数组
# @param versionA int整型
# @param versionB int整型
# @return int整型
# Git 是一个常用的分布式代码管理工具，Git 通过树的形式记录文件的更改历史（例如示例图），树上的每个节点表示一个版本分支，工程师经常需要找到两个分支的最近的分割点。
# 例如示例图中 3,4 版本的分割点是 1。3,5 版本的分割点是 0。
# 给定一个用邻接矩阵 matrix 表示的树，请你找到版本 versionA 和 versionB 最近的分割点并返回编号。
#
from typing import List, DefaultDict
from collections import defaultdict


class Solution:
    def Git(self, matrix: List[str], versionA: int, versionB: int) -> int:
        # write code here
        def dfs(cur: int, parent: int, depth: int) -> None:
            parentMap[cur] = parent
            levelMap[cur] = depth
            for next in adjMap[cur]:
                if next == parent:
                    continue
                dfs(next, cur, depth + 1)

        def LCA(
            root1: int, root2: int, level: DefaultDict[int, int], parent: DefaultDict[int, int]
        ) -> int:
            if level[root1] < level[root2]:
                root1, root2 = root2, root1
            diff = level[root1] - level[root2]
            for _ in range(diff):
                root1 = parent[root1]
            while root1 != root2:
                root1 = parent[root1]
                root2 = parent[root2]
            return root1

        root1, root2 = versionA, versionB
        adjMap = defaultdict(list)
        for cur, row in enumerate(matrix):
            nums = [int(char) for char in row]
            for next, isConnected in enumerate(nums):
                if isConnected:
                    adjMap[cur].append(next)
                    adjMap[next].append(cur)

        levelMap, parentMap = defaultdict(lambda: -1), defaultdict(lambda: -1)
        dfs(0, -1, 0)

        lca = LCA(root1, root2, levelMap, parentMap)
        return lca

