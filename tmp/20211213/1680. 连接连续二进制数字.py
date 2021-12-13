# 请你将 1 到 n 的二进制表示连接起来，并返回连接结果对应的 十进制 数字对 109 + 7 取余的结果。
# 1 <= n <= 105


MOD = int(1e9 + 7)


# 遇到二的整数次幂移位量就加一
class Solution:
    def concatenatedBinary(self, n: int) -> int:
        res = 0
        shiftLen = 0
        for num in range(1, n + 1):
            if num & (num - 1) == 0:
                shiftLen += 1
            res = (res << shiftLen) + num
            res %= MOD

        return res


print(Solution().concatenatedBinary(3))
# 输出：27
# 解释：二进制下，1，2 和 3 分别对应 "1" ，"10" 和 "11" 。
# 将它们依次连接，我们得到 "11011" ，对应着十进制的 27 。
