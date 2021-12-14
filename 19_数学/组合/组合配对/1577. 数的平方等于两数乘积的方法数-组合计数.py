from typing import List
from collections import Counter

# 给你两个整数数组 nums1 和 nums2 ，请你返回根据以下规则形成的三元组的数目

# 1 <= nums1.length, nums2.length <= 1000

# 哈希表两数之和思想:固定目标，枚举小的数
class Solution:
    def numTriplets(self, nums1: List[int], nums2: List[int]) -> int:
        def countTriplets(curCounter: Counter, targetCounter: Counter) -> int:
            res = 0
            for t1, c1 in targetCounter.items():
                target = t1 ** 2
                for t2, c2 in curCounter.items():
                    if target % t2 != 0:
                        continue
                    t3 = target // t2

                    if t2 == t3:
                        res += c1 * c2 * (c2 - 1) // 2

                    # 注意这里防止重复计算的细节
                    # 1.两种情况只算一次
                    elif t2 < t3 and t3 in curCounter:
                        c3 = curCounter[t3]
                        res += c1 * c2 * c3
            return res

        c1, c2 = Counter(nums1), Counter(nums2)
        return countTriplets(c1, c2) + countTriplets(c2, c1)


print(Solution().numTriplets(nums1=[1, 1], nums2=[1, 1, 1]))
# 输出：9
# 解释：所有三元组都符合题目要求，因为 1^2 = 1 * 1
# 类型 1：(0,0,1), (0,0,2), (0,1,2), (1,0,1), (1,0,2), (1,1,2), nums1[i]^2 = nums2[j] * nums2[k]
# 类型 2：(0,0,1), (1,0,1), (2,0,1), nums2[i]^2 = nums1[j] * nums1[k]
