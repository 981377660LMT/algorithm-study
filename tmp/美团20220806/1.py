# 美团笔试1
# !正解是线性规划 求直线两两的交点处的最大值

# 2*x+y<=A
# x+2*y<=B
# 求x+y的最大值
# !可以直接min(A,B,(A+B)//3)

import sys

sys.setrecursionlimit(int(1e6))
input = lambda: sys.stdin.readline().rstrip("\r\n")


T = int(input())
for _ in range(T):
    x, y = map(int, input().split())  # x个A点心 y个B点心
    upper = (x + y) // 3
    print(min(x, y, upper))
