class Solution:
    def solve(self, n, lower, upper):
        # 所有的数在[lower,upper]之间
        # 中间向两边扩展
        # 起点要最大
        if n > (upper - lower) * 2 + 1:
            return []

        left = right = upper - 1
        n -= 3
        while right > lower and n:
            right -= 1
            n -= 1
        while left > lower and n:
            left -= 1
            n -= 1

        return list(range(left, upper)) + list(range(upper, right - 1, -1))


print(Solution().solve(n=5, lower=2, upper=6))
print(Solution().solve(n=3, lower=2, upper=2))
