# https://codeforces.com/problemset/problem/712/C
# 给你两个整数 x 和 y，满足 3<=y<x<=1e5。
# 从边长为 x 的等边三角形出发，每次操作你可以将其中一条边的长度修改为某个整数，要求修改后的三条边仍能组成一个三角形。
# 将三角形修改成边长为 y 的等边三角形，最少需要操作多少次？

# https://codeforces.com/contest/712/submission/158586227
# 正难则反。
# 从边长为 y 的等边三角形出发，每次将最短的边修改为另外两条边的和减一，直到最短的边不低于 x。
x, y = map(int, input().split())
a, b, c = y, y, y
res = 0
while True:
    a, b, c = sorted([a, b, c])
    if a >= x:
        break
    a = b + c - 1
    res += 1
print(res)
