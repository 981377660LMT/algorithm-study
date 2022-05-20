# 将 s 分割成 2 个字符串 p 和 q ，
# 它们连接起来等于 s 且 p 和 q 中不同字符的数目相同。
# 请你返回 s 中好分割的数目。


from collections import Counter
from typing import List


class Solution:
    def numSplits(self, s: str) -> int:
        """前后缀分解"""

        def countType(string: str) -> List[int]:
            res = [0] * len(string)
            counter = Counter()
            for i, char in enumerate(string):
                counter[char] += 1
                res[i] = len(counter)
            return res

        pre, suf = countType(s), countType(s[::-1])[::-1]
        return sum(pre[i] == suf[i + 1] for i in range(len(s) - 1))
