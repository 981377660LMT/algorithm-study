# 1081. 包含所有不同字符的字典序最小子序列
# 返回 s 字典序最小的子序列，该子序列包含 s 的所有不同字符，且只包含一次。

# 需要保证每个数最后正好出现一次
# !是否能加：使用visited保证最多出现一次
# !是否能删：使用remain防止多删

# !使用栈存放答案，并用一个集合记录栈内字母，保证每个字母只出现一次
# 对于每个字符：
# 如果该字符已经在栈中，跳过；
# 否则，在栈不空、当前字符小于栈顶字符、且栈顶字符后面还会出现时，弹出栈顶字符，从集合中移除。
# 最后将当前字符入栈，并加入集合。

from collections import Counter


class Solution:
    def smallestSubsequence(self, s: str) -> str:
        """返回 s 字典序最小的子序列，该子序列包含 s 的所有不同字符，且只包含一次。"""
        remain = Counter(s)
        visited = set()
        stack = []

        for cur in s:
            remain[cur] -= 1
            # !是否能加  (最后这个肯定要删的 既然要删 不如早删)
            if cur in visited:
                continue

            # !是否能删  (如果后面凑不齐这个字符了 就不删)
            while stack and stack[-1] > cur and remain[stack[-1]] > 0:
                visited.remove(stack.pop())
            stack.append(cur)
            visited.add(cur)

        return "".join(stack)
