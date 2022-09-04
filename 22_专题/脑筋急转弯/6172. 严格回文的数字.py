# !如果一个整数 n 在 b 进制下（b 为 2 到 n - 2 之间的所有整数）
# !对应的字符串 全部 都是 回文的 ,
# 那么我们称这个数 n 是 严格回文 的。

# 给你一个整数 n ,如果 n 是 严格回文 的,请返回 true ,
# 否则返回 false


class Solution:
    def isStrictlyPalindromic(self, n: int) -> bool:
        """
        数字 4 在二进制下不是回文的。
        对于 n ≥ 5,它们的 (n - 2) 进制表示都是 12,因此也都不是回文的。
        直接返回 false 即可
        """
        return False

    def isStrictlyPalindromic2(self, n: int) -> bool:
        def toString(num: int, radix: int) -> str:
            """将数字转换为指定进制的字符串"""
            if num < 0:
                return "-" + toString(-num, radix)

            if num == 0:
                return "0"

            res = []
            while num:
                div, mod = divmod(num, radix)
                res.append(str(mod))
                num = div
            return "".join(res)[::-1] or "0"

        for r in range(2, n - 2 + 1):
            s = toString(n, r)
            if s != s[::-1]:
                return False
        return True
