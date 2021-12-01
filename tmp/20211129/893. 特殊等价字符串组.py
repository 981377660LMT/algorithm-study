from typing import List

# 如果两个字符串，奇数位的字母一样，偶数位的字母一样，这样的字符串相互等价# 返回 words 中 特殊等价字符串组 的数量。
# 1 <= words.length <= 1000
# 1 <= words[i].length <= 20
class Solution:
    def numSpecialEquivGroups(self, words: List[str]) -> int:
        res = set()
        for sub in words:
            sub = ''.join(sorted(sub[::2]) + sorted(sub[1::2]))
            res.add(sub)
        return len(res)
