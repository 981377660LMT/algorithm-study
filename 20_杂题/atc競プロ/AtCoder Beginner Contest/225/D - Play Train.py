"""火车车厢模拟"""

# 题意:给定数字1,2,3,.. .,N,需要进行一下操作
# 1. 以x结尾的组和以y开头的组合并(前后顺序不变),
# 2. x,y所在的组分开(x, y中间是断点)
# 2. 查询x所在的组的所有人(按顺序输出)

# 我们只需要维护两个数组即可,记为pre, next
# pre和next数组模拟链表


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":
    n, q = map(int, input().split())
    pre, next = [-1] * n, [-1] * n
    for _ in range(q):
        t, *rest = map(int, input().split())
        if t == 1:
            x, y = rest
            x, y = x - 1, y - 1
        elif t == 2:
            x, y = rest
            x, y = x - 1, y - 1
        else:
            x = rest[0] - 1
