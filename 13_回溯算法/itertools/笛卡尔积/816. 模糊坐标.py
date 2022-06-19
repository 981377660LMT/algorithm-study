from itertools import product
from typing import Generator, List

# 4 <= S.length <= 12.
# S[0] = "(", S[S.length - 1] = ")", 且字符串 S 中的其他元素都是数字。


class Solution:
    def ambiguousCoordinates(self, s: str) -> List[str]:
        def gen(string: str) -> Generator[str, None, None]:
            n = len(string)
            for i in range(1, n + 1):
                left = string[:i]
                right = string[i:]
                # !1. 整数无前导零
                # !2. 小数无尾随零
                if (not left.startswith('0') or left == '0') and (not right.endswith('0')):
                    yield left + ('.' if i != n else '') + right

        s = s[1:-1]
        return [
            "({}, {})".format(*cand)
            for i in range(1, len(s))
            for cand in product(gen(s[:i]), gen(s[i:]))
        ]


print(Solution().ambiguousCoordinates(s="(100)"))
print(Solution().ambiguousCoordinates(s="(123)"))
print(Solution().ambiguousCoordinates(s="(0123)"))
