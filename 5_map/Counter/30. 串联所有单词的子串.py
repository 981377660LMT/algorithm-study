from typing import List
from collections import Counter

# 长度相同 的单词 words
# 找出 s 中恰好可以由 words 中所有单词串联形成的子串的起始位置。
class Solution:
    def findSubstring(self, s: str, words: List[str]) -> List[int]:
        if not s or not words:
            return []
        res = []
        n = len(words)
        word_len = len(words[0])
        window_len = n * word_len
        target = Counter(words)

        # 对每个可能的起始位置
        i = 0
        while i + window_len <= len(s):
            sliced = []
            start = i
            for _ in range(n):
                sliced.append(s[start : start + word_len])
                start += word_len
            print(Counter(sliced), target)
            if Counter(sliced) == target:
                res.append(i)
            i += 1
        return res


s = Solution()
print(s.findSubstring("barfoothefoobarman", ["foo", "bar"]))

# 输出：[0,9]
# 解释：
# 从索引 0 和 9 开始的子串分别是 "barfoo" 和 "foobar" 。
# 输出的顺序不重要, [9,0] 也是有效答案。
# Counter({'bar': 1, 'foo': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'arf': 1, 'oot': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'rfo': 1, 'oth': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'foo': 1, 'the': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'oot': 1, 'hef': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'oth': 1, 'efo': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'the': 1, 'foo': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'hef': 1, 'oob': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'efo': 1, 'oba': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'foo': 1, 'bar': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'oob': 1, 'arm': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'oba': 1, 'rma': 1}) Counter({'foo': 1, 'bar': 1})
# Counter({'bar': 1, 'man': 1}) Counter({'foo': 1, 'bar': 1})
# [0, 9]
