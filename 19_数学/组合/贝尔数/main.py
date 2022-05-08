MAX = 2010
bell = [[0] * MAX for _ in range(MAX)]
bell[1][1] = 1
for i in range(2, MAX):
    bell[i][1] = bell[i - 1][i - 1]
    for j in range(2, i + 1):
        bell[i][j] = bell[i - 1][j - 1] + bell[i][j - 1]


def cal(n):
    return bell[n][n]


print(cal(2))


