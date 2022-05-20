# 需要保证每个数最后正好出现一次
# 不多：使用visited保证最多出现一次
# 不少：使用remain_counter防止多删


from collections import Counter


class Solution:
    def removeDuplicateLetters(self, s) -> str:
        """返回 s 字典序最小的子序列，该子序列包含 s 的所有不同字符，且只包含一次。"""
        stack = []
        visited = set()
        remain = Counter(s)
        for char in s:
            if char not in visited:
                while stack and stack[-1] > char and remain[stack[-1]]:
                    visited.discard(stack.pop())
                stack.append(char)
                visited.add(char)
            remain[char] -= 1
        return ''.join(stack)


print(Solution().removeDuplicateLetters("cbacdcbc"))
