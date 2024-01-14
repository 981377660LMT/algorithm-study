# 求小于等于n的哈沙德数(Harshad number)的个数
# harshad number: n能被n的各位数码之和整除
# n<=1e14
# https://atcoder.jp/contests/abc336/editorial/9055
# !需要枚举数码之和 1~9*m 才能确定当前数模 digitSum 的模数


import sys

input = lambda: sys.stdin.readline().rstrip("\r\n")
n = int(input())

if __name__ == "__main__":
    from functools import lru_cache

    def cal(targetDigitSum: int) -> int:
        @lru_cache(None)
        def dfs(pos: int, hasLeadingZero: bool, isLimit: bool, digitSum: int, curMod: int) -> int:
            if digitSum > targetDigitSum:
                return 0
            if pos == m:
                return int(not hasLeadingZero and curMod == 0 and digitSum == targetDigitSum)

            res = 0
            up = nums[pos] if isLimit else 9
            for cur in range(up + 1):
                if hasLeadingZero and cur == 0:
                    res += dfs(pos + 1, True, (isLimit and cur == up), digitSum, curMod)
                else:
                    nextDigitSum = digitSum + cur
                    nextMod = (curMod * 10 + cur) % targetDigitSum
                    res += dfs(pos + 1, False, (isLimit and cur == up), nextDigitSum, nextMod)
            return res

        res = dfs(0, True, True, 0, 0)
        dfs.cache_clear()
        return res

    nums = list(map(int, str(n)))
    m = len(nums)
    res = sum(cal(digitSum) for digitSum in range(1, 9 * m + 1))
    print(res)
