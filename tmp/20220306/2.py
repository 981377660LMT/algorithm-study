from bisect import bisect_left, bisect_right
from typing import List, Tuple


MOD = int(1e9 + 7)

# 有重复元素
# JS 的 number 上限比 int64 要小。这是个坑。


class Solution:
    def minimalKSum(self, nums: List[int], k: int) -> int:
        def findMex(arr: List[int], k: int) -> int:
            """二分搜索缺失的第k个正整数,lc1539. 第 k 个缺失的正整数"""
            # MEX:Min Excluded
            left, right = 0, len(arr) - 1
            while left <= right:
                mid = (left + right) >> 1
                diff = arr[mid] - (mid + 1)
                if diff >= k:
                    right = mid - 1
                else:
                    left = mid + 1
            return left + k

        nums = sorted(set(nums))
        mex = findMex(nums, k)
        index = bisect_right(nums, mex)
        allsum = (mex + 1) * (mex) // 2
        print(index, allsum)
        return allsum - sum(nums[:index])


# 500000000500000001
print(Solution().minimalKSum(nums=[1000000000], k=1000000000))

print(Solution().minimalKSum(nums=[1, 4, 25, 10, 25], k=2))
print(
    Solution().minimalKSum(
        nums=[
            96,
            44,
            99,
            25,
            61,
            84,
            88,
            18,
            19,
            33,
            60,
            86,
            52,
            19,
            32,
            47,
            35,
            50,
            94,
            17,
            29,
            98,
            22,
            21,
            72,
            100,
            40,
            84,
        ],
        k=35,
    )
)
# 824
# 预期：
# 794
