from typing import List, Tuple

# n<=18
# m<=2**n-1
# !是否存在一棵生成树 满足任意边的顶点异或值不存在于badXor中
# 线性基、格雷码


def encounterAndFarewell(n: int, m: int, badXor: List[int]) -> Tuple[List[Tuple[int, int]], bool]:
    """出会いと別れ  韻を踏んでなくない???"""
    ...


if __name__ == "__main__":
    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")

    n, m = map(int, input().split())
    badXor = list(map(int, input().split()))
    res, ok = encounterAndFarewell(n, m, badXor)
    if not ok:
        print(-1)
        exit(0)
    print(*res, sep="\n")
