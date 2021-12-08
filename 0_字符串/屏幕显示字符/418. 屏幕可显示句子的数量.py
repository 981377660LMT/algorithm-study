from typing import List

# 请你计算出给定句子可以在屏幕上完整显示的次数。
# 给你一个 rows x cols 的屏幕和一个用 非空 的单词列表组成的句子

# 在一行中 的两个连续单词必须用一个空格符分隔。
# 请你计算出给定句子可以在屏幕上完整显示的次数。
class Solution:
    def wordsTyping(self, sentence: List[str], rows: int, cols: int) -> int:
        word_len = [len(word) for word in sentence]
        n = len(word_len)
        all_len = sum(word_len) + n  # 所有单词长度(加上空格,最后单词也有)
        res = 0
        idx = 0

        for _ in range(rows):
            remain_cols = cols
            while remain_cols > 0:
                if word_len[idx] <= remain_cols:
                    remain_cols -= word_len[idx]
                    if remain_cols > 0:
                        remain_cols -= 1
                    idx += 1

                    # 到头了，剩余的列的位置有放一个 sentence
                    if idx == n:
                        div, mod = divmod(remain_cols, all_len)
                        res += div + 1
                        remain_cols = mod
                        idx = 0
                else:
                    break

        return res


print(Solution().wordsTyping(rows=3, cols=6, sentence=["a", "bcd", "e"]))
# 输出：
# 2

# 解释：
# a-bcd-
# e-a---
# bcd-e-

# 字符 '-' 表示屏幕上的一个空白位置。
