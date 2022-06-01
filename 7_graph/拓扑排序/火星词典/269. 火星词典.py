from collections import defaultdict, deque
from itertools import pairwise
from typing import DefaultDict, Hashable, List, Set, Tuple, TypeVar

T = TypeVar('T', bound=Hashable)


def toposort2(
    adjMap: DefaultDict[T, Set[T]], deg: DefaultDict[T, int], /, allVertex: Set[T]
) -> Tuple[int, List[T]]:
    """"返回有向图拓扑排序方案数和拓扑排序结果"""
    for v in allVertex:  # !初始化所有顶点的入度
        deg.setdefault(v, 0)
    queue = deque([v for v in allVertex if deg[v] == 0])
    res, topoCount = [], 1
    while queue:
        topoCount *= len(queue)
        cur = queue.popleft()
        res.append(cur)
        for next in adjMap[cur]:
            deg[next] -= 1
            if deg[next] == 0:
                queue.append(next)
    if len(res) != len(allVertex):
        return 0, []
    return topoCount, res


class Solution:
    def alienOrder(self, words: List[str]) -> str:
        adjMap = defaultdict(set)
        allVertex = set(char for word in words for char in word)
        deg = defaultdict(int, {v: 0 for v in allVertex})
        for pre, cur in pairwise(words):
            allVertex.update(pre, cur)  # update加入多个iterable
            for char1, char2 in zip(pre, cur):
                if char1 != char2:
                    if char2 not in adjMap[char1]:  # !如果adjMap用set 需要保证deg不重复计算(即重边只算一次)
                        adjMap[char1].add(char2)
                        deg[char2] += 1
                    break
            else:
                if len(pre) > len(cur):
                    return ''

        _, res = toposort2(adjMap, deg, allVertex)
        return ''.join(res)


# print(Solution().alienOrder(["wrt", "wrf", "er", "ett", "rftt"]))
# print(Solution().alienOrder(["ab", "adc"]))
print(Solution().alienOrder(["ac", "ab", "zc", "zb"]))

