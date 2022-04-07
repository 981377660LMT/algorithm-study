# 你前面至少a个人 后面至少b个人
# 求你可能在那个位置 给出可能的个数


class Solution:
    def solve(self, n, a, b):
        return min(n - a, b + 1)


print(Solution().solve(10, 3, 4))

