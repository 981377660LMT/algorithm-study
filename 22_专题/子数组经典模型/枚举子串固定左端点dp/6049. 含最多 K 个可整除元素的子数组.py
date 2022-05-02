from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= nums.length <= 200


class Solution:
    def countDistinct(self, nums: List[int], k: int, p: int) -> int:
        """按照起点枚举子数组
        
        tuple + set 去重
        """
        res = set()
        for start in range(len(nums)):
            count = 0
            for end in range(start, len(nums)):
                if nums[end] % p == 0:
                    count += 1
                if count <= k:
                    res.add(tuple((nums[start : end + 1])))
                else:
                    break

        return len(res)


print(Solution().countDistinct(nums=[2, 3, 3, 2, 2], k=2, p=2))

# [2]、[2,3]、[2,3,3]、[2,3,3,2]、[3]、[3,3]、[3,3,2]、[3,3,2,2]、[3,2]、[3,2,2] 和 [2,2] 。
