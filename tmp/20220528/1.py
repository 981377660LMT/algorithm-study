from collections import Counter
from typing import List, Tuple

MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def digitCount(self, num: str) -> bool:
        for i, char in enumerate(num):
            if int(char) != num.count(str(i)):
                return False
        return True


# 修正：
class Solution:
    def digitCount(self, num: str) -> bool:
        counter = Counter(num)
        for i, char in enumerate(num):
            if int(char) != counter[str(i)]:
                return False
        return True
