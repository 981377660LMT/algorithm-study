from typing import List


class Solution:
    def __maxDistance(self, colors: List[int]) -> int:
        return max(
            abs(i - j)
            for i in range(len(colors))
            for j in range(i + 1, len(colors))
            if colors[i] != colors[j]
        )

    def maxDistance(self, colors: List[int]) -> int:
        if colors[0] != colors[-1]:
            return len(colors) - 1
        first, last = colors[0], colors[-1]
        res1, res2 = 0, 0
        for i in range(1, len(colors)):
            if colors[i] != last:
                res1 = i
                break
        for i in range(len(colors) - 1, 0, -1):
            if colors[i] != first:
                res2 = i
                break
        return max(len(colors) - res1 - 1, res2)


print(Solution().maxDistance([1, 1, 1, 6, 1, 1, 1]))

# 思路：贪心，从右边招和head不一样的 从左边找和tail不一样的
