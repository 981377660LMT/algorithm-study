# 0 ≤ n < 2 ** 31

# !5个空瓶可换一瓶新酒
# 现在有n瓶新酒 问最多可喝多少瓶
# 不可以借空瓶
class Solution:
    def solve(self, full):
        empty = full
        res = full

        while empty >= 5:
            count = empty // 5
            empty -= count * 5
            empty += count
            res += count
        return res


print(Solution().solve(10))
