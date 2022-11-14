# !统计回文子串个数
from Manacher import Manacher


class Solution:
    def countSubstrings(self, s: str) -> int:
        M = Manacher(s)
        return sum(M.oddRadius[i] + M.evenRadius[i] for i in range(len(s)))


assert Solution().countSubstrings(s="abc") == 3
assert Solution().countSubstrings(s="aaa") == 6
