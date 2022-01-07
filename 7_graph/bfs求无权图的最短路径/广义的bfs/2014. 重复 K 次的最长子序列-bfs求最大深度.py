from collections import deque

# 请你找出字符串 s 中 重复 k 次的 最长子序列 。
# 如果存在多个满足的子序列，则返回 字典序最大 的那个。如果不存在这样的子序列，返回一个 空 字符串。

# 1. 找出重复至少k次的字符
# 2. bfs 验证合起来的子序列是否出现k次
class Solution:
    def longestSubsequenceRepeatedK(self, s: str, k: int) -> str:
        def check(cand: str) -> bool:
            """Return True if cand is a k-repeated sub-sequence of s."""
            matchIndex, hit = 0, 0
            for char in s:
                if char == cand[matchIndex]:
                    matchIndex += 1
                    if matchIndex == len(cand):
                        hit += 1
                        if hit == k:
                            return True
                        matchIndex = 0
            return False

        counter = [0] * 26
        for char in s:
            counter[ord(char) - 97] += 1
        cand = [chr(i + 97) for i, count in enumerate(counter) if count >= k]

        res = ''
        queue = deque([''])
        while queue:
            cur = queue.popleft()
            for char in cand:
                next = cur + char
                if check(next):
                    # 如果存在多个满足的子序列，则返回 字典序最大 的那个
                    res = next
                    queue.append(next)

        return res


print(Solution().longestSubsequenceRepeatedK(s="letsleetcode", k=2))
# 输出："let"
# 解释：存在两个最长子序列重复 2 次：let" 和 "ete" 。
# "let" 是其中字典序最大的一个。
