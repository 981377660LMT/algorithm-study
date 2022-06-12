from collections import defaultdict
from typing import List, Optional, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)

# 1 <= sub.length <= s.length <= 5000
class Solution:
    def matchReplacement(self, s: str, sub: str, mappings: List[List[str]]) -> bool:
        adjMap = defaultdict(set)
        for pre, cur in mappings:
            adjMap[pre].add(cur)
            adjMap[pre].add(pre)

        # !枚举起点
        for start in range(len(s)):
            i = start
            j = 0
            while i < len(s) and j < len(sub):
                if s[i] in adjMap[sub[j]] or s[i] == sub[j]:
                    j += 1
                    i += 1
                else:
                    break
            if j == len(sub):
                return True

        return False


print(
    Solution().matchReplacement(
        s="fool3e7bar", sub="leet", mappings=[["e", "3"], ["t", "7"], ["t", "8"]]
    )
)
