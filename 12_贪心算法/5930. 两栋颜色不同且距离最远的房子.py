from typing import List


class Solution:
    def maxDistance1(self, colors: List[int]) -> int:
        return max(
            abs(i - j)
            for i in range(len(colors))
            for j in range(i + 1, len(colors))
            if colors[i] != colors[j]
        )

    # 只需找到第一个和首尾不同的房子
    def maxDistance(self, colors: List[int]) -> int:
        n = len(colors)
        cand1 = next((i for i in range(n) if colors[i] != colors[-1]), 0)
        cand2 = next((i for i in range(n - 1, -1, -1) if colors[i] != colors[0]), n - 1)
        return max(n - 1 - cand1, cand2)


print(Solution().maxDistance([1, 1, 1, 6, 1, 1, 1]))

# 思路：贪心，从右边招和head不一样的 从左边找和tail不一样的
