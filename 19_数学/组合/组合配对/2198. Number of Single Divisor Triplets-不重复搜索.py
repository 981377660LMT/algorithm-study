# 3 <= nums.length <= 105
# 1 <= nums[i] <= 100

# 求出三元组(i,j,k)个数，
# 使得nums[i]+nums[j]+nums[k]只能被nums[i],nums[j],nums[k]中的一个整除
# single divisor triplet

# 直接搜答案


from typing import List
from collections import Counter
from itertools import product


class Solution:
    def singleDivisorTriplet(self, nums: List[int]) -> int:
        counter = Counter(nums)
        res = 0
        for n1, n2, n3 in product(counter.keys(), repeat=3):
            # 不重复搜素排列,非常关键!
            if not n1 <= n2 <= n3:
                continue
            sum_ = n1 + n2 + n3
            if sum(sum_ % num == 0 for num in (n1, n2, n3)) == 1:
                good1 = next((n for n in (n1, n2, n3) if sum_ % n == 0))
                bad1, bad2 = (n for n in (n1, n2, n3) if sum_ % n != 0)
                if bad1 == bad2:
                    res += counter[good1] * counter[bad1] * (counter[bad1] - 1) * 3
                else:
                    res += counter[good1] * counter[bad1] * counter[bad2] * 6
        return res


print(Solution().singleDivisorTriplet(nums=[4, 6, 7, 3, 2]))
print(Solution().singleDivisorTriplet(nums=[1, 2, 2]))
print(Solution().singleDivisorTriplet(nums=[1, 1, 1]))
