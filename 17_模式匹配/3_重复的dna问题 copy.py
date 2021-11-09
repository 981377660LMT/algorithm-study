from typing import List
from collections import Counter

LEN = 10


# 遍历所有长度为10的子串，统计哈希个数
class Solution:
    def findRepeatedDnaSequences(self, s: str) -> List[str]:
        return [
            k for k, v in Counter(s[i : i + LEN] for i in range(len(s) - LEN + 1)).items() if v > 1
        ]

