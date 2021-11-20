# 设 X 代表小 A 和小 B 的浏览顺序中出现在同一位置的简历数，求 X 的期望。

from typing import List
from collections import Counter


class Solution:
    def expectNumber(self, scores: List[int]) -> int:
        return len(set(scores))


# 输入：scores = [1,1]

# 输出：1

# 解释：设两位面试者的编号为 0, 1。由于他们的能力值都是 1，
# 小 A 和小 B 的浏览顺序都为从全排列 [[0,1],[1,0]] 中等可能地取一个。
# 如果小 A 和小 B 的浏览顺序都是 [0,1] 或者 [1,0] ，那么出现在同一位置的简历数为 2 ，
# 否则是 0 。所以 X 的期望是 (2+0+2+0) * 1/4 = 1

# 总结:对同一个数 他们的贡献之和为1
