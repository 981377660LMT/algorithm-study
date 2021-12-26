from typing import List
from collections import Counter

# 1 <= changed.length <= 10^5
# 954. 二倍数对数组.py
class Solution:
    def findOriginalArray(self, changed: List[int]) -> List[int]:
        n = len(changed)
        if n & 1:
            return []

        counter = Counter(changed)
        res = []
        for key in sorted(counter):
            # 处理0的特殊情况
            if key == 0:
                if counter[key] & 1:
                    return []
                res.extend([key] * (counter[key] // 2))
                continue

            # 配对相减
            elif counter[key] > counter[key * 2]:
                return []
            res.extend([key] * counter[key])
            counter[key * 2] -= counter[key]

        return res


print(Solution().findOriginalArray(changed=[1, 3, 4, 2, 6, 8]))
# 输出：[1,3,4]
# 解释：一个可能的 original 数组为 [1,3,4] :
# - 将 1 乘以 2 ，得到 1 * 2 = 2 。
# - 将 3 乘以 2 ，得到 3 * 2 = 6 。
# - 将 4 乘以 2 ，得到 4 * 2 = 8 。
# 其他可能的原数组方案为 [4,3,1] 或者 [3,1,4] 。
print(Solution().findOriginalArray(changed=[0, 0, 0, 0]))

