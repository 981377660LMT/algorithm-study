# 0 ≤ n < 2 ** 31

# 三个空瓶可换一瓶新酒
# 现在又n瓶新酒 问最多可喝多少瓶
class Solution:
    def solve(self, full):
        empty = 0
        druken = 0

        while full:
            druken += full
            empty += full
            full = empty // 3
            empty = empty % 3
        return druken

