# 检查是否存在负权环的方法为：记录一个点的入队次数，如果超过V-1次说明存在负权环，
# 因为最短路径上除自身外至多V-1个点，故一个点不可能被更新超过V-1次


from collections import defaultdict, deque
from typing import DefaultDict, Hashable, Set, Tuple, TypeVar, Union


T = TypeVar("T", bound=Hashable)


def spfa(
    adjMap: DefaultDict[T, DefaultDict[T, int]], allVertex: Set[T]
) -> Union[Tuple[DefaultDict[T, int], bool], Tuple[None, bool]]:
    """spfa求单源最长路并断正环"""
    dist = defaultdict(int)
    queue = deque(allVertex)
    inQueue = defaultdict(lambda: True)
    count = defaultdict(int)  # 边数更新次数

    while queue:
        cur = queue.popleft()
        inQueue[cur] = False

        for next, weight in adjMap[cur].items():
            if dist[cur] + weight < dist[next]:
                count[next] = count[cur] + 1
                if count[next] >= len(allVertex):
                    return None, True
                dist[next] = dist[cur] + weight
                if not inQueue[next]:
                    inQueue[next] = True
                    if queue and dist[next] < dist[queue[0]]:  # !酸辣粉优化
                        queue.appendleft(next)
                    else:
                        queue.append(next)

    return dist, False
