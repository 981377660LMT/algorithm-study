from typing import List, Tuple


def sccTourByIndegs(indeg: List[int]) -> Tuple[int, List[int]]:
    """给定竞赛图(Tournament, 有向完全图)的每个点的入度, 返回SCC的个数和每个点所属的SCC的编号"""
    n = len(indeg)
    order = sorted(range(n), key=lambda x: indeg[x])
    sccId = [0] * n
    degSum = 0
    nextId = 0
    for i in range(n):
        cur = order[i]
        sccId[cur] = nextId
        degSum += indeg[cur]
        diff = degSum - (i + 1) * i // 2  # 前i个点构成了一个强连通分量
        if diff == 0:
            nextId += 1
    return nextId, sccId
