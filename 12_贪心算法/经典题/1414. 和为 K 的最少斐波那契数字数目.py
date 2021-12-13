# 1 <= k <= 10^9
# 给你数字 k ，请你返回和为 k 的斐波那契数字的最少数目，其中，每个斐波那契数字都可以被使用多次。
class Solution:
    def findMinFibonacciNumbers(self, k: int) -> int:
        # 斐波那契数字为：1，1，2，3，5，8，13，……
        fibo = []
        a, b = 1, 1
        while True:
            if a > k:
                break
            fibo.append(a)
            a, b = b, a + b

        res = 0
        for fi in reversed(fibo):
            if k >= fi:
                k -= fi
                res += 1
        return res


print(Solution().findMinFibonacciNumbers(7))
# 输出：2
# 解释：斐波那契数字为：1，1，2，3，5，8，13，……
# 对于 k = 7 ，我们可以得到 2 + 5 = 7 。
