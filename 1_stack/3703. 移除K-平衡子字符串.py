# 3703. 移除K-平衡子字符串
# https://leetcode.cn/problems/remove-k-balanced-substrings/
# 给你一个只包含 '(' 和 ')' 的字符串 s，以及一个整数 k。
# 如果一个 字符串 恰好是 k 个 连续 的 '(' 后面跟着 k 个 连续 的 ')'，即 '(' * k + ')' * k ，那么称它是 k-平衡 的。
# 例如，如果 k = 3，k-平衡字符串是 "((()))"。
# 你必须 重复地 从 s 中移除所有 不重叠 的 k-平衡子串，然后将剩余部分连接起来。持续这个过程直到不存在 k-平衡 子串 为止。
# 返回所有可能的移除操作后的最终字符串。


class Solution:
    def removeSubstring(self, s: str, k: int) -> str:
        stack = []  # (char, count)
        for c in s:
            if stack and stack[-1][0] == c:
                stack[-1][1] += 1
            else:
                stack.append([c, 1])

            if c == ")" and len(stack) >= 2 and stack[-1][1] == k and stack[-2][1] >= k:
                stack.pop()
                stack[-1][1] -= k
                if stack[-1][1] == 0:
                    stack.pop()

        return "".join([c * count for c, count in stack])
