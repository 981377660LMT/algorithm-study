from typing import List

# 给定一组 互不相同 的单词， 找出所有 不同 的索引对 (i, j)，使得列表中的两个单词， words[i] + words[j] ，可拼接成回文串。
# 1 <= words.length <= 5000
# 0 <= words[i].length <= 300


class Solution:
    def palindromePairs(self, words: List[str]) -> List[List[int]]:
        adjMap = {w[::-1]: i for i, w in enumerate(words)}
        res = set()

        for i, word in enumerate(words):
            for j in range(len(word) + 1):
                # 枚举后缀回文
                cur1 = word[j:]
                if cur1 == cur1[::-1]:
                    need1 = word[:j]
                    if need1 in adjMap and adjMap[need1] != i:
                        res.add((i, adjMap[need1]))
                # 枚举前缀回文
                cur2 = word[:j]
                if cur2 == cur2[::-1]:
                    need2 = word[j:]
                    if need2 in adjMap and adjMap[need2] != i:
                        res.add((adjMap[need2], i))

        return list(res)


print(Solution().palindromePairs(['abcd', 'dcba', 'lls', 's', 'sssll']))
