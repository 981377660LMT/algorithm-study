# - 最大两段子段和：求每个位置上的前缀最大子段和和后缀最大子段和
#   `等价于允许翻转一段子区间的最大子段和`
from typing import List


def solve(nums: List[int]) -> int:
    def getDp(nums: List[int]) -> List[int]:
        res = [0] * len(nums)
        curMax, allMax = -int(1e20), -int(1e20)
        for i, num in enumerate(nums):
            if curMax < 0:
                curMax = 0
            curMax += num
            allMax = max(allMax, curMax)  # 以当前元素结尾的最大子数组和
            res[i] = allMax
        return res

    pre = getDp(nums)
    suf = getDp(nums[::-1])[::-1]
    return max(pre[i] + suf[i + 1] for i in range(len(nums) - 1))


print(solve([-2, 1, -3, 4, -1, 2, 1, -5, 4]))

