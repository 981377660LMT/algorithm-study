# len(s)<=300
# 可以不断删除'iwi' 求操作的最大步数
# iwiwii 先に 3 文字目から 5 文字目の iwi を取り除くのが最適である。

# !注意到消除'iwii'或者'iiwi'之后，对左右两边没有影响 可以先消除这两种
# !最后所有的iwi全部为 'wiwiw'
# !之后再消除'iwi'
import sys


sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)


s = input()
n = len(s)

res = 0

stack1, need1 = [], set(["iwii", "iiwi"])
for char in s:
    stack1.append(char)
    while "".join(stack1[-4:]) in need1:
        for _ in range(4):
            stack1.pop()
        stack1.append("i")
        res += 1

stack2 = []
for char in stack1:
    stack2.append(char)
    while "".join(stack2[-3:]) == "iwi":
        for _ in range(3):
            stack2.pop()
        res += 1


print(res)
