# 472. 连接词
# https://leetcode.cn/problems/concatenated-words/description/
# 给你一个 不含重复 单词的字符串数组 words ，请你找出并返回 words 中的所有 连接词 。
# 连接词 定义为：一个完全由给定数组中的至少两个较短单词（不一定是不同的两个单词）组成的字符串。


from typing import List, Set


class Solution:
    def findAllConcatenatedWordsInADict(self, words: List[str]) -> List[str]:
        words.sort(key=len)
        wordset: Set[str] = set()
        res: List[str] = []

        for w in words:
            if not w:
                continue
            if self._check(w, wordset):
                res.append(w)
            wordset.add(w)

        return res

    def _check(self, word: str, wordset: Set[str]) -> bool:
        """
        判断 word 是否能由 wordset 中的两个或更多单词拼接而成。
        用 dp[i] 表示 word[:i] 能否被拆分，dp[0]=True。
        对每个 i 满足 dp[i]，尝试所有 j>i：
          如果 word[i:j] 在 wordset 中，则标记 dp[j]=True。
        一旦 dp[n] 为 True，就说明 word = ...+... 拼接而成。
        """
        n = len(word)
        dp = [False] * (n + 1)
        dp[0] = True

        for i in range(n):
            if not dp[i]:
                continue
            for j in range(i + 1, n + 1):
                if word[i:j] in wordset:
                    dp[j] = True
                    if j == n:
                        return True
        return False


if __name__ == "__main__":
    sol = Solution()
    words = [
        "cat",
        "cats",
        "catsdogcats",
        "dog",
        "dogcatsdog",
        "hippopotamus",
        "rat",
        "ratcatdogcat",
    ]
    print(sol.findAllConcatenatedWordsInADict(words))
