# 二十四小时进制下，在S时0分开灯，在T时0分关灯，问X时30分是否灯亮?


import sys

sys.setrecursionlimit(int(1e9))
input = lambda: sys.stdin.readline().rstrip("\r\n")
MOD = 998244353
INF = int(4e18)

if __name__ == "__main__":

    def toMinute(h: int, m: int) -> int:
        return h * 60 + m

    s, t, x = map(int, input().split())
    if t < s:
        t += 24
    m1 = toMinute(s, 0)
    m2 = toMinute(t, 0)
    m3 = toMinute(x, 30)
    m4 = toMinute(x + 24, 30)
    if m1 <= m3 < m2 or m1 <= m4 < m2:
        print("Yes")
    else:
        print("No")
