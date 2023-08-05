# 2802. 找出第K个幸运数字
# https://leetcode.cn/problems/find-the-k-th-lucky-number/
# 我们知道 4 和 7 是 幸运 数字。同时，如果一个数字只包含幸运数字，那么它被称为幸运数字。
# 给定一个整数 k，返回第 k 个幸运数字，并将其表示为一个 字符串 。


# 1<=k<=1e9
# 保持大小顺序的映射：将4映射到0,7映射到1
# 相当于求第k个二进制数


class Solution:
    def kthLuckyNumber(self, k: int) -> str:
        def getLen(kth: int):
            """第kth个数幸运数字有几位."""
            res = 1
            while True:
                curCount = 1 << res
                if kth <= curCount:
                    break
                kth -= curCount
                res += 1
            return res

        n = getLen(k)
        for i in range(n):  # 减去前面的数(偏移量)
            k -= 1 << i
        kthBin = bin(k)[2:].zfill(n)
        kthBin = kthBin.replace("0", "4").replace("1", "7")
        return kthBin


print(Solution().kthLuckyNumber(1))
print(Solution().kthLuckyNumber(4))
