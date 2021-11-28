from typing import List
from itertools import product

# S 的长度不超过12。


class Solution:
    def letterCasePermutation(self, s: str) -> List[str]:
        available = [set([char.lower(), char.upper()]) for char in s]
        # print(*product(*available))
        return [''.join(t) for t in product(*available)]


print(Solution().letterCasePermutation(s="a1b2"))
# 输出：["a1b2", "a1B2", "A1b2", "A1B2"]
