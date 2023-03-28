# https://yukicoder.me/problems/no/205


from heapq import heapify, heappop, heappush
from typing import List


def minLexMerge(words: List[str]) -> str:
    """字典序最小的合并字符串"""
    pq = [w + chr(200) for w in words]
    heapify(pq)
    res = []
    while pq:
        min_ = heappop(pq)
        res.append(min_[0])
        min_ = min_[1:]
        if len(min_) >= 2:
            heappush(pq, min_)
    return "".join(res)


if __name__ == "__main__":
    N = int(input())
    words = [input() for _ in range(N)]
    print(minLexMerge(words))
