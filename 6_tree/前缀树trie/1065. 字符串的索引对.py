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
    def indexPairs(self, text: str, words: List[str]) -> List[List[int]]:
        res = []
        for w in words:
            starts = findAll(text, w)
            res.extend([[start, start + len(w) - 1] for start in starts])
        return sorted(res)


print(Solution().indexPairs("ababa", ["aba", "ab"]))
