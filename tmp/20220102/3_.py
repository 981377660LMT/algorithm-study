from typing import List, Tuple
from collections import defaultdict, deque, Counter
from heapq import heapify, heappop, heappush
from sortedcontainers import SortedList
from bisect import bisect_left, bisect_right
from functools import lru_cache
from itertools import accumulate, groupby, combinations
from math import gcd

MOD = int(1e9 + 7)
INF = 0x7FFFFFFF


class Solution:
    def asteroidsDestroyed(self, mass: int, asteroids: List[int]) -> bool:
        asteroids.sort()

        cur = mass
        for num in asteroids:
            if cur < num:
                return False
            cur += num
        return True


print(Solution().asteroidsDestroyed(mass=10, asteroids=[3, 9, 19, 5, 21]))
print(Solution().asteroidsDestroyed(mass=5, asteroids=[4, 9, 23, 4]))
