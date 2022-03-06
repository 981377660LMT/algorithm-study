from typing import List

# 注意ord 与chr的使用以及看范围

# s.length == 5
# 'A' <= s[0] <= s[3] <= 'Z'
# '1' <= s[1] <= s[4] <= '9'


class Solution:
    def cellsInRange(self, s: str) -> List[str]:
        c1, c2 = ord(s[0]), ord(s[3])
        r1, r2 = int(s[1]), int(s[4])
        res = []
        for c in range(c1, c2 + 1):
            for r in range(r1, r2 + 1):
                res.append(chr(c) + str(r))
        return res

