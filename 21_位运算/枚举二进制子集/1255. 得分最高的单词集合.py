from typing import List
from collections import Counter

# 1 <= words.length <= 14
# 就是检查每个单词取不取 2^14种
class Solution:
    def maxScoreWords(self, words: List[str], letters: List[str], score: List[int]) -> int:
        n = len(words)
        wordScores = [sum(score[ord(c) - ord('a')] for c in word) for word in words]
        freqCounter = Counter(letters)

        res = 0
        for state in range(1 << n):
            needs, select = Counter(), []
            for wordIndex in range(n):
                if (state >> wordIndex) & 1:
                    needs += Counter(words[wordIndex])
                    select.append(wordIndex)
            if freqCounter & needs == needs:
                curScore = sum(wordScores[i] for i in select)
                res = max(res, curScore)
        return res


print(
    Solution().maxScoreWords(
        words=["dog", "cat", "dad", "good"],
        letters=["a", "a", "c", "d", "d", "d", "g", "o", "o"],
        score=[1, 0, 9, 5, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0],
    )
)

