from typing import List
from collections import Counter

# 奇数位等于偶数位乘以二
# 1-2 2-4 4-6 -8--16
class Solution:
    def canReorderDoubled(self, arr: List[int]) -> bool:
        counter = Counter(arr)
        for key in sorted(counter, key=abs):
            # 映射
            if counter[key] > counter[key * 2]:
                return False
            counter[key * 2] -= counter[key]
        return True


print(Solution().canReorderDoubled([4, -2, 2, -4]))
# 输出：true
# 解释：可以用 [-2,-4] 和 [2,4] 这两组组成 [-2,-4,2,4] 或是 [2,4,-2,-4]
