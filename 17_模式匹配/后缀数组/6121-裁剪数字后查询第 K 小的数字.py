#  1. 将 nums 中每个数字 裁剪 到剩下 最右边 trimi 个数位。

#  2. 在裁剪过后的数字中，找到 nums 中第 ki 小数字对应的 下标 。
#  如果两个裁剪后数字一样大，那么下标 更小 的数字视为更小的数字。

#  3. 裁剪后的数字中，找到 nums 中第 ki 小数字对应的 下标 。

#  !把所有数字字符串拼接后，用后缀数组计算出rank，再扫一遍就能算出每个trim对应的相对位置。
#  !由于数字相同时，序号小的排前面，那么就用序号作为每个字符串的结尾，就能得到题目要求的顺序了

from collections import defaultdict
from typing import List
from SuffixArray import sa_is


class Solution:
    def smallestTrimmedNumbers(self, nums: List[str], queries: List[List[int]]) -> List[int]:
        n, ords = len(nums[0]), []
        for i, word in enumerate(nums):
            ords.extend(list(map(ord, word)))
            ords.append(i)  # 相当于每个字符串的长度是n+1

        sa = sa_is(ords, max(ords))
        indexMap = defaultdict(list)  # 每个单词后缀长度 => 从小到大排列的单词序号
        for saIndex in sa:
            wordIndex, suffixLen = divmod(saIndex, n + 1)
            indexMap[suffixLen].append(wordIndex)

        return [indexMap[n - trim][k - 1] for k, trim in queries]


print(
    Solution().smallestTrimmedNumbers(
        nums=["102", "473", "251", "814"], queries=[[1, 1], [2, 3], [4, 2], [1, 2]]
    )
)
