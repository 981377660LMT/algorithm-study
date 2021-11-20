# s 可以由 words 中的前 k（k 为 正数 ）个字符串按顺序相连得到，
from typing import List


class Solution:
    def isPrefixString(self, s: str, words: List[str]) -> bool:
        return any(s == ''.join(words[:i]) for i in range(1, len(words) + 1))

