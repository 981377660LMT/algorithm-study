from typing import List
from collections import Counter

# lower[i] = arr[i] - k
# higher[i] = arr[i] + k

# 已知lower与higher的混合，求出可能的arr
# 直接枚举最大的那个数对应的lower中的数即可。  (枚举)

# O(n^2logn)  枚举 O(n^2)+排序O(nlogn)
# 1 <= n <= 1000
# 1 <= nums[i] <= 10^9


class Solution:
    def recoverArray(self, nums: List[int]) -> List[int]:
        n = len(nums)
        nums = sorted(nums)
        keys = sorted(set(nums))

        # 枚举间隔k
        for i in range(n - 2, -1, -1):
            diff = nums[-1] - nums[i]
            if diff == 0 or diff & 1:
                continue

            k = diff // 2
            counter = Counter(nums)
            res = []
            for key in keys:
                if counter[key + 2 * k] < counter[key]:
                    break
                res.extend([key + k] * counter[key])
                counter[key + 2 * k] -= counter[key]
            else:
                return res

        return []


# print(Solution().recoverArray([2, 10, 6, 4, 8, 12]))
# print(Solution().recoverArray([1, 1, 3, 3]))
# print(Solution().recoverArray([11, 6, 3, 4, 8, 7, 8, 7, 9, 8, 9, 10, 10, 2, 1, 9]))
# # [2,3,4,5,7,8,8,9]
# # 预期：[2,3,7,8,8,9,9,10]
# print(Solution().recoverArray([8, 4, 5, 1, 9, 8, 6, 5, 6, 9, 7, 3, 8, 3, 6, 7, 10, 11, 6, 4]))
# [1, 3, 3, 4, 4, 5, 5, 6, 6, 6, 6, 7, 7, 8, 8, 8, 9, 9, 10, 11]
# [2,4,5,5,6,7,7,8,9,10]
print(Solution().recoverArray([8, 4, 5, 1, 9, 8, 6, 5, 6, 9, 7, 3, 8, 3, 6, 7, 10, 11, 6, 4]))

# 养成好习惯：数组必定先求length
