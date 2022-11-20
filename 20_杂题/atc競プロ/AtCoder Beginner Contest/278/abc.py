from collections import defaultdict, deque
import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)
# 高橋君が運営する SNS「Twidai」にはユーザー 1 からユーザー N までの N 人のユーザーがいます。 Twidai では、ユーザーは別のユーザーをフォローすることや、フォローを解除することができます。

# Twidai がサービスを開始してから、Q 回の操作が行われました。 i 回目 (1≤i≤Q) の操作は 3 つの整数 T
# i
# ​
#  ,A
# i
# ​
#  ,B
# i
# ​
#   で表され、それぞれ次のような操作を表します。

# T
# i
# ​
#  =1 のとき：ユーザー A
# i
# ​
#   がユーザー B
# i
# ​
#   をフォローしたことを表す。この操作の時点でユーザー A
# i
# ​
#   がユーザー B
# i
# ​
#   をフォローしている場合、ユーザーのフォロー状況に変化はない。
# T
# i
# ​
#  =2 のとき：ユーザー A
# i
# ​
#   がユーザー B
# i
# ​
#   のフォローを解除したことを表す。この操作の時点でユーザー A
# i
# ​
#   がユーザー B
# i
# ​
#   をフォローしていない場合、ユーザーのフォロー状況に変化はない。
# T
# i
# ​
#  =3 のとき：ユーザー A
# i
# ​
#   とユーザー B
# i
# ​
#   が互いにフォローしているかをチェックすることを表す。この操作の時点でユーザー A
# i
# ​
#   がユーザー B
# i
# ​
#   をフォローしており、かつユーザー B
# i
# ​
#   がユーザー A
# i
# ​
#   をフォローしているとき、このチェックに対して Yes と答え、そうでないときこのチェックに対して No と答える必要がある。
# サービス開始時には、どのユーザーも他のユーザーをフォローしていません。

# すべての T
# i
# ​
#  =3 であるような操作に対して、i が小さいほうから順番に正しい答えを出力してください。

if __name__ == "__main__":
    n, q = map(int, input().split())
    follow = defaultdict(set)
    for _ in range(q):
        kind, a, b = map(int, input().split())
        if kind == 1:
            follow[a].add(b)
        elif kind == 2:
            follow[a].discard(b)
        else:
            if b in follow[a] and a in follow[b]:
                print("Yes")
            else:
                print("No")
