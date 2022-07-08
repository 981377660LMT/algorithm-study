from collections import defaultdict
import sys
import os
from typing import DefaultDict, List, Set

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353


def solve(adjMap: DefaultDict[int, Set[int]], values: List[int]) -> None:
    """
    树上两点距离定义为`路径上点的个数*路径上点的值的gcd`
    求树中所有pair的距离之和模MOD (0<=i<j<n)

    「全ての頂点対についての距離の和」なので、「各辺が和に何度寄与するか」を考える
    计算边的贡献
    """
    ...


def main() -> None:
    n = int(input())
    values = list(map(int, input().split()))
    adjMap = defaultdict(set)
    for _ in range(n - 1):
        u, v = map(int, input().split())
        u, v = u - 1, v - 1
        adjMap[u].add(v)
        adjMap[v].add(u)
    solve(adjMap, values)


if __name__ == "__main__":
    if os.environ.get("USERNAME", " ") == "caomeinaixi":
        while True:
            main()
    else:
        main()
