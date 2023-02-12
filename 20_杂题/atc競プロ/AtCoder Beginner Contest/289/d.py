from collections import deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 無限に続く階段があります。 一番下は
# 0 段目で、
# 1 段のぼるごとに
# 1 段目、
# 2 段目と続きます。

# 0 段目に階段登りロボットがいます。 階段登りロボットは、一回の動作で
# A
# 1
# ​
#  ,A
# 2
# ​
#  ,…,A
# N
# ​
#   段ぶん階段をのぼることができます。 つまり、階段登りロボットが
# i 段目にいるとき、一回動作をした後は
# i+A
# 1
# ​
#   段目、
# i+A
# 2
# ​
#   段目、⋯、
# i+A
# N
# ​
#   段目のいずれかにいることができます。 それ以外の段数を一回の動作でのぼることはできません。 階段登りロボットは階段を下ることもできません。

# 階段の
# B
# 1
# ​
#  ,B
# 2
# ​
#  ,…,B
# M
# ​
#   段目にはモチが設置されています。 モチが設置されている段へのぼるとロボットは動けなくなり、他の段に移動することができなくなります。

# 階段登りロボットは階段のちょうど
# X 段目にのぼりたいです。 階段登りロボットが階段のちょうど
# X 段目にのぼることが可能か判定してください。
if __name__ == "__main__":
    n = int(input())
    A = list(map(int, input().split()))
    m = int(input())
    B = set(map(int, input().split()))
    X = int(input())

    visited = set()
    queue = deque([0])
    while queue:
        cur = queue.popleft()
        if cur == X:
            print("Yes")
            exit()
        for a in A:
            nxt = cur + a
            if nxt in B or nxt in visited or nxt > X:
                continue
            visited.add(nxt)
            queue.append(nxt)

    print("No")
