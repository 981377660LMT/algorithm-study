from typing import List
from collections import Counter


# 长度相同 的单词 words
# 找出 s 中恰好可以由 words 中所有单词串联形成的子串的起始位置。
class Solution:
    def findSubstring(self, s: str, words: List[str]) -> List[int]:
        # O(n^2) 暴力
        if not s or not words:
            return []

        res = []
        n = len(words)
        wordLen = len(words[0])
        windowLen = n * wordLen
        target = Counter(words)

        # 对每个可能的起始位置判断
        for i in range(len(s) - windowLen + 1):
            cur = Counter([s[left : left + wordLen] for left in range(i, i + windowLen, wordLen)])
            if cur == target:
                res.append(i)

        return res

    def findSubstring2(self, s: str, words: List[str]) -> List[int]:
        # O(n) 滑窗 `如果当前单词超出需要 那么左端点回来`
        res = []
        n = len(words)
        wordLen = len(words[0])
        target = Counter(words)

        # 分组滑窗
        for start in range(wordLen):
            left, count = start, 0
            counter = Counter()
            for right in range(start, len(s) - wordLen + 1, wordLen):
                cur = s[right : right + wordLen]
                counter[cur], count = counter[cur] + 1, count + 1
                while counter[cur] > target[cur]:
                    pre = s[left : left + wordLen]
                    counter[pre], count = counter[pre] - 1, count - 1
                    left += wordLen  # 这个单词不能要了
                if count == n:
                    res.append(left)

        return res


s = Solution()
print(s.findSubstring2("barfoothefoobarman", ["foo", "bar"]))

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
