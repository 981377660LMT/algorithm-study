# 有向图存在欧拉回路<=>所有点出度等于入度
# 无向图存在欧拉回路<=>所有度数为偶数


from typing import DefaultDict, List, Set, Tuple


def getEulerPath(
    n: int, adjMap: DefaultDict[int, Set[int]], *, isDirected: bool
) -> Tuple[bool, List[int]]:
    """求欧拉回路，出发点任意"""
    start = next(iter(adjMap.keys()))
    
    if isDirected:
        indegree, outdegree = [0] * n, [0] * n
        for cur, nexts in adjMap.items():
            outdegree[cur] += len(nexts)
            for next_ in nexts:
                indegree[next_] += 1
        if any(indegree[cur] != outdegree[cur] for cur in adjMap.keys()):
            return False, []
    else:
        if any(len(adjMap[cur]) & 1 for cur in adjMap.keys()):
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
