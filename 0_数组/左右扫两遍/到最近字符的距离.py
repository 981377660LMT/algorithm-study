from typing import List

# 蜡烛之间的盘子


class Solution:
    def solve(self, s: str, c: str) -> List[int]:
        """
        求s每个位置距离c的最近的距离
        左右扫两边即可
        这种问题是`到最近的`
        """
        left = [-int(1e20)] * len(s)
        for i in range(len(s)):
            if s[i] == c:
                left[i] = i
            elif i > 0:
                left[i] = left[i - 1]

        right = [int(1e20)] * len(s)
        for i in range(len(s) - 1, -1, -1):
            if s[i] == c:
                right[i] = i
            elif i < len(s) - 1:
                right[i] = right[i + 1]

        res = []
        for i in range(len(s)):
            min_ = min(i - left[i], right[i] - i)
            res.append(min_)

        return res


print(Solution().solve(s="aabaab", c="b"))
