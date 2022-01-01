from typing import List

# 给你一个「有效括号字符串」 seq，请你将其分成`两个`不相交的有效括号字符串，A 和 B，并使这两个字符串的深度最小。
# 划分方案用一个长度为 seq.length 的答案数组 answer 表示，编码规则如下：

# 为0的部分对应seq的括号是A字符串，为1的部分对应seq的括号是B字符串
class Solution:
    def maxDepthAfterSplit(self, seq: str) -> List[int]:
        res = []
        level = 0
        for char in seq:
            if char == '(':
                level += 1
                res.append(level % 2)
            if char == ')':
                res.append(level % 2)
                level -= 1
        return res


print(Solution().maxDepthAfterSplit("()(())()"))
# 输出：[0,0,0,1,1,0,1,1]
# 解释：本示例答案不唯一。
# 按此输出 A = "()()", B = "()()", max(depth(A), depth(B)) = 1，它们的深度最小。
# 像 [1,1,1,0,0,1,1,1]，也是正确结果，其中 A = "()()()", B = "()", max(depth(A), depth(B)) = 1 。

# 解题关键：
# the index parity of a pair '(', ')' must be opposite, say '(' at an odd position ,
# the corresponding ')' must be in even position.
# So, we assign '(' at odd position and ')' at even position to B, others to A.
