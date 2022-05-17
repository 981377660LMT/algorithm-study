# 最多能分成大小为k的几组
# 每个组内元素不可相同


class Solution:
    def solve(self, counts, k):
        def check(mid: int) -> bool:
            """"是否能分成mid组"""
            res = 0
            for num in counts:
                res += min(mid, num)
            return res >= k * mid

        left, right = 0, sum(counts)
        while left <= right:
            mid = (left + right) >> 1
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right


print(Solution().solve([3, 3, 2, 5], 2))

# Let's name the four item types [a, b, c, d] respectively. We can have the following groups of two where all elements are distinct types:

# (a, d)
# (b, d)
# (a, b)
# (a, b)
# (c, d)
# (c, d)
