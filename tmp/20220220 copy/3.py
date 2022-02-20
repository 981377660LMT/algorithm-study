from typing import List, Tuple
from collections import Counter
from heapq import heappop, heappush, heapify
from collections import deque


class Solution:
    def repeatLimitedString(self, s: str, repeatLimit: int) -> str:
        counter = Counter(s)
        pq = [(-ord(char), char, freq) for char, freq in counter.items()]
        heapify(pq)

        res = []
        repeat = 0
        while pq:
            _, char, freq = heappop(pq)
            if repeat == repeatLimit and res and res[-1] == char:
                if pq:
                    _, nextChar, nextFreq = heappop(pq)
                    res.append(nextChar)
                    if nextFreq > 1:
                        heappush(pq, (-ord(nextChar), nextChar, nextFreq - 1))
                else:
                    break
            if res and res[-1] == char:
                repeat += 1
            else:
                repeat = 1
            if freq > 1:
                heappush(pq, (-ord(char), char, freq - 1))

            res.append(char)

        return ''.join(res)


print(Solution().repeatLimitedString(s="robnsdvpuxbapuqgopqvxdrchivlifeepy", repeatLimit=2))
