from bisect import bisect_left
from typing import List
from collections import defaultdict


MOD = int(1e9 + 7)
INF = int(1e20)


# !总结:
# !1. (a & b) << 1 + (a ^ b) = a + b
# !2. (a | b).bit_count() + (a & b).bit_count() = a.bit_count() + b.bit_count()
# !3. 前缀和/二分:固定某个数，求 sum(某个数>=k) 的个数

# 脑筋急转弯 + 两数之和


class Solution:
    def countExcellentPairs(self, nums: List[int], k: int) -> int:
        visited = set(nums)
        counter = defaultdict(int)
        for num in visited:
            counter[num.bit_count()] += 1

        # 前缀和查找
        preSum = [0]
        for i in range(70):
            preSum.append(preSum[-1] + counter[i])

        res = 0
        for cur in range(70):
            count1 = counter[cur]
            count2 = preSum[-1] - preSum[max(0, k - cur)]
            res += count1 * count2
        return res

    def countExcellentPairs2(self, nums: List[int], k: int) -> int:
        nums = list(set(nums))
        counts = []
        for num in nums:
            counts.append(num.bit_count())
        counts.sort()

        # 二分查找也可以
        res = 0
        for count in counts:
            pos = bisect_left(counts, k - count)
            res += len(nums) - pos
        return res


# print(Solution().countExcellentPairs(nums=[1, 2, 3, 1], k=3))
# print(Solution().countExcellentPairs(nums=[5, 1, 1], k=10))
# print(Solution().countExcellentPairs(nums=[1, 2, 3, 1, 536870911], k=3))
print(
    Solution().countExcellentPairs(
        nums=[
            423436147,
            29690092,
            724837828,
            339900252,
            819138876,
            559797269,
            337577818,
            347372617,
            568172510,
            434046210,
            233465903,
            73777015,
            995100887,
            952551841,
            314703814,
            588503612,
            5824363,
            105686599,
            5167368,
            154358365,
            497653021,
            450975800,
            431388582,
            607991479,
            856148544,
            982787927,
            513430676,
            918344731,
            98092726,
            690894469,
            396191705,
            848402861,
            593468334,
            563155911,
            715586102,
            739434236,
            387304407,
            927581316,
            779272764,
            558853665,
            215920106,
            631709145,
            726054493,
            415057810,
            708860839,
            401596916,
            795418594,
            462963799,
            835785708,
            670198432,
            171214014,
            162179684,
            27901422,
            717744871,
            603604788,
            664478320,
            915525044,
            818068724,
            564705733,
            490294265,
            804021123,
            688892990,
            741612165,
            590640255,
            535167444,
            228105610,
            197887678,
            963803394,
            698521654,
            794863135,
            712203903,
            16780599,
            583378338,
            927863644,
            628601885,
            878322079,
            632547981,
            426926648,
        ],
        k=55,
    )
)
# 预期12
