# 我们称一个数字字符串是 好数字 当它满足（下标从 0 开始）
# !偶数 下标处的数字为 偶数 且
# !奇数 下标处的数字为 质数 （2，3，5 或 7）。
# 给你一个整数 n ，请你返回长度为 n 且为好数字的数字字符串 总数 。

MOD = int(1e9 + 7)


class Solution:
    def countGoodNumbers(self, n: int) -> int:
        oddCount = n // 2
        evenCount = n - oddCount
        return pow(4, oddCount, MOD) * pow(5, evenCount, MOD) % MOD

