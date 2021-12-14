from typing import List


class Solution:
    def checkArithmeticSubarrays(self, nums: List[int], l: List[int], r: List[int]) -> List[bool]:
        res = []

        def find_diffs(arr):

            arr.sort()

            dif = []

            for i in range(len(arr) - 1):
                dif.append(arr[i] - arr[i + 1])

            return len(set(dif)) == 1

        for i, j in zip(l, r):
            res.append(find_diffs(nums[i : j + 1]))

        return res


# 输入：nums = [4,6,5,9,3,7], l = [0,0,2], r = [2,3,5]
# 输出：[true,false,true]
# 解释：
# 第 0 个查询，对应子数组 [4,6,5] 。可以重新排列为等差数列 [6,5,4] 。
# 第 1 个查询，对应子数组 [4,6,5,9] 。无法重新排列形成等差数列。
# 第 2 个查询，对应子数组 [5,9,3,7] 。可以重新排列为等差数列 [3,5,7,9] 。

