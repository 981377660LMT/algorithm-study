# 起点在(0,0），终点在(x,y)，
# 每步能走ai，但下一步和上一步要成90度角，
# 第一步是确定了只能走到(a1,0)，
# 在x或y方向上每一步走a[i]或-a[i]，问最终能否达到(x,y).
# !观察发现（下标从1开始）下标为奇数的步是走在x方向的，偶数步是走在y方向上的，
# 因此行和列可以分开dp
# !n<=1e3 ai<=10 -1e4x,y<=1e4

from typing import List


def cal(start: int, end: int, moves: List[int]) -> bool:
    """从start出发,每步可以走moves[i]或-moves[i],问能否到达end"""
    dp = set([start])
    for move in moves:
        ndp = set()
        for pre in dp:
            ndp.add(pre + move)
            ndp.add(pre - move)
        dp = ndp
    return end in dp


n, endX, endY = map(int, input().split())
nums = list(map(int, input().split()))
startX, startY = nums[0], 0
nums = nums[1:]
yMoves = nums[::2]
xMoves = nums[1::2]
print("Yes" if cal(startX, endX, xMoves) and cal(startY, endY, yMoves) else "No")
