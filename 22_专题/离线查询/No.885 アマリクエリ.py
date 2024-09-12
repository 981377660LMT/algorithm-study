# No.885 アマリクエリ(modSum)
# https://yukicoder.me/problems/no/885
# 给定一个数组和q个查询，每个查询形如Xi，将每个数模Xi后，求数组的和。
# 查询之间互相影响。


from heapq import heappop, heappush
import sys


input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    Q = int(input())
    X = list(map(int, input().split()))

    res = sum(A)
    pq = []  # 大根堆
    for v in A:
        heappush(pq, -v)

    for v in X:
        while True:
            cur = -heappop(pq)
            if cur < v:
                heappush(pq, -cur)
                break
            else:
                mod = cur % v
                res -= cur - mod
                heappush(pq, -mod)

        print(res)
