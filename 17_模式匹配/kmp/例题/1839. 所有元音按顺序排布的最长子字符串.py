# 只包含元音
from itertools import groupby
from typing import List


def findAll(string: str, target: str) -> List[int]:
    """找到所有匹配的字符串起始位置"""
    start = 0
    res = []
    while True:
        pos = string.find(target, start)
        if pos == -1:
            break
        else:
            res.append(pos)
            start = pos + 1

    return res


class Solution:
    def longestBeautifulSubstring(self, word: str) -> int:
        res, length, type = 0, 1, 1
        for right in range(1, len(word)):
            if word[right] >= word[right - 1]:
                length += 1
            if word[right] > word[right - 1]:
                type += 1
            if word[right] < word[right - 1]:
                length, type = 1, 1
            if type == 5:
                res = max(res, length)

        return res

    def longestBeautifulSubstring2(self, word: str) -> int:
        groups = [[char, len(list(group))] for char, group in groupby(word)]
        chars = ''.join([c for c, _ in groups])
        starts = findAll(chars, 'aeiou')
        preSum = [0]
        for i, (_, count) in enumerate(groups):
            preSum.append(preSum[i] + count)

        res = 0
        for start in starts:
            res = max(res, preSum[start + 5] - preSum[start])
        return res


print(Solution().longestBeautifulSubstring2(word="aeiaaioaaaaeiiiiouuuooaauuaeiu"))
# 输出：13
# 解释：最长子字符串是 "aaaaeiiiiouuu" ，长度为 13 。
