'''
给你一个整数数组 nums​​​ 和一个整数 k 。你需要将这个数组划分到 k 个相同大小的子集中，
使得同一个子集里面没有两个相同的元素。
请你返回将数组分成 k 个子集后，各子集 不兼容性(是该子集里面最大值和最小值的差) 的 和 的 最小值 ，如果无法分成分成 k 个子集，返回 -1 。
'''

from typing import List
from collections import Counter


class Solution:
    def minimumIncompatibility(self, nums: List[int], k: int) -> float:
        if (max(Counter(nums).values())) > k:
            return -1
        res = float('inf')
        nums.sort(reverse=True)
        arr = [[] for _ in range(k)]
        upper = len(nums) // k

        def bt(index: int):
            if index == len(nums):
                nonlocal res
                res = min(res, sum(arr[i][0] - arr[i][-1] for i in range(k)))
                return True
            flag = 0
            for j in range(k):
                if not arr[j] or len(arr[j]) < upper and arr[j][-1] != nums[index]:
                    arr[j].append(nums[index])
                    if bt(index + 1):
                        flag += 1
                    arr[j].pop()
                    # nums[i] can be assigned to arr[j] and arr[j+1]
                if flag >= 2:
                    break
            return flag != 0

        bt(0)
        return res


print(Solution().minimumIncompatibility([6, 3, 8, 1, 3, 1, 2, 2], 4))
