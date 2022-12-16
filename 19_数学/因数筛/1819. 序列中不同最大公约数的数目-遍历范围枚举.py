from typing import List
from collections import defaultdict

# 子序列的 GCD 的不同个数
# 题目提示:
# 1 <= nums.length <= 1e5
# !1 <= nums[i] <= 2 * 1e5

# enumerate all possibilities
# https://leetcode-cn.com/problems/number-of-different-subsequences-gcds/solution/pythonjian-dan-si-lu-he-dai-ma-by-semiro-urs1/
# 例如：[2, 4, 6], 对从1开始的所有的数遍历，取出所有他们的倍数存入列表
# 1 -> [2, 4, 6]
# 2 -> [2, 4, 6]
# 3 -> [6]
# 4 -> [4]
# 5 -> []
# 6 -> [6]

# !考虑每一个数是否可以成为某个子序列的最大公因数即可
# 然后把所有的列表变成tuple用set去重。得到的长度就是答案。


# 筛法的复杂度为n/1 + n/2 + n/3 + … + n/n 渐进为O(n * logn)
class Solution:
    def countDifferentSubsequenceGCDs(self, nums: List[int]) -> int:
        """计算并返回 nums 的所有 非空 子序列中 不同 最大公约数的 数目 。"""
        numSet = set(nums)
        max_ = max(numSet)
        adjMap = defaultdict(list)

        # 枚举因子
        for factor in range(1, max_ + 1):
            for multi in range(factor, max_ + 1, factor):
                if multi in numSet:
                    adjMap[factor].append(multi)

        return len(set(tuple(sub) for sub in adjMap.values()))


print(Solution().countDifferentSubsequenceGCDs(nums=[6, 10, 3]))
# 输出：5
# 解释：上图显示了所有的非空子序列与各自的最大公约数。
# 不同的最大公约数为 6 、10 、3 、2 和 1 。

print(Solution().countDifferentSubsequenceGCDs(nums=[3, 6]))
