# https://techblog.finatext.com/lets-do-competitive-programming-with-python-9c8b834769f6
# https://qiita.com/c-yan/items/dbf2838cdd89864ef5ac
# https://qiita.com/shoji9x9/items/e7d19bd6f54e960f46be

from operator import itemgetter


from random import random
import time
from timeit import timeit

####################################################################################
# !sys.stdin.readline() 比 input 快很多
####################################################################################
# list的随机访问很慢，尽量避免
# !List の効率的な使い方
# !index によるアクセスは遅い

n = 10**7
a = list(range(n))


# 0.5065463000000818
print(timeit("for i in range(n): a[i]", number=1, globals=globals()))

# 0.11476179999954184
print(timeit("for ai in a: ai", number=1, globals=globals()))

# 0.37129209999966406
print(timeit("for i, ai in enumerate(a): ai", number=1, globals=globals()))
####################################################################################
# !二重循环里缓存数组
# for i in range(Y):
#   for j in range(X):
#     a[i + 1][j] = ... a[i][j] ...

# ↓

# for i in range(Y):
#   ai = a[i]
#   ai1 = a[i + 1]
#   for j in range(X):
#     ai1[j] = ... ai[j] ...
####################################################################################
# !使用切片代替循环赋值
# for x in range(i, j):
#   a[x] = y

# ↓

# a[i:j] = [y] * (j - i)
####################################################################################
# for j in range(1 << n):
#     if dp[j] + a < dp[j | c]:
#         dp[j | c] = dp[j] + a
#     # dp[j | c] = min(dp[j] + a, dp[j | c])
# 15,16行目(if文)	  => 758ms
# 17行目(毎回更新)	=> 1660ms

# !Listへの不必要な書き込みがないか注意しながら書いてみましょう。
####################################################################################
# 利用itemgetter代替key排序会快一些
# !ListのListのソートにはoperator.itemgetter
n = 10**6
a = [[random(), random(), random()] for _ in range(n)]
# 0.36658799999986513
print(timeit("a.sort(key=lambda x: x[1])", number=1, globals=globals()))
# 0.2570695999993404
print(timeit("a.sort(key=itemgetter(1))", number=1, globals=globals()))
####################################################################################
# local変数を利用した高速化
# 把全局变量转换为局部变量，可以提高运行速度???
# !在python中用main函数会快一些，但是在pypy中会慢一些.
####################################################################################
# !pypy比python慢的场合(PyPyだと遅い処理)
# - 递归
# - 字符串处理
# - tuple的比较
####################################################################################
# !数値演算より if 文を使う
# 与其他语言不同，python的数值运算很慢
# !因此使用if语句代替数值运算可以提高速度
# x += speed * direction

# ↓

# if direction == 1:
#   x += speed
# else:
#   x -= speed
####################################################################################
# !min, max より if 文を使う
# 関数呼び出しは重いので、if 文で書くほうが速くなる.


# dp[i + a] = max(dp[i + a], dp[i] + b)

# ↓

# if dp[i] + b > dp[i + a]:
#   dp[i + a] = dp[i] + b
####################################################################################
# !遍历list比遍历range快
# くり返し同じループをするときはリスト(配列)を使う
# くり返し同じループをするときは、range をリスト(配列)に変換すると速くなる.
# range は1ループ毎に long オブジェクトを生成するため malloc をするのだが、
# リスト(配列)にしてしまえばそのコストはループ一周分だけで済む.
# for i in range(Y):
#   for j in range(X):
#     ...

# ↓

# listX = list(range(X))
# for i in range(Y):
#   for j in listX:
#     ...
####################################################################################
# !插入时使用array会快一些
import bisect, array

# 遅い（というか普通はこう書く） -> 451 ms
time1 = time.time()
q = list(range(0, 10**5, 2))
for i in range(1, 10**5, 2):
    bisect.insort_left(q, i)
print(time.time() - time1)

# 速い -> 142 ms
time1 = time.time()
q = array.array("i", list(range(0, 10**5, 2)))
for i in range(1, 10**5, 2):
    bisect.insort_left(q, i)
print(time.time() - time1)
####################################################################################
# !尽量减少循环内的判断/使用批量处理
# !ループ内の処理はシンプルに保つ（if文での判定は行わずループ後にうまく処理するなど）
# # ループの中でグループ数を数え判定しようとしたところTLEになった
# for c, a, b in G:
#     if uft.group_count() <= K:
#         break
#     if not uft.same(a, b):
#         uft.union(a, b)
#         ans += c

# # ループ後に必要な数値のみ切り出し答えとしたところACした
# for c, a, b in G:
#     if not uft.same(a, b):
#         uft.union(a, b)
#         cost.append(c)
# ans = sum(cost[:N-K])
####################################################################################
####################################################################################
