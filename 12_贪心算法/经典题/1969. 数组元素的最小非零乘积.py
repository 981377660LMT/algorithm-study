# 这个数组包含范围 [1, 2^p - 1] 内所有整数的二进制形式（两端都 包含）
# 从 nums 中选择两个元素 x 和 y  。
# 选择 x 中的一位与 y 对应位置的位交换。对应位置指的是两个整数 相同位置 的二进制位。

# 请你算出进行以上操作 任意次 以后，nums 能得到的 最小非零 乘积。将乘积对 109 + 7 取余 后返回。


# 总结：配对
# 1,2,3,4,5,6,7 => 1,6 2,5 3,4 7
# 配对后的和为(2^p-1) 即全为1 一个位置为1 另一个位置必定为0  交换后的结果为 1 和 2^p-2
# 一共有（2^(p-1)-1)对 再算上 不配对的2^p-1

MOD = int(1e9 + 7)


class Solution:
    def minNonZeroProduct(self, p: int) -> int:
        return ((1 << p) - 1) * pow((1 << p) - 2, (1 << (p - 1)) - 1, MOD) % MOD

    def minNonZeroProduct2(self, p: int) -> int:
        def qpow(base: int, exp: int, mod: int) -> int:
            res = 1

            while exp:
                if exp & 1:
                    res *= base
                    res %= mod

                exp >>= 1
                base **= 2
                base %= mod

            return res

        return ((1 << p) - 1) * qpow((1 << p) - 2, (1 << (p - 1)) - 1, MOD) % MOD


print(Solution().minNonZeroProduct2(p=3))
# 输入：p = 3
# 输出：1512
# 解释：nums = [001, 010, 011, 100, 101, 110, 111]
# - 第一次操作中，我们交换第二个和第五个元素最左边的数位。
#     - 结果数组为 [001, 110, 011, 100, 001, 110, 111] 。
# - 第二次操作中，我们交换第三个和第四个元素中间的数位。
#     - 结果数组为 [001, 110, 001, 110, 001, 110, 111] 。
# 数组乘积 1 * 6 * 1 * 6 * 1 * 6 * 7 = 1512 是最小乘积。
