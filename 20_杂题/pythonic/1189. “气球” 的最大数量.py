from collections import Counter


class Solution:
    def maxNumberOfBalloons(self, text: str) -> int:
        return min((d := Counter(text))['b'], d['a'], d['l'] // 2, d['o'] // 2, d['n'])
