from typing import List
from itertools import zip_longest


class Solution:
    def printVertically(self, s: str) -> List[str]:
        words = s.split()
        return [''.join(col).rstrip() for col in zip_longest(*words, fillvalue=' ')]


print(Solution().printVertically(s="TO BE OR NOT TO BE"))
# 输出：["TBONTB","OEROOE","   T"]
# 解释：题目允许使用空格补位，但不允许输出末尾出现空格。
# "TBONTB"
# "OEROOE"
# "   T"

