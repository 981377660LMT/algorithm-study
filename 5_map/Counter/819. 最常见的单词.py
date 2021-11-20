# 返回出现次数最多，同时不在禁用列表中的单词。

from collections import Counter
from re import findall
from typing import List


class Solution:
    def mostCommonWord(self, paragraph: str, banned: List[str]) -> str:
        for (key, _) in Counter(findall(r"\w+", paragraph.lower())).most_common():
            if key not in banned:
                return key
        return ''

