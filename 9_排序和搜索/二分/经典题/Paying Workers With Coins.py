from bisect import bisect_left


MOD = int(1e9 + 7)


class Solution:
    def solve(self, coins, salaries):
        """返回发工资的方案数"""
        coins = sorted(coins)
        salaries = sorted(salaries, reverse=True)

        res = 1
        for i, salary in enumerate(salaries):
            inValid = bisect_left(coins, salary)
            valid = len(coins) - inValid
            res *= valid - i
            res %= MOD

        return res


print(Solution().solve(coins=[1, 2, 3], salaries=[1, 2]))

