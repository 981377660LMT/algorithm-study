from typing import List
from itertools import zip_longest


class Solution:
    def printVertically(self, s: str) -> List[str]:
        return [''.join(i).rstrip() for i in zip_longest(*s.split(), fillvalue=' ')]


print(Solution().printVertically(s="TO BE OR NOT TO BE"))
# 输出：["TBONTB","OEROOE","   T"]
# 解释：题目允许使用空格补位，但不允许输出末尾出现空格。
# "TBONTB"
# "OEROOE"
# "   T"

