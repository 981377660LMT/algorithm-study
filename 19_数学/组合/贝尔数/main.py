MOD = 998244353

N = 400
bell = [[0] * N for _ in range(N)]
bell[1][1] = 1
for i in range(2, N):
    bell[i][1] = bell[i - 1][i - 1]
    for j in range(2, i + 1):
        bell[i][j] = bell[i - 1][j - 1] + bell[i][j - 1]
        bell[i][j] %= MOD


def cal(n):
    return bell[n][n]


print(cal(5))
# B 0 = 1, B 1 = 1, B 2 = 2, B 3 = 5, B 4 = 15, B 5 = 52, B 6 = 203
