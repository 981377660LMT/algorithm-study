from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 和一个正整数 k 。

# 用 vowels 和 consonants 分别表示字符串中元音字母和辅音字母的数量。

# 如果某个字符串满足以下条件，则称其为 美丽字符串 ：

# vowels == consonants，即元音字母和辅音字母的数量相等。
# (vowels * consonants) % k == 0，即元音字母和辅音字母的数量的乘积能被 k 整除。
# 返回字符串 s 中 非空美丽子字符串 的数量。

# 子字符串是字符串中的一个连续字符序列。

# 英语中的 元音字母 为 'a'、'e'、'i'、'o' 和 'u' 。


# 英语中的 辅音字母 为除了元音字母之外的所有字母。

VOWELS = set(["a", "e", "i", "o", "u"])


class Solution:
    def beautifulSubstrings(self, s: str, k: int) -> int:
        nums = [1 if c in VOWELS else 0 for c in s]
        res = 0
        for i in range(len(nums)):
            v, c = 0, 0
            for j in range(i, len(nums)):
                if nums[j] == 1:
                    v += 1
                else:
                    c += 1
                if v == c and v * c % k == 0:
                    res += 1
        return res
