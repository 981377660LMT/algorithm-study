# F - Smooth Occlusion (磨牙)
# https://atcoder.jp/contests/abc395/tasks/abc395_f
# 高橋君有 2N 颗牙齿，其中 N 颗为上牙，N 颗为下牙。第 i 颗上牙的长度为 U_i，第 i 颗下牙的长度为 D_i。
#
# 他的牙齿「咬合良好」需要同时满足两个条件：
#
# !存在一个整数 H，使得对于所有 1 ≤ i ≤ N，都有 U_i + D_i = H。
# !对于所有 1 ≤ i < N，都有 |U_i − U_{i+1}| ≤ X。
# !高橋君可以反复执行如下操作：花 1 日元使用磨牙机，将任一牙齿（要求牙长始终为正）长度减 1。只能通过这种操作改变牙齿长度。
#
# 问题要求求出使牙齿「咬合良好」所需花费的最小金额。


INF = int(4e18)


def min2(a: int, b: int) -> int:
    return a if a < b else b


if __name__ == "__main__":
    N, X = map(int, input().split())
    U, D = [0] * N, [0] * N
    for i in range(N):
        U[i], D[i] = map(int, input().split())

    minU, minD, minH = INF, INF, INF
    for u, d in zip(U, D):
        minU = min2(minU + X, u)
        minD = min2(minD + X, d)
        minH = min2(minH, minU + minD)

    print(sum(U) + sum(D) - N * minH)
