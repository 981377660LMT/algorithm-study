# 拓扑排序树
# !树的拓扑排序,从叶子开始向中间侵蚀,
# 记录每个结点的层数信息

from collections import deque
from itertools import accumulate
from typing import List


def topoSortTree(n: int, tree: List[List[int]], deg: List[int], directed=False) -> List[int]:
    """从所有叶子结点向中心侵蚀的拓扑排序过程.
    返回每个结点是在第几轮遍历到的(0-based).
    """
    leafDeg = 1 if not directed else 0
    res = [0] * n
    queue = deque()
    for i in range(n):
        if deg[i] == leafDeg:
            queue.append(i)

    round = 0
    while queue:
        len_ = len(queue)
        for _ in range(len_):
            u = queue.popleft()
            res[u] = round
            for v in tree[u]:
                deg[v] -= 1
                if deg[v] == leafDeg:
                    queue.append(v)
        round += 1
    return res


if __name__ == "__main__":
    # https://leetcode.cn/problems/collect-coins-in-a-tree/
    # 树中收集金币的最短路径
    # 一开始，你需要选择树中任意一个节点出发。你可以执行下述操作任意次：
    #  !- 收集距离当前节点距离为 2 以内的所有金币，或者
    #  !- 移动到树中一个相邻节点。
    # 你需要收集树中所有的金币，并且回到出发节点，请你返回最少经过的边数。

    # !进阶问题:对每个距离为2,3,...,k,求出所有答案
    # !统计每条边端点的level最小值
    class Solution:
        def collectTheCoins(self, coins: List[int], edges: List[List[int]]) -> int:
            if sum(coins) <= 1:
                return 0
            n = len(coins)
            adjList = [[] for _ in range(n)]
            deg = [0] * n
            for cur, next in edges:
                adjList[cur].append(next)
                adjList[next].append(cur)
                deg[cur] += 1
                deg[next] += 1

            # !第一次拓扑排序:去除没有金币的叶子
            queue = deque()
            for i in range(n):
                if deg[i] == 1 and coins[i] == 0:
                    queue.append(i)
            while queue:
                cur = queue.popleft()
                for next in adjList[cur]:
                    deg[next] -= 1
                    if deg[next] == 1 and coins[next] == 0:
                        queue.append(next)

            # !第二次拓扑排序:记录每个点的层数
            level = [0] * n
            queue = deque()
            for i in range(n):
                if deg[i] == 1 and coins[i] == 1:
                    queue.append(i)
            while queue:
                cur = queue.popleft()
                for next in adjList[cur]:
                    deg[next] -= 1
                    if deg[next] == 1:
                        queue.append(next)
                        level[next] = level[cur] + 1

            # 进阶问题:对每个距离为2,3,...,k,求出所有答案
            edgeLevel = [0] * n
            for x, y in edges:
                cur = min(level[x], level[y])
                edgeLevel[cur] += 1
            sufSum = ([0] + list(accumulate(edgeLevel[::-1])))[::-1]
            return 2 * sufSum[2]

    # coins = [1,0,0,0,0,1], edges = [[0,1],[1,2],[2,3],[3,4],[4,5]]
    print(Solution().collectTheCoins([1, 0, 0, 0, 0, 1], [[0, 1], [1, 2], [2, 3], [3, 4], [4, 5]]))
