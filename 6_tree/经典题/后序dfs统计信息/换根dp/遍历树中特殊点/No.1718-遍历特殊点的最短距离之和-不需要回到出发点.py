# https://yukicoder.me/problems/no/1718
# 从每个结点出发,求遍历所有特殊点的最短距离之和
# n<=1e5
# TODO

from typing import List, Tuple
from Rerooting import Rerooting


def randomSquirrel(n: int, edges: List[Tuple[int, int]], sweets: List[int]) -> List[int]:
    ...


if __name__ == "__main__":
    n, k = map(int, input().split())
    edges = []
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        edges.append((u, v))
    sweets = list(map(int, input().split()))
    print(*randomSquirrel(n, edges, sweets), sep="\n")
