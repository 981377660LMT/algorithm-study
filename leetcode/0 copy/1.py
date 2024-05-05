import string
from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 有效单词 需要满足以下几个条件：

# 至少 包含 3 个字符。
# 由数字 0-9 和英文大小写字母组成。（不必包含所有这类字符。）
# 至少 包含一个 元音字母 。
# 至少 包含一个 辅音字母 。
# 给你一个字符串 word 。如果 word 是一个有效单词，则返回 true ，否则返回 false 。

# 注意：


# 'a'、'e'、'i'、'o'、'u' 及其大写形式都属于 元音字母 。
# 英文中的 辅音字母 是指那些除元音字母之外的字母。
VOWEL = "aeiouAEIOU"
NOT_VOWEL = "".join(set(string.ascii_letters) - set(VOWEL))
VALID = string.ascii_letters + string.digits


class Solution:
    def isValid(self, word: str) -> bool:
        n = len(word)
        if n < 3:
            return False
        if any(c not in VALID for c in word):
            return False
        if all(c not in VOWEL for c in word):
            return False
        if all(c not in NOT_VOWEL for c in word):
            return False
        return True
