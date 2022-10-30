from typing import List


# 给你两个字符串数组 queries 和 dictionary 。
# 数组中所有单词都只包含小写英文字母，且长度都相同。
# 一次 编辑 中，你可以从 queries 中选择一个单词，
# !将任意一个字母修改成任何其他字母。
# 从 queries 中找到所有满足以下条件的字符串：不超过 两次编辑内，字符串与 dictionary 中某个字符串相同。
# 请你返回 queries 中的单词列表，这些单词距离 dictionary 中的单词 编辑次数 不超过 两次 。
# 单词返回的顺序需要与 queries 中原本顺序相同。


# !因为只能`修改`,所以只需计算两个字符串的各位diff数之和
def calDist(word1: str, word2: str) -> int:
    return sum(c1 != c2 for c1, c2 in zip(word1, word2))


class Solution:
    def twoEditWords(self, queries: List[str], dictionary: List[str]) -> List[str]:
        S = set(dictionary)
        res = []
        for q in queries:
            for word in S:
                if len(q) != len(word):  # !长度不同,跳过
                    continue
                if calDist(q, word) <= 2:
                    res.append(q)
                    break
        return res
