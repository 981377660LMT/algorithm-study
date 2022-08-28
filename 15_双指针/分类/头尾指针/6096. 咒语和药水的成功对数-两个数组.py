from typing import List, Tuple


MOD = int(1e9 + 7)
INF = int(1e20)

# 一个咒语和药水的能量强度 相乘 如果 大于等于 success ，那么它们视为一对 成功 的组合。
# 请你返回一个长度为 n 的整数数组 pairs，其中 pairs[i] 是能跟第 i 个咒语成功组合的 药水 数目。

# !如果不能选一样的 怎么办
# !哈希表存储每种强度 减去comb(num,2)
# 15_双指针/分类/头尾指针/法力值大于k的配对数.py

# 头尾双指针
class Solution:
    def successfulPairs(self, spells: List[int], potions: List[int], success: int) -> List[int]:
        queries = sorted([(num, i) for i, num in enumerate(spells)], key=lambda x: x[0])
        potions.sort()

        right = len(potions) - 1
        res = [0] * len(queries)
        for qv, qi in queries:
            while right >= 0 and potions[right] * qv >= success:
                right -= 1
            res[qi] = len(potions) - 1 - right

        return res


print(Solution().successfulPairs(spells=[5, 1, 3], potions=[1, 2, 3, 4, 5], success=7))
