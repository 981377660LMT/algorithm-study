# https://blog.csdn.net/qq_35577488/article/details/114221767
# 循环出现的时间段 - 中国剩余定理
# 有两个循环出现的时段，第一个范围为：
# [k(2x+2y)+x,k(2x+2y)+x+y),k≥0
# 第二个范围是[c(P+Q)+P,c(P+Q)+P+Q).
# 现在问你是否存在一个时刻使得同时属于两个时段?求出最小的时刻.
# x,P<=1e9
# y,Q<=500


# !容易发现两个时段的长度最大为500.暴力枚举两个时段里的点是可以接受的.
# 假设这两个时刻为i,j(相对线段位置来说)。那么我们假设这个时刻为t.
# 则有同余方程组
# t≡i(mod 2x+2y)
# t≡j(mod P+Q)

from 中国剩余定理 import excrt

INF = int(1e20)


def overSleeping(X: int, Y: int, P: int, Q: int) -> int:
    res = INF
    for i in range(X, X + Y):
        for j in range(P, P + Q):
            cand = excrt([1, 1], [i, j], [2 * X + 2 * Y, P + Q])
            if cand is not None:
                res = min(res, cand[0])
    return res


if __name__ == "__main__":
    T = int(input())
    for _ in range(T):
        X, Y, P, Q = map(int, input().split())
        res = overSleeping(X, Y, P, Q)
        if res == INF:
            print("infinity")
        else:
            print(res)
