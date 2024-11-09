# 给定一个正整数 s，令 A 为一个 n × n × n 的三维数组，其中每个元素 A[i][j][k] 定义为：
# A[i][j][k] = i * (j OR k)，其中 0 <= i, j, k < n。
# 返回使数组 A 中所有元素的和不超过 s 的 最大的 n。


def calc(upper: int, k: int) -> int:
    """[0, upper]中二进制第k(k>=0)位为1的数的个数.
    即满足 `num & (1 << k) > 0` 的数的个数
    """
    if k >= upper.bit_length():
        return 0
    res = upper // (1 << (k + 1)) * (1 << k)
    upper %= 1 << (k + 1)
    if upper >= 1 << k:
        res += upper - (1 << k) + 1
    return res


class Solution:
    def maxSizedArray(self, s: int) -> int:
        def check(mid: int) -> bool:
            pairs = mid * mid
            sumOr = 0
            # 对于每一位 bit，计算在所有 (j OR k) 中该位为 1 的总次数
            b = 0
            while True:
                ones = calc(mid - 1, b)
                if ones == 0:
                    break
                zeros = mid - ones
                anyOnes = pairs - zeros * zeros
                sumOr += anyOnes * (1 << b)
                b += 1
            sumI = mid * (mid - 1) // 2
            return sumI * sumOr <= s

        left, right = 1, int(1e16)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                left = mid + 1
            else:
                right = mid - 1
        return right
