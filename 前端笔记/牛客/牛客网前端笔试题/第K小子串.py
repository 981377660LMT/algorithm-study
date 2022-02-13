# 1<=k<=5 很小，适合用堆而不是直接sort
from heapq import heappush, heappop


s = input()
k = int(input())


class WrappedStr(str):
    def __init__(self, string: str):
        super().__init__()
        self.string = string

    def __lt__(self, other: 'WrappedStr') -> bool:
        return self.string > other.string

    def __eq__(self, other: 'WrappedStr') -> bool:
        return self.string == other.string


pq: list[WrappedStr] = []
visited = set()

for length in range(1, k + 1):
    for i in range(len(s) - length + 1):
        if s[i : i + length] in visited:
            continue
        visited.add(s[i : i + length])
        heappush(pq, WrappedStr(s[i : i + length]))
        if len(pq) > k:
            heappop(pq)

print(pq[0].string)

