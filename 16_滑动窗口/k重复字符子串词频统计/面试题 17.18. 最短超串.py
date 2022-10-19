from collections import Counter
from typing import List

INF = int(1e18)

# 找到长数组中包含短数组所有的元素的最短子数组，其出现顺序无关紧要。
# 如有多个满足条件的子数组，返回左端点最小的一个
# 若不存在，返回空数组。

# big.length <= 100000
# 1 <= small.length <= 100000


# 最短超串


class Solution:
    def shortestSeq(self, big: List[int], small: List[int]) -> List[int]:
        need = set(small)
        left, n = 0, len(big)
        counter, kind = Counter(), 0  # 记录当前窗口中的元素 以及 当前窗口中的元素种类
        resLen, resLeft = INF, -1

        for right in range(n):
            if big[right] in need:
                counter[big[right]] += 1
                if counter[big[right]] == 1:
                    kind += 1

            while left <= right and kind == len(need):
                cand = right - left + 1
                if cand < resLen:
                    resLen = cand
                    resLeft = left
                if big[left] in need:
                    counter[big[left]] -= 1
                    if counter[big[left]] == 0:
                        kind -= 1
                left += 1

        return [] if resLen == INF else [resLeft, resLeft + resLen - 1]
