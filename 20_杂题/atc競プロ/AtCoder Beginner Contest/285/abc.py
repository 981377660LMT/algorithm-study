from collections import defaultdict
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# あなたの運営する Web サービスには
# N 人のユーザがいます。

# i 番目のユーザの現在のユーザ名は
# S
# i
# ​
#   ですが、
# T
# i
# ​
#   への変更を希望しています。
# ここで、
# S
# 1
# ​
#  ,…,S
# N
# ​
#   は相異なり、
# T
# 1
# ​
#  ,…,T
# N
# ​
#   も相異なります。

# ユーザ名を変更する順序を適切に定めることで、以下の条件を全て満たすように、全てのユーザのユーザ名を希望通り変更することができるか判定してください。

# ユーザ名の変更は
# 1 人ずつ行う
# どのユーザもユーザ名の変更は一度だけ行う
# ユーザ名の変更を試みる時点で他のユーザが使っているユーザ名に変更することはできない


from typing import List
from collections import deque


# TODO adjMap版的拓扑排序
if __name__ == "__main__":
    a, b = map(int, input().split())
    print("Yes" if a == b // 2 else "No")
