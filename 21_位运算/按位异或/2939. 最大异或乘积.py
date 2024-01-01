# 2939. 最大异或乘积
# https://leetcode.cn/problems/maximum-xor-product/

# 给你三个整数 a ，b 和 n ，请你返回 (a XOR x) * (b XOR x) 的 最大值 且 x 需要满足 0 <= x < 2n。
# 由于答案可能会很大，返回它对 109 + 7 取余 后的结果。
#
# !贪心：遍历每一位，如果异或后的乘积比之前的最大值大，就异或.


MOD = int(1e9 + 7)


class Solution:
    def maximumXorProduct(self, a: int, b: int, n: int) -> int:
        res = a * b
        for i in range(n):
            cur = (a ^ (1 << i)) * (b ^ (1 << i))
            if cur > res:
                a ^= 1 << i
                b ^= 1 << i
                res = cur
        return res % MOD

    def maximumXorProduct2(self, a: int, b: int, n: int) -> int:
        """
        O(1).
        https://leetcode.cn/problems/maximum-xor-product/solutions/2532915/o1-zuo-fa-wei-yun-suan-de-qiao-miao-yun-lvnvr/
        """
        if a < b:
            a, b = b, a  # 保证 a >= b

        mask = (1 << n) - 1
        ax = a & ~mask  # 第 n 位及其左边，无法被 x 影响，先算出来
        bx = b & ~mask
        a &= mask  # 低于第 n 位，能被 x 影响
        b &= mask

        left = a ^ b  # 可分配：a XOR x 和 b XOR x 一个是 1 另一个是 0
        one = mask ^ left  # 无需分配：a XOR x 和 b XOR x 均为 1
        ax |= one  # 先加到异或结果中
        bx |= one

        # 现在要把 left 分配到 ax 和 bx 中
        # 根据基本不等式（均值定理），分配后应当使 ax 和 bx 尽量接近，乘积才能尽量大
        if left > 0 and ax == bx:
            # 尽量均匀分配，例如把 1111 分成 1000 和 0111
            high_bit = 1 << (left.bit_length() - 1)
            ax |= high_bit
            left ^= high_bit
        # 如果 a & ~mask 更大，则应当全部分给 bx（注意最上面保证了 a>=b）
        bx |= left

        return ax * bx % MOD
