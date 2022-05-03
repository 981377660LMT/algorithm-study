from typing import List


# 缩写规则：
# 1. 初始缩写由起始字母+省略字母的数量+结尾字母组成。
# 2. 若存在冲突，则使用更长的前缀代替首字母，直到从单词到缩写的映射唯一
# 3. 若缩写并不比原单词更短，则保留原样。

# 贪心：
# 首先给每个单词选择最短的缩写。然后我们对于所有重复的单词，我们增加这些重复项的长度。
class Solution:
    def wordsAbbreviation(self, words: List[str]) -> List[str]:
        def compress(word: str, start: int = 0) -> str:
            if len(word) - start <= 3:
                return word
            return word[: start + 1] + str(len(word) - start - 2) + word[-1]

        n = len(words)
        res = list(map(compress, words))
        needStartFrom = [0] * n

        for i in range(n):
            while True:
                dup = set()
                for j in range(i + 1, n):
                    if res[i] == res[j]:
                        dup.add(j)

                if not dup:
                    break

                # 重复前缀的单词start+1 重新压缩
                dup.add(i)
                for dupeIndex in dup:
                    needStartFrom[dupeIndex] += 1
                    res[dupeIndex] = compress(words[dupeIndex], needStartFrom[dupeIndex])

        return res


print(
    Solution().wordsAbbreviation(
        words=[
            "like",
            "god",
            "internal",
            "me",
            "internet",
            "interval",
            "intension",
            "face",
            "intrusion",
        ]
    )
)

# 输出: ["l2e","god","internal","me","i6t","interval","inte4n","f2e","intr4n"]
