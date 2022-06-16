# 每次操作会使排列变成它的前一个排列，
# 题目等价于求从小到大数第几个排列-1，
# 或者相当于求有多少个排列比它小
# 求在所有的`排列组合中`，当前这个是字典序第k小 (k从0开始)
# 有点类似康托展开的思想

from collections import Counter
from functools import lru_cache


# 1 <= s.length <= 3000

MOD = int(1e9 + 7)


@lru_cache(None)
def fac(n: int) -> int:
    """n的阶乘"""
    if n == 0:
        return 1
    return n * fac(n - 1) % MOD


@lru_cache(None)
def ifac(n: int) -> int:
    """n的阶乘的逆元"""
    return pow(fac(n), MOD - 2, MOD)


class Solution:
    def makeStringSorted(self, s: str) -> int:
        """求在所有不重复的组合中,当前这个是字典序第几小"""
        res, n = 0, len(s)
        counter = Counter(s)
        keys = sorted(counter)

        for char in s:
            # !后面位置的不重复组合数*当前位置比自己小的数个数
            mul1 = fac(n - 1)
            for count in counter.values():
                mul1 *= ifac(count)
                mul1 %= MOD

            for smaller in keys:
                if smaller >= char:
                    break
                res += counter[smaller] * mul1
                res %= MOD

            counter[char] -= 1
            n -= 1

        return res % MOD


print(Solution().makeStringSorted(s="cba"))
# 输出：5
# 解释：模拟过程如下所示：
# 操作 1：i=2，j=2。交换 s[1] 和 s[2] 得到 s="cab" ，然后反转下标从 2 开始的后缀字符串，得到 s="cab" 。
# 操作 2：i=1，j=2。交换 s[0] 和 s[2] 得到 s="bac" ，然后反转下标从 1 开始的后缀字符串，得到 s="bca" 。
# 操作 3：i=2，j=2。交换 s[1] 和 s[2] 得到 s="bac" ，然后反转下标从 2 开始的后缀字符串，得到 s="bac" 。
# 操作 4：i=1，j=1。交换 s[0] 和 s[1] 得到 s="abc" ，然后反转下标从 1 开始的后缀字符串，得到 s="acb" 。
# 操作 5：i=2，j=2。交换 s[1] 和 s[2] 得到 s="abc" ，然后反转下标从 2 开始的后缀字符串，得到 s="abc" 。
