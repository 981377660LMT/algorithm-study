# 一共会进行Q次操作，每次操作有三种类型：

# 1.在一个空球上面写一个整数X,并把这个球放入包内。
# !2.对于包内的所有球，将每个球上面的整数加上X 。
# 3.输出包中所有球上的最小的数字，并把这个球扔掉。

# !在每个数字加入到堆中之前，先减去累加值，取出的时候再加上新的累加值
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    q = int(input())
    offset = 0
    pq = []
    for _ in range(q):
        kind, *rest = map(int, input().split())
        if kind == 1:
            num = rest[0]
            heappush(pq, num - offset)
        elif kind == 2:
            delta = rest[0]
            offset += delta
        else:
            top = heappop(pq)
            print(top + offset)
