# 每收集一个硬币，上下左右的硬币就会消失
# 求最大的分

# 二维打家劫舍


def max(x, y):
    if x > y:
        return x
    return y


class Solution:
    def solve(self, matrix):
        def solve1D(row):
            """一维打家劫舍"""
            dp0, dp1 = 0, row[0]
            for i in range(1, len(row)):
                dp0, dp1 = max(dp0, dp1), max(dp0 + row[i], dp1)
            return max(dp0, dp1)

        return solve1D([solve1D(row) for row in matrix])


print(Solution().solve(matrix=[[1, 7, 6, 5], [9, 9, 3, 1], [4, 8, 1, 2]]))
