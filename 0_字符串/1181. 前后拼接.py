from typing import List


# 请你按字典序排列并返回「新短语」列表，列表中的字符串应该是 不重复的 。
class Solution:
    def beforeAndAfterPuzzles(self, phrases: List[str]) -> List[str]:
        n = len(phrases)

        # 统计每个字符串（index表示）的第一个word和最后一个word
        pre_and_suf = [[]] * n
        for i, p in enumerate(phrases):
            words = p.split(' ')
            pre_and_suf[i] = [words[0], words[-1]]

        # 匹配
        match = set()
        for i in range(n):
            for j in range(n):
                if i == j:
                    continue
                if pre_and_suf[i][1] == pre_and_suf[j][0]:
                    new_p = phrases[i] + phrases[j][len(pre_and_suf[j][0]) :]  # 前后桥接的单词只出现1次
                    match.add(new_p)

        return sorted(list(match))


print(Solution().beforeAndAfterPuzzles(phrases=["writing code", "code rocks"]))
# 输出：["writing code rocks"]
