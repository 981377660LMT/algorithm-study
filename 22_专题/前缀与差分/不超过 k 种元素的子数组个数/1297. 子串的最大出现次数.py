# 1 <= s.length <= 10^5
# 请你返回满足以下条件且出现次数最大的 任意 子串的出现次数：

# 子串中不同字母的数目必须小于等于 maxLetters 。
# 子串的长度必须大于等于 minSize 且小于等于 maxSize 。
from collections import defaultdict


class Solution:
    def maxFreq(self, s: str, maxLetters: int, minSize: int, maxSize: int) -> int:
        counter = defaultdict(int)
        for start in range(len(s) - minSize + 1):
            cur = s[start : start + minSize]
            if len(set(cur)) <= maxLetters:
                counter[s[start : start + minSize]] += 1
        return max(counter.values(), default=0)

