"""火车车厢模拟链表"""

# 题意:给定数字1,2,3,.. .,N,需要进行一下操作

# 1. 以x结尾的组和以y开头的组合并(前后顺序不变),
# !pre[y] = x, next[x] = y
# 2. x,y所在的组分开(x, y中间是断点)
# !pre[y]=-1, next[x]=-1
# 3. 查询x所在的组的所有人(按顺序输出)

# 我们只需要维护两个数组即可,记为pre, next
# !pre和next数组模拟链表


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
        if t == 1:  # !merge
            x, y = rest
            x, y = x - 1, y - 1
            pre[y] = x
            next[x] = y
        elif t == 2:  # !split
            x, y = rest
            x, y = x - 1, y - 1
            pre[y] = -1
            next[x] = -1
        else:  # !query
            x = rest[0] - 1

            group1 = []
            cur = x
            while pre[cur] != -1:
                cur = pre[cur]
                group1.append(cur)

            group2 = []
            cur = x
            while next[cur] != -1:
                cur = next[cur]
                group2.append(cur)

            res = group1[::-1] + [x] + group2
            print(len(res), *[i + 1 for i in res])
