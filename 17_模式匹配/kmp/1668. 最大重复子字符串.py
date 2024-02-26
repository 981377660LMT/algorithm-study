# 最大重复子字符串
# https://leetcode.cn/problems/maximum-repeating-substring/description/
# 单词 word 的 最大重复值 是单词 word 在 sequence 中最大的重复值。
# 给你一个字符串 sequence 和 word ，请你返回 最大重复值 k。

from kmp import indexOfAll


class Solution:
    def maxRepeating(self, sequence: str, word: str) -> int:
        # !O(n) dp 先findAll/indexAll找到合法位置，再dp记录每个位置的最大重复值
        n1, n2 = len(sequence), len(word)
        if n1 < n2:
            return 0

        indexes = indexOfAll(sequence, word)
        ok = [False] * n1
        for i in indexes:
            ok[i] = True

        dp = [0] * n1
        for i in range(n1):
            if ok[i]:
                dp[i] = (dp[i - n2] + 1) if i >= n2 else 1
        return max(dp)


if __name__ == "__main__":
    assert Solution().maxRepeating("ababc", "ab") == 2
    assert Solution().maxRepeating("ababc", "ba") == 1
