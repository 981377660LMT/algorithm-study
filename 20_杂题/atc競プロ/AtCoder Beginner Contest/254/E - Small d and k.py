# 有n个顶点，m条无向边，顶点的编号为1-n (n<=1.5e5,m<=2.25e5)
# !并且保证所有点的度数不超过3。
# 有q个询问，每次询问给你一个顶点编号ver和最大深度k(k ≤3)，问你从ver点出发，离ver的深度不能超过k ,所能经过的所有的点的编号之和。
# 思路:
# 由于题目已经保证了询问的深度不超过3，并且图的每个点的度数不超过3.
# 考虑极限情况下，我们最多能遍历的个数不会超过40个。
from collections import defaultdict, deque
import sys
import os

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = int(1e9 + 7)


def main() -> None:
    def bfs(start: int, depthLimit: int) -> int:
        queue = deque([start])
        visited = set([start])
        res = 0
        depth = 0
        while queue and depth <= depthLimit:
            len_ = len(queue)
            for _ in range(len_):
                cur = queue.popleft()
                res += cur
                for next in adjMap[cur]:
                    if next not in visited:
                        visited.add(next)
                        queue.append(next)
            depth += 1

        return res

    n, m = map(int, input().split())
    adjMap = defaultdict(set)
    for _ in range(m):
        u, v = map(int, input().split())
        adjMap[u].add(v)
        adjMap[v].add(u)
    q = int(input())
    for _ in range(q):
        start, maxDepth = map(int, input().split())
        print(bfs(start, maxDepth))


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
