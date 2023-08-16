from typing import List, Tuple, Optional
from collections import defaultdict, Counter
from sortedcontainers import SortedList

try:
    import sys

    sys.set_int_max_str_digits(0)

except AttributeError:
    pass


MOD = int(1e9 + 7)
INF = int(1e20)


class Solution:
    def a(self, s: str) -> List[str]:
        ...
