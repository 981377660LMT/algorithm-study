from heapq import heapify, heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# AtCoder 王国では、
# N 種類のたこ焼きが売られています。
# i 種類目のたこ焼きの値段は
# A
# i
# ​
#   円です。

# 高橋君は、合計で
# 1 個以上のたこ焼きを買います。このとき、同じたこ焼きを複数個買うことも許されます。

# 高橋君が支払う金額としてあり得るもののうち、安い方から
# K 番目の金額を求めてください。ただし、同じ金額を支払う方法が複数存在する場合は
# 1 回だけ数えます。
if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))
    visited = set()
    pq = nums[:]
    heapify(pq)
    while True:
        num = heappop(pq)
        if num in visited:
            continue
        visited.add(num)
        if len(visited) == k:
            print(num)
            exit(0)
        for i in range(n):
            next_ = num + nums[i]
            if next_ in visited:
                continue
            heappush(pq, next_)
