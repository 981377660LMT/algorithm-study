class Solution:
    def solve(self, n):
        # 有多少组连续正整数和为n
        res = 0
        if n == 0:
            return 1

        itemCount = 1  # 项数
        while itemCount * (itemCount - 1) // 2 < n:
            base = n - itemCount * (itemCount - 1) // 2
            if base % itemCount == 0:
                res += 1
            itemCount += 1
        return res


print(Solution().solve(n=9))
# n<2**31
# The possible lists are: [2, 3, 4], [4, 5], and [9].
