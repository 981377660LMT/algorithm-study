# 返回 s 字典序最小的子序列，该子序列包含 s 的所有不同字符，且只包含一次。

# 需要保证每个数最后正好出现一次
# 是否能加：使用visited保证最多出现一次
# 是否能删：使用remain防止多删


from collections import Counter, defaultdict


class Solution:
    def removeDuplicateLetters(self, s) -> str:
        """返回 s 字典序最小的子序列，该子序列包含 s 的所有不同字符，且只包含一次。"""
        stack = []
        visited = defaultdict(int)
        remain = Counter(s)

        for char in s:
            # !是否能加  (最后这个肯定要删的 既然要删 不如早删)
            if visited[char] == 1:
                remain[char] -= 1
                continue

            # !是否能删  (如果后面凑不齐这个字符了 就不删)
            while stack and stack[-1] > char and remain[stack[-1]] >= 1:
                top = stack.pop()
                visited[top] -= 1

            stack.append(char)
            visited[char] += 1
            remain[char] -= 1

        return "".join(stack)


print(Solution().removeDuplicateLetters("cbacdcbc"))
