# 418. 屏幕可显示句子的数量
# https://leetcode.cn/problems/sentence-screen-fitting/
# 请你计算出给定句子可以在屏幕上完整显示的次数。
# 给你一个 rows x cols 的屏幕和一个用 非空 的单词列表组成的句子
# 在一行中 的两个连续单词必须用一个空格符分隔。
# 请你计算出给定句子可以在屏幕上完整显示的次数。

from typing import List


class Solution:
    def wordsTyping(self, sentence: List[str], rows: int, cols: int) -> int:
        """
        计算给定句子在 rows x cols 屏幕上能完整显示的次数。
        时间复杂度 O(n + rows)，空间复杂度 O(n)。
        """
        n = len(sentence)
        nextIndex = [0] * n  # 当前行以 sentence[i] 开头时，下一行应从哪个单词开始
        times = [0] * n  # 当前行以 sentence[i] 开头时，完整放下整句的次数

        for i in range(n):
            count, width = 0, 0  # 完整放下整句的次数，当前行已用宽度
            j = i
            while width + len(sentence[j]) <= cols:
                width += len(sentence[j]) + 1
                j += 1
                if j == n:
                    j = 0
                    count += 1
            nextIndex[i] = j
            times[i] = count

        # 优化：转移可以倍增/循环节加速
        res = 0
        wordIndex = 0
        for _ in range(rows):
            res += times[wordIndex]
            wordIndex = nextIndex[wordIndex]
        return res


if __name__ == "__main__":
    sol = Solution()
    print(sol.wordsTyping(["hello", "world"], 2, 8))  # → 1
    print(sol.wordsTyping(["a", "bcd", "e"], 3, 6))  # → 2
    print(sol.wordsTyping(["I", "had", "apple", "pie"], 4, 5))  # → 1
