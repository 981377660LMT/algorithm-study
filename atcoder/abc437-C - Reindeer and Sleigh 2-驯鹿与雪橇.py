# https://atcoder.jp/contests/abc437/editorial/14862
# 有N个驯鹿和1个雪橇。  第i只驯鹿体重 Wi ​ 和力量 Pi ​ 。
# 对于每只驯鹿，选择“拉雪橇”或“乘坐雪橇”。
# 然而，驯鹿拉雪橇时的力总和必须大于或等于乘坐雪橇的驯鹿重量之和。
# 雪橇最多能骑多少只驯鹿？

T = int(input())
for _ in range(T):
    N = int(input())
    W, P = [0] * N, [0] * N
    for i in range(N):
        W[i], P[i] = map(int, input().split())

    order = sorted(range(N), key=lambda x: P[x] + W[x])
    sumP = sum(P)
    res = N
    for c, i in enumerate(order, 1):
        w, p = W[i], P[i]
        res += w + p
        if res > sumP:
            res = c - 1
            break
    print(res)
