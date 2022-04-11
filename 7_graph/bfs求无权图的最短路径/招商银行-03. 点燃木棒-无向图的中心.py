from collections import defaultdict, deque
from typing import DefaultDict, Deque, List, Optional, Set, Tuple, Union, overload

MOD = int(1e9 + 7)
INF = int(1e20)

# 招商银行-03. 点燃木棒-无向图的中心  O(n^2)


@overload
def bfs(adjMap: DefaultDict[int, Set[int]], start: int) -> DefaultDict[int, int]:
    ...


@overload
def bfs(adjMap: DefaultDict[int, Set[int]], start: int, end: int) -> int:
    ...


def bfs(
    adjMap: DefaultDict[int, Set[int]], start: int, end: Optional[int] = None
) -> Union[int, DefaultDict[int, int]]:
    """时间复杂度O(V+E)"""
    dist = defaultdict(lambda: INF, {key: INF for key in adjMap.keys()})
    dist[start] = 0
    queue: Deque[Tuple[int, int]] = deque([(0, start)])
    while queue:
        curDist, cur = queue.popleft()
        if end is not None and cur == end:
            return curDist
        for next in adjMap[cur]:
            if dist[next] > dist[cur] + 1:
                dist[next] = dist[cur] + 1
                queue.append((dist[next], next))

    return INF if end is not None else dist


# 求无向图的中心点(到所有点的距离最大值最小) 注意图中可能存在环
class Solution:
    def lightSticks(self, height: int, width: int, indices: List[int]) -> List[int]:
        def getId(x: int, y: int) -> int:
            x, y = sorted([x, y])
            row, col = x // COL, x % COL
            # 横的
            if y - x == 1:
                return (width * 2 + 1) * row + col
            # 竖的
            else:
                return (width * 2 + 1) * row + col + width

        ROW, COL = height + 1, width + 1
        bad = set(indices)
        adjMap = defaultdict(set)

        for r in range(ROW):
            for c in range(COL):
                cur = r * COL + c
                if r + 1 < ROW:
                    next = (r + 1) * COL + c
                    curId = getId(cur, next)
                    if curId not in bad:
                        adjMap[cur].add(next)
                        adjMap[next].add(cur)
                if c + 1 < COL:
                    next = r * COL + c + 1
                    curId = getId(cur, next)
                    if curId not in bad:
                        adjMap[cur].add(next)
                        adjMap[next].add(cur)

        # print(getId(3, 6))
        # 求无向图的中心点(到所有点的距离最大值最小) 注意图中可能存在环
        # 注意到n<=2500 可以从每个点bfs
        max_ = INF
        res = []
        for start in adjMap.keys():
            dist = bfs(adjMap, start)
            curMax = max(dist.values())
            if curMax == INF:
                return []
            if curMax < max_:
                max_ = curMax
                res = [start]
            elif curMax == max_:
                res.append(start)

        return sorted(res)


print(Solution().lightSticks(height=2, width=2, indices=[2, 5, 6, 7, 8, 10, 11]))
print(Solution().lightSticks(height=1, width=2, indices=[3]))
print(Solution().lightSticks(height=1, width=1, indices=[0, 3]))
print(
    Solution().lightSticks(
        height=17, width=18, indices=[81, 376, 558, 470, 308, 395, 386, 315, 249, 499]
    )
)
# [161,180]
