from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList, SortedDict, SortedSet
from bisect import bisect_left, bisect_right
from functools import lru_cache, reduce
from itertools import accumulate, groupby, combinations, permutations, product, chain, islice
from math import gcd, sqrt, ceil, floor, comb
from string import ascii_lowercase, ascii_uppercase, ascii_letters, digits
from operator import xor, or_, and_, not_

MOD = int(1e9 + 7)
INF = 2 ** 64
EPS = int(1e-8)
dirs4 = [[-1, 0], [0, 1], [1, 0], [0, -1]]
dirs8 = [[-1, 0], [-1, 1], [0, 1], [1, 1], [1, 0], [1, -1], [0, -1], [-1, -1]]


class Solution:
    def smallestNumber(self, num: int) -> int:
        if num == 0:
            return 0
        digits = [char for char in str(num)]
        if all(digit == '0' for digit in digits):
            return 0
        if num > 0:
            digits.sort()
            first = next(i for i, digit in enumerate(digits) if digit != '0')
            firstChar = digits[first]
            digits.pop(first)
            return int(firstChar + ''.join(digits))
        else:
            digits.sort(reverse=True)
            first = next(i for i, digit in reversed(list(enumerate(digits))) if digit != '0')
            firstChar = digits[first]
            digits.pop(first)
            return int(firstChar + ''.join(digits))


print(Solution().smallestNumber(num=310))
print(Solution().smallestNumber(num=-7605))

