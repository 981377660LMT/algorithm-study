# 无向图存在欧拉路径<=> 度数为奇数的点只能有 0 或 2 个，不存在出度入度之说
# 有向图存在欧拉路径<=> 出度减入度不为 0 的点只能有 0 或 2 个，且 2 个时`起点出度比入度多 1，终点入度比出度多 1`

from typing import DefaultDict, List, Set, Tuple


def getEulerPath(
    n: int, adjMap: DefaultDict[int, Set[int]], *, isDirected: bool
) -> Tuple[bool, List[int]]:
    """求欧拉路径，需要寻找出发点"""
    start = next(iter(adjMap.keys()))

    if isDirected:
        indegree, outdegree = [0] * n, [0] * n
        minusOne, one = 0, 0
        for cur, nexts in adjMap.items():
            outdegree[cur] += len(nexts)
            for next_ in nexts:
                indegree[next_] += 1
        for cur in adjMap.keys():
            diff = outdegree[cur] - indegree[cur]
            if diff == 0:
                continue
            elif diff == 1:
                start = cur
                one += 1
            elif diff == -1:
                minusOne += 1
            else:
                return False, []
        if not (minusOne == one == 1) or not (minusOne == one == 0):
            return False, []
    else:
        oddCount = 0
        for cur in adjMap.keys():
            if len(adjMap[cur]) & 1:
                oddCount += 1
                start = cur
        if oddCount not in (0, 2):
            return False, []

    res = []
    stack = [start]
    cur = start
    while stack:
        if adjMap[cur]:
            stack.append(cur)
            next_ = adjMap[cur].pop()
            if not isDirected:
                adjMap[next_].remove(cur)  # 无向图 要删两条边
            cur = next_
        else:
            res.append(cur)
            cur = stack.pop()

    return True, res[::-1]

