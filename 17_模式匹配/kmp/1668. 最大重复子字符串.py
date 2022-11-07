# 最大重复子字符串
# 单词 word 的 最大重复值 是单词 word 在 sequence 中最大的重复值。
# 给你一个字符串 sequence 和 word ，请你返回 最大重复值 k。


from typing import List


def findAll(string: str, target: str) -> List[int]:
    """找到所有匹配的字符串起始位置"""
    start = 0
    res = []
    while True:
        pos = string.find(target, start)
        if pos == -1:
            break
        else:
            res.append(pos)
            start = pos + 1

    return res


class Solution:
    def maxRepeating(self, sequence: str, word: str) -> int:
        # O(nlongn)
        left, right = 1, len(sequence) // len(word) + 1
        while left <= right:
            mid = (left + right) // 2
            if word * mid in sequence:
                left = mid + 1
            else:
                right = mid - 1
        return right

    def maxRepeating2(self, sequence: str, word: str) -> int:
        # !O(n) dp 先findAll/indexAll找到合法位置，再dp记录每个位置的最大重复值
        n, m = len(sequence), len(word)
        if n < m:
            return 0

        indexes = findAll(sequence, word)
        ok = [False] * n
        for i in indexes:
            ok[i] = True

        dp = [0] * n
        for i in range(n):
            if ok[i]:
                dp[i] = (dp[i - m] + 1) if i >= m else 1
        return max(dp)


print(Solution().maxRepeating("ababc", "ab"))
print(Solution().maxRepeating2("ababc", "ab"))
# 输出：2
# 解释："abab" 是 "ababc" 的子字符串。
