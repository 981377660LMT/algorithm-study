#  求至少走过target次的区间数
#  if you walk right 2 times, then you walked on blocks [0, 1] and [1, 2] once each.
#  If you walk left once, then you'd walk on block [1, 2] again.
# Return the number of blocks that's been walked on at least target number of times.
from collections import defaultdict

# 0 ≤ n ≤ 100,000 n为区间个数
# 没有告诉值域范围


class Solution:
    def solve(self, walks, target):
        pos = 0
        diff = defaultdict(int)
        for num in walks:
            if num > 0:
                diff[pos] += 1
                diff[pos + num] -= 1
            elif num < 0:
                diff[pos + num] += 1
                diff[pos] -= 1
            pos += num

        res = 0
        curSum = 0
        prePos = 0
        for pos, cur in sorted(diff.items()):
            if curSum >= target:
                res += pos - prePos
            curSum += cur
            prePos = pos

        return res


print(Solution().solve(walks=[2, -4, 1], target=2))
# We move right 2 steps right and then 4 steps left and then 1 step right.
#  So we step on blocks [-2, -1], [0, 1] and [1, 2] 2 times.
