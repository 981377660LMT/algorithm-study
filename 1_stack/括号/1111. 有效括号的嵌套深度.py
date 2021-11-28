from typing import List

# 给你一个「有效括号字符串」 seq，请你将其分成`两个`不相交的有效括号字符串，A 和 B，并使这两个字符串的深度最小。
# 划分方案用一个长度为 seq.length 的答案数组 answer 表示，编码规则如下：
class Solution:
    def maxDepthAfterSplit(self, seq: str) -> List[int]:
        return [i & 1 ^ (seq[i] == '(') for i in range(len(seq))]


print(Solution().maxDepthAfterSplit("()(())()"))
# 输出：[0,0,0,1,1,0,1,1]
# 解释：本示例答案不唯一。
# 按此输出 A = "()()", B = "()()", max(depth(A), depth(B)) = 1，它们的深度最小。
# 像 [1,1,1,0,0,1,1,1]，也是正确结果，其中 A = "()()()", B = "()", max(depth(A), depth(B)) = 1 。

# 解题关键：
# the index parity of a pair '(', ')' must be opposite, say '(' at an odd position ,
# the corresponding ')' must be in even position.
# So, we assign '(' at odd position and ')' at even position to B, others to A.
