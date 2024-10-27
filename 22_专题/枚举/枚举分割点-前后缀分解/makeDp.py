# 前后缀分解模板

from typing import List


def solve(nums: List[int]) -> int:
    def makeDp(arr: List[int]) -> List[int]:
        n = len(arr)
        dp = [0] * (n + 1)
        for i in range(1, n + 1):
            cur = arr[i - 1]
            #
        return dp

    preDp, sufDp = makeDp(nums), makeDp(nums[::-1])[::-1]
    res = 0
    for i in range(len(nums) + 1):
        res += preDp[i] * sufDp[i]  # [0,i) [i,n)
