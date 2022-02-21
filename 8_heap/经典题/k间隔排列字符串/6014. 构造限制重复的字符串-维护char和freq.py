from typing import List, Tuple
from collections import Counter
from heapq import heappop, heappush, heapify
from collections import deque

# 用 s 中的字符构造一个新字符串 repeatLimitedString ，
# 使任何字母 连续 出现的次数都不超过 repeatLimit 次。
# 你不必使用 s 中的全部字符。
# 返回 字典序最大的 repeatLimitedString 。

# 1405. 最长快乐字符串
class Solution:
    def repeatLimitedString(self, s: str, repeatLimit: int) -> str:
        counter = Counter(s)
        pq = [(-ord(char), char, freq) for char, freq in counter.items()]
        heapify(pq)

        res = []
        repeat = 0
        while pq:
            _, char, freq = heappop(pq)
            # 必须来下一个
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

            res.append(char)

            if freq > 1:
                heappush(pq, (-ord(char), char, freq - 1))

        return ''.join(res)


print(Solution().repeatLimitedString(s="robnsdvpuxbapuqgopqvxdrchivlifeepy", repeatLimit=2))
