# 给你两个正整数 n 和 k，二进制字符串  Sn 的形成规则如下：

# S1 = "0"
# 当 i > 1 时，Si = Si-1 + "1" + reverse(invert(Si-1))


# 其实就是一张纸对折
# mid = pow(2, n - 1) = 1 << (n - 1)


class Solution:
    def findKthBit(self, n: int, k: int) -> str:
        if n == 1:
            return '0'
        mid = 1 << (n - 1)
        if k == mid:
            return '1'
        elif k < mid:
            return self.findKthBit(n - 1, k)
        else:
            k = 2 * mid - k
            return '0' if self.findKthBit(n - 1, k) == '1' else '1'


print(Solution().findKthBit(n=4, k=11))
# 输出："1"
# 解释：S4 为 "011100110110001"，其第 11 位为 "1" 。
