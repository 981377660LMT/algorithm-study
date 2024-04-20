import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

# 英大小文字と ( 、 ) からなる文字列
# S=S
# 1
# ​
#  S
# 2
# ​
#  S
# 3
# ​
#  …S
# ∣S∣
# ​
#   が与えられます。
# 文字列
# S 中の括弧は、対応が取れています。

# 次の操作を、操作ができなくなるまで繰り返します。

# まず、以下の条件を全て満たす整数組
# (l,r) をひとつ選択する。
# l<r
# S
# l
# ​
#  = (
# S
# r
# ​
#  = )
# S
# l+1
# ​
#  ,S
# l+2
# ​
#  ,…,S
# r−1
# ​
#   は全て英大文字または英小文字である
# T=
# S
# r−1
# ​
#  S
# r−2
# ​
#  …S
# l+1
# ​

# ​
#   とする。
# 但し、
# x
#   は
# x の大文字と小文字を反転させた文字列を指す。
# その後、
# S の
# l 文字目から
# r 文字目までを削除し、削除した位置に
# T を挿入する。
# 詳細は入出力例を参照してください。

# 上記の操作を使って全ての ( と ) を除去することができ、最終的な文字列は操作の方法や順序によらないことが証明できます。
# このとき、最終的な文字列を求めてください。


if __name__ == "__main__":
    s = input()
    ords = [ord(c) for c in s]
    n = len(s)
    swapPair = [0] * n
    stack = []
    for i, v in enumerate(s):
        if v == "(":
            stack.append(i)
        elif v == ")":
            j = stack.pop()
            swapPair[i] = j
            swapPair[j] = i

    sb = []
    index, direction, flag = 0, 1, False
    while index < n:
        if s[index] == "(" or s[index] == ")":
            index = swapPair[index]
            direction = -direction
            flag = not flag
        else:
            sb.append(s[index] if not flag else chr(ords[index] ^ 32))
        index += direction
    print("".join(sb))
