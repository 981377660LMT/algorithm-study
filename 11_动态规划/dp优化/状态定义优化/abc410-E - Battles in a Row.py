# https://atcoder.jp/contests/abc410/tasks/abc410_e
# 高桥君依次与 N 个怪物战斗，初始体力 H、魔力 M。第 i 个怪物可选：
#
# 用体力战斗：前提体力 ≥ A[i]，消耗体力 A[i]；
# 用魔力战斗：前提魔力 ≥ B[i]，消耗魔力 B[i]。
# 只要面对某个怪物时两种方式都不满足，游戏结束。求最多能连续击败多少只怪物。
#
# N,H,M<=3000
#
# !dp[h] 表示“已消耗体力 h 时，剩余魔力最多能是多少”。

if __name__ == "__main__":
    N, H, M = map(int, input().split())
    A, B = [0] * N, [0] * N
    for i in range(N):
        A[i], B[i] = map(int, input().split())

    dp = [M] * (H + 1)
    for i in range(N):
        ndp = [-1] * (H + 1)
        a, b = A[i], B[i]
        for x in range(H + 1):
            if x >= a:
                ndp[x - a] = max(ndp[x - a], dp[x])
            if dp[x] >= b:
                ndp[x] = max(ndp[x], dp[x] - b)
        dp = ndp
        if max(dp) < 0:
            print(i)
            exit()

    print(N)
