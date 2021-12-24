from collections import Counter
from heapq import heappop, heappush, heapify
from collections import deque

# 1.先把频率最大的字母扔进窗口(贪心)
# 2.借助最大堆和队列(k长度的滑窗)


class Solution:
    def rearrangeString(self, s: str, k: int) -> str:
        if k == 0:
            return s

        counter = Counter(s)
        pq = [(-freq, char) for char, freq in counter.items()]
        heapify(pq)

        res = []
        window = deque()

        while pq:
            freq, char = heappop(pq)
            freq *= -1

            res.append(char)
            window.append((freq - 1, char))
            # 凑齐k个
            if len(window) == k:
                freq, char = window.popleft()
                if freq > 0:
                    heappush(pq, (-freq, char))

        return ''.join(res) if len(res) == len(s) else ''


print(Solution().rearrangeString(s="aabbcc", k=3))
# 输出: "abcabc
# 在新的字符串中间隔至少 3 个单位距离。
