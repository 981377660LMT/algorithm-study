# 给定一个序列a，现在有两次操作，
# 第一次可以将前x个数字变成left,
# 第二次操作可以将后y个数字变成right,
# 询问在两次操作后该序列的最小值

# !正反dp一遍 再前后缀分解
# 移除违禁火车车厢
import sys
from typing import List

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

n, left, right = map(int, input().split())
nums = list(map(int, input().split()))


def getDp(nums: List[int], target: int) -> List[int]:
    dp = [0]
    for i, num in enumerate(nums, 1):
        cand1 = i * target
        cand2 = dp[-1] + num
        dp.append(min(cand1, cand2))
    return dp


dp1 = getDp(nums, left)
dp2 = getDp(nums[::-1], right)[::-1]
print(min(dp1[i] + dp2[i] for i in range(n + 1)))
