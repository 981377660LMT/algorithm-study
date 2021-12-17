from typing import List
from collections import defaultdict

# 题目提示:
# 1 <= nums.length <= 105
# 1 <= nums[i] <= 2 * 105

# enumerate all possibilities
# 直接枚举因子,看他是数组里哪些数的因子,如果tuple相同,取大的那个
# https://leetcode-cn.com/problems/number-of-different-subsequences-gcds/solution/pythonjian-dan-si-lu-he-dai-ma-by-semiro-urs1/
# 例如：[2, 4, 6], 对从1开始的所有的数遍历，取出所有他们的倍数存入列表
# 1 -> [2, 4, 6]
# 2 -> [2, 4, 6]
# 3 -> [6]
# 4 -> [4]
# 5 -> []
# 6 -> [6]

# 然后把所有的列表变成tuple用set去重。得到的长度就是答案(3)。


# 筛法的复杂度为n/1 + n/2 + n/3 + … + n/n 渐进为O(n * logn)
class Solution:
    def countDifferentSubsequenceGCDs(self, nums: List[int]) -> int:
        numSet = set(nums)
        maxVal = max(numSet)
        commonFactorOfSub = defaultdict(list)

        # 类似于素数筛
        for fac in range(1, maxVal + 1):
            for multi in range(fac, maxVal + 1, fac):
                if multi in numSet:
                    commonFactorOfSub[fac].append(multi)

        return len(set([tuple(sub) for sub in commonFactorOfSub.values()]))


print(Solution().countDifferentSubsequenceGCDs(nums=[6, 10, 3]))
# 输出：5
# 解释：上图显示了所有的非空子序列与各自的最大公约数。
# 不同的最大公约数为 6 、10 、3 、2 和 1 。

