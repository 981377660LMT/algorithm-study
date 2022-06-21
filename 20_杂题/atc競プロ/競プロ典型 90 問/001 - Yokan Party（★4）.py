import sys


# !二分
input = sys.stdin.readline

# 左右の長さが L [cm] のようかんがあります。
# あなたは N 個の切れ目のうち K 個を選び、ようかんを K+1 個のピースに分割したいです
# K+1 個のピースのうち、最も短いものの長さ  を スコア とします
# スコアが最大となるように分割する場合に得られるスコアを求めてください。

N, L = map(int, input().split())
K = int(input())
splits = list(map(int, input().split()))
splits = [0] + splits + [L]
nums = [cur - pre for pre, cur in zip(splits, splits[1:])]


def check(mid: int) -> bool:
    """选择k+1块 每块长度>=mid"""
    res = 0
    curSum = 0
    for num in nums:
        curSum += num
        if curSum >= mid:
            res += 1
            curSum = 0
    return res >= K + 1


left, right = 0, L
while left <= right:
    mid = (left + right) // 2
    if check(mid):
        left = mid + 1
    else:
        right = mid - 1
print(right)
