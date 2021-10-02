# 例如，当相似单词对是 pairs = [["great", "fine"], ["acting","drama"], ["skills","talent"]]的时候，
# "great acting skills" 和 "fine drama talent" 是相似的。

# 注意相似关系是'不'具有传递性的。例如，如果 "great" 和 "fine" 是相似的，
# "fine" 和 "good" 是相似的，但是 "great" 和 "good" 未必是相似的。


from typing import List


class Solution:
    def areSentencesSimilar(
        self, sentence1: List[str], sentence2: List[str], similarPairs: List[List[str]]
    ) -> bool:
        if len(sentence1) != len(sentence2):
            return False
        pairset = set(map(tuple, similarPairs))  # 相似单词对转元组后去重
        return all(
            w1 == w2 or (w1, w2) in pairset or (w2, w1) in pairset
            for w1, w2 in zip(sentence1, sentence2)
        )

