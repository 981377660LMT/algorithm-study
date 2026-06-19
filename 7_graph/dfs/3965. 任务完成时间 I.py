# 3965. 任务完成时间 I
# https://leetcode.cn/problems/finish-time-of-tasks-i/description/
# 给你一个整数 n，表示项目中的任务数量，编号从 0 到 n - 1。这些任务以任务 0 为根的 树 的形式连接。这由一个长度为 n - 1 的二维整数数组 edges 表示，其中 edges[i] = [ui, vi] 表示任务 ui 是任务 vi 的父节点。
#
# 同时给你一个长度为 n 的数组 baseTime，其中 baseTime[i] 表示完成任务 i 所需的时间。
#
# 每个任务的 完成时间 计算如下：
#
# 叶子任务：完成时间为 baseTime[i]。
# 非叶子任务：
# 令 earliest 为其子节点中的 最小 完成时间，latest 为其子节点中的 最大 完成时间。
# 令 ownDuration 为 (latest - earliest) + baseTime[i]。
# 任务 i 的完成时间为 latest + ownDuration。
# 返回根任务 0 的完成时间。

from typing import List

INF = int(2e18)


class Solution:
    def finishTime(self, n: int, edges: List[List[int]], baseTime: List[int]) -> int:
        adjList = [[] for _ in range(n)]
        for x, y in edges:
            adjList[x].append(y)

        def dfs(cur: int) -> int:
            if not adjList[cur]:
                return baseTime[cur]
            earliest, latest = INF, 0
            for next_ in adjList[cur]:
                childRes = dfs(next_)
                earliest = min(childRes, earliest)
                latest = max(childRes, latest)
            return latest * 2 - earliest + baseTime[cur]

        return dfs(0)
