# 预处理组合数 C(n,k)=C(n-1,k)+C(n-1,k-1)
comb = [[0] * 36 for _ in range(36)]
for i in range(36):
    comb[i][0] = 1
    for j in range(1, i + 1):
        comb[i][j] = comb[i - 1][j - 1] + comb[i - 1][j]

print(comb[10][2])

