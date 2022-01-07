from collections import Counter

# 返回 s 中长度为 k 且 字典序最小 的子序列，
# 该子序列同时应满足字母 letter 出现 至少 repetition 次。
# 生成的测试用例满足 letter 在 s 中出现 至少 repetition 次。


class Solution:
    def smallestSubsequence(self, s: str, k: int, letter: str, repetition: int) -> str:
        need = repetition
        remain = sum(char == letter for char in s)
        stack = []
        for i, char in enumerate(s):
            # 还可以继续删:长度保证+重复次数保证
            while (
                stack
                and stack[-1] > char
                and len(stack) + len(s) - i > k
                and (stack[-1] != letter or need < remain)
            ):
                top = stack.pop()
                if top == letter:
                    need += 1

            # 还可以继续加
            if len(stack) < k and (char == letter or len(stack) + need < k):
                stack.append(char)
                if char == letter:
                    need -= 1
            if char == letter:
                remain -= 1

        return ''.join(stack)


print(Solution().smallestSubsequence(s="leet", k=3, letter="e", repetition=1))
# 输出："eet"
# 解释：存在 4 个长度为 3 ，且满足字母 'e' 出现至少 1 次的子序列：
# - "lee"（"leet"）
# - "let"（"leet"）
# - "let"（"leet"）
# - "eet"（"leet"）
# 其中字典序最小的子序列是 "eet" 。
