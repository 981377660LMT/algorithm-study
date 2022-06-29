# n个开关 m个灯泡
# 每个开关会使若干灯泡反转状态
# 每个开关只能按一次 灯泡变为全关 求按开关的方案数
# n,m<=300

# !行列の掃き出し法 (高斯消元)
# 消元后 每一列都只会剩下一个1 即一个开关控制一个维度(灯泡) 从上往下直接消去即可
import sys
from typing import List

sys.setrecursionlimit(int(1e9))
input = sys.stdin.readline
MOD = 998244353


def findBase(matrix: List[List[int]]) -> List[int]:
    """求线性基 输入为01矩阵"""
    ROW, COL = len(matrix), len(matrix[0])
    nums = [int("".join(map(str, row)), 2) for row in matrix]

    rowCount = 0  # 非0行计数
    for col in range(COL - 1, -1, -1):
        # 高斯消元，找到第一个1，把下面消成0
        for row in range(rowCount, ROW):
            if (nums[row] >> col) & 1:
                nums[row], nums[rowCount] = nums[rowCount], nums[row]
                break

        # 说明当前这一位所有向量都是0，这一位没有携带任何信息，直接跳过这一位
        if (nums[rowCount] >> col) & 1 == 0:
            continue

        for row in range(ROW):
            if row != rowCount and (nums[row] >> col) & 1:
                nums[row] ^= nums[rowCount]

        # 新组出来一个非0行向量，非0向量的个数是基的大小，不可能超过总的向量数
        rowCount += 1
        if rowCount == ROW:
            break

    return nums[:rowCount]


n, m = map(int, input().split())
matrix = [[0] * m for _ in range(n)]  # 每个开关控制哪些灯泡
for r in range(n):
    input()
    bulbs = tuple(int(x) - 1 for x in input().split())
    for c in bulbs:
        matrix[r][c] = 1


guass = findBase(matrix)  # 线性基(高斯消元后的矩阵)

initBulbs = list(map(int, input().split()))
state = int(''.join(map(str, initBulbs)), 2)
for base in guass:
    state = min(state, base ^ state)  # 反转首位

if state != 0:
    print(0)
else:
    print(pow(2, n - len(guass), MOD))  # 剩下高斯矩阵的0行 即可以任意开关

