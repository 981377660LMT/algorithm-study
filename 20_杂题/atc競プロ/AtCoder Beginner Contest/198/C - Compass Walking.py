# 圆规行走
# 从原点出发 每次必须走到距离当前点为r的点(不必走到网格上的整点)
# 问最少走多少步能走到tr,tc
# r<=1e5 x,y<=1e5
# !注意不能bfs 会TLE 需要数学解法

# !计算(x,y)到原点距离后向上取整，特别地，对于距离小于r的，需要两步才能到。


from math import ceil


def compassWalk(r: int, tr: int, tc: int) -> int:
    dist = (tr**2 + tc**2) ** 0.5
    if dist < r:
        return 2  # 折返
    return ceil(dist / r)


if __name__ == "__main__":

    import sys

    sys.setrecursionlimit(int(1e9))
    input = lambda: sys.stdin.readline().rstrip("\r\n")
    R, tr, tc = map(int, input().split())
    print(compassWalk(R, tr, tc))
