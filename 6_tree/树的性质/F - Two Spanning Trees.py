"""
无向图G中求两个生成树T1和T2
当1为根时
要求G中不在T1里任意两条边都有祖孙关系
要求G中不在T2里任意两条边都没有祖孙关系
输出这两棵生成树的边
# n,m<=2e5

利用dfs与bfs的性质

"""


from collections import defaultdict, deque
import sys
import os
from typing import Set

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    def dfs(cur: int, visited: Set[int]) -> None:
        """dfs生成树"""
        for next in adjMap[cur]:
            if next not in visited:
                visited.add(next)
                print(f"{cur} {next}")
                dfs(next, visited)

    def bfs(start: int, visited: Set[int]) -> None:
        """bfs生成树"""
        queue = deque([start])
        count = 0
        while queue:
            cur = queue.popleft()
            for next in adjMap[cur]:
                if next not in visited:
                    visited.add(next)
                    print(f"{cur} {next}")
                    count += 1
                    if count == n - 1:
                        break
                    queue.append(next)

    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        adjMap[u].add(v)
        adjMap[v].add(u)

    dfs(1, set([1]))
    bfs(1, set([1]))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
