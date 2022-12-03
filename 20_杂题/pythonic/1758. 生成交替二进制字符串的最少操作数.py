# 1758. 生成交替二进制字符串的最少操作数
class Solution:
    def minOperations(self, s: str) -> int:
        evenZero = s[::2].count("0")
        oddOne = s[1::2].count("1")
        res = evenZero + oddOne
        return min(res, len(s) - res)


# 输入：s = "0100"
# 输出：1
# 解释：如果将最后一个字符变为 '1' ，s 就变成 "0101" ，即符合交替字符串定义。
