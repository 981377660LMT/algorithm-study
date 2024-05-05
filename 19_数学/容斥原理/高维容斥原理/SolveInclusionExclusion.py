from typing import List


MOD = int(1e9 + 7)


def solveInclusionExclusion(nums: List[int]) -> int:
    n = len(nums)
    res = 0
    for state in range(1 << n):
        cur = 0
        count = 0
        for i, v in enumerate(nums):
            if state >> i & 1 > 0:
                # 视情况而定，有时候包含元素 i 表示考虑这种情况，有时候表示不考虑这种情况
                _ = v
                count += 1
        if count & 1 > 0:
            cur = -cur  # 某些题目是 == 0
        res = (res + cur) % MOD

    return res % MOD
