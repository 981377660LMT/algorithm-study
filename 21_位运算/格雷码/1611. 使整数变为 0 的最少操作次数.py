# 格雷码转二进制

# 给你一个整数 n，你需要重复执行多次下述操作将其转换为 0 ：
# 翻转 n 的二进制表示中最右侧位（第 0 位）。
# 如果第 (i-1) 位为 1 且从第 (i-2) 位到第 0 位都为 0，则翻转 n 的二进制表示中的第 i 位。

# 很像玩九连环的过程。把最高位消掉以后，后面要循环，挨着一个一个慢慢消
class Solution:
    def minimumOneBitOperations(self, n: int) -> int:
        res = 0
        while n > 0:
            res ^= n
            n >>= 1
        return res


print(Solution().minimumOneBitOperations(n=3))
# 输入：n = 3
# 输出：2
# 解释：3 的二进制表示为 "11"
# "11" -> "01" ，执行的是第 2 种操作，因为第 0 位为 1 。
# "01" -> "00" ，执行的是第 1 种操作。
