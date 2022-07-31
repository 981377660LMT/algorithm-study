# bfs模板 类似于dijkstra


if __name__ == "__main__":

    # from typing import Hashable, TypeVar
    # T = TypeVar("T", bound=Hashable)

    from collections import defaultdict, deque

    INF = int(1e20)

    def bfs(adjMap: defaultdict[int, set[int]], start: int) -> defaultdict[int, int]:
        """时间复杂度O(V+E)"""
        dist = defaultdict(lambda: INF, {start: 0})
        queue: deque[tuple[int, int]] = deque([(0, start)])

        while queue:
            _, cur = queue.popleft()
            for next in adjMap[cur]:
                if dist[next] > dist[cur] + 1:
                    dist[next] = dist[cur] + 1
                    queue.append((dist[next], next))

        return dist
