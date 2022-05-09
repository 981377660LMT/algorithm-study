# 无向图存在欧拉回路<=>所有度数为偶数
# 有向图存在欧拉回路<=>所有点出度等于入度


from typing import DefaultDict, List, Set, Tuple


def getEulerPath(
    allVertex: Set[int], adjMap: DefaultDict[int, Set[int]], *, isDirected: bool
) -> Tuple[bool, List[int]]:
    """求欧拉回路，出发点任意，保证输入的图是连通图"""

    if isDirected:
        indegree, outdegree = {v: 0 for v in allVertex}, {v: 0 for v in allVertex}
        for cur, nexts in adjMap.items():
            outdegree[cur] += len(nexts)
            for next_ in nexts:
                indegree[next_] += 1

        for cur in allVertex:
            diff = outdegree[cur] - indegree[cur]
            if diff != 0 or outdegree[cur] == 0:  # 所有点出度等于入度/入度为 0，出度也为 0，不是联通图
                return False, []

    else:
        for cur in allVertex:
            degree = len(adjMap[cur])
            if (degree & 1) or degree == 0:  # 所有点度数为偶数/度数为 0，不是联通图
                return False, []

    start = next(iter(adjMap.keys()))
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
