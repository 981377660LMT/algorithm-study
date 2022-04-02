# Given three integer lists a, b and c,
# find the minimum possible abs(a[i] - b[j]) + abs(b[j] - c[k]) for any i, j and k.
# ni<=1e5

# 思路：对每个b[j] 找离他最近的元素
from bisect import bisect_right


class Solution:
    def solve(self, a, b, c):
        def getNearest(nums, target):
            pos = bisect_right(nums, target)
            left = nums[pos - 1] if pos > 0 else int(1e20)
            right = nums[pos] if pos < len(nums) else int(1e20)
            if abs(target - left) < abs(target - right):
                return left
            return right

        a, b, c = sorted(a), sorted(b), sorted(c)
        res = int(1e20)
        for i in range(len(b)):
            cur = b[i]
            n1, n2 = getNearest(a, cur), getNearest(c, cur)
            res = min(res, abs(n1 - cur) + abs(n2 - cur))
        return res


print(Solution().solve(a=[1, 8, 5], b=[2, 9], c=[5, 4]))
# We can pick a[0], b[0] and c[1]
