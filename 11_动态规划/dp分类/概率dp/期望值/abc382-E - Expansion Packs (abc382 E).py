# E - Expansion Packs
# https://atcoder.jp/contests/abc382/tasks/abc382_e
# 一个卡牌包里有n张卡牌，第i张卡牌是稀有牌的概率是pi。
# 现在将一包一包的开牌，直到开出X张稀有牌为止，问开的卡牌包数的期望是多少
# N,X<=5000,1<=Pi<=100
#
# dp[i][j]为一包中拆前i张牌，有j张是稀有牌的概率


if __name__ == "__main__":
    N, X = map(int, input().split())
    P = list(map(int, input().split()))

    p = [pi / 100 for pi in P]
    q = [0.0] * (N + 1)
    q[0] = 1.0

    for i in range(N):
        pi = p[i]
        for k in range(i + 1, 0, -1):
            q[k] = q[k] * (1 - pi) + q[k - 1] * pi
        q[0] *= 1 - pi

    f = [0.0] * (X + 1)
    f[0] = 0.0
    for s in range(1, X + 1):
        tmp_sum = 1.0
        for k in range(1, min(s, N) + 1):
            tmp_sum += q[k] * f[s - k]
        f[s] = tmp_sum / (1 - q[0])

    print("%.16f" % f[X])
