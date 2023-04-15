from collections import defaultdict, deque
from heapq import heappop, heappush
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 文字列
# S があり、初め
# S= 1 です。
# 以下の形式のクエリが
# Q 個与えられるので順に処理してください。

# 1 x :
# S の末尾に数字
# x を追加する
# 2 :
# S の先頭の数字を削除する
# 3 :
# S を十進数表記の数とみなした値を
# 998244353 で割った余りを出力する
# 制約
# 1≤Q≤6×10
# 5

# 1 番目の形式のクエリについて、
# x∈{1,2,3,4,5,6,7,8,9}
# 2 番目の形式のクエリは
# S が
# 2 文字以上の時にのみ与えられる
# 3 番目の形式のクエリが
# 1 個以上存在する


if __name__ == "__main__":
    q = int(input())
    len_ = 1
    cur = 1
    queue = deque([1])
    for _ in range(q):
        query = list(map(int, input().split()))
        if query[0] == 1:
            x = query[1]
            len_ += 1
            queue.append(x)
            cur = (cur * 10 + x) % MOD
        elif query[0] == 2:
            y = queue.popleft()
            cur = cur - y * pow(10, len_ - 1, MOD)
            len_ -= 1
            cur %= MOD
        else:
            print(cur % MOD)
