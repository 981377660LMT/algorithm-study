from collections import Counter
from typing import List

# 如果某个单词在其中一个句子中恰好出现一次，在另一个句子中却 没有出现 ，那么这个单词就是 不常见的 。


# 可以理解成拼接字符串A+B，然后返回拼接后的字符串中只出现过一次的单词
class Solution:
    def uncommonFromSentences(self, s1: str, s2: str) -> List[str]:
        return [k for k, v in Counter(s1.split() + s2.split()).items() if v == 1]

