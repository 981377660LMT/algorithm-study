from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList


MOD = int(1e9 + 7)
INF = int(1e20)

# 给你一个字符串 s 和一个整数 k。请你使用以下算法加密字符串：


# 对于字符串 s 中的每个字符 c，用字符串中 c 后面的第 k 个字符替换 c（以循环方式）。
# 返回加密后的字符串。


class Solution:
    def getEncryptedString(self, s: str, k: int) -> str:
        n = len(s)
        sb = [""] * n
        for i in range(n):
            sb[i] = s[(i + k) % n]
        return "".join(sb)


#  s = "dart", k = 3

print(Solution().getEncryptedString("dart", 3))
