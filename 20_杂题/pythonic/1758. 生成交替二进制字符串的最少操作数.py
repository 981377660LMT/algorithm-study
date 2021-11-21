class Solution:
    def minOperations(self, s: str) -> int:
        order = s[::2].count('0') + s[1::2].count('1')
        return min(order, len(s) - order)


# 输入：s = "0100"
# 输出：1
# 解释：如果将最后一个字符变为 '1' ，s 就变成 "0101" ，即符合交替字符串定义。
