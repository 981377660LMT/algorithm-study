import math

# 给定一个正整数 N，试求有多少组连续正整数满足所有数字之和为 N?
class Solution:
    def consecutiveNumbersSum(self, n: int) -> int:
        res = 1  ## 自身
        for i in range(2, math.ceil(math.sqrt(2 * n))):
            if (n - i * (i - 1) / 2) % i == 0:
                res += 1
        return res


# 2*N=(2*a1+n-1)*n
# 即 a1*n=N-n*(n-1)/2 n为长度 n肯定不超过sqrt(2*N) 枚举n即可
