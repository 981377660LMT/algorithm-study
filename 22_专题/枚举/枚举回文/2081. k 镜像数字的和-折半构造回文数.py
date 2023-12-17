# 一个 k 镜像数字 指的是一个在十进制和 k 进制下从前往后读和从后往前读都一样的 没有前导 0 的 正 整数。
# 2 <= k <= 9
# 1 <= n <= 30

from enumeratePalindrome import emumeratePalindromeByLength


def check(num: int, radix: int) -> bool:
    """判断num在radix进制下是否为回文"""
    assert radix >= 2

    digits = []
    while num:
        div, mod = divmod(num, radix)
        digits.append(mod)
        num = div
    return digits == digits[::-1]


class Solution:
    def kMirror(self, k: int, n: int) -> int:
        res = []
        iter = emumeratePalindromeByLength(1, int(1e20))
        while len(res) < n:
            palindrome = int(next(iter))
            if check(palindrome, k):
                res.append(palindrome)
        return sum(res)


if __name__ == "__main__":
    print(Solution().kMirror(k=3, n=7))
    # 输出：499
    # 解释：
    # 7 个最小的 3 镜像数字和它们的三进制表示如下：
    #   十进制       三进制
    #     1          1
    #     2          2
    #     4          11
    #     8          22
    #     121        11111
    #     151        12121
    #     212        21212
    # 它们的和为 1 + 2 + 4 + 8 + 121 + 151 + 212 = 499 。
