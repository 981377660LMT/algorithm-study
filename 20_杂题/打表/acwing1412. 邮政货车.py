n = int(input())
n = n - 1
f = [0, 1, 2, 6]
for i in range(4, n + 1):
    f.append(f[i - 1] * 2 + f[i - 2] * 2 - f[i - 3] * 2 + f[i - 4])
print(f[n] * 2)
print(f)
# https://oeis.org/

# Number of Hamiltonian cycles in P_4 X P_n.
# https://oeis.org/search?q=0%2C+1%2C+2%2C+6%2C+14%2C+37%2C+92%2C+236%2C+596%2C+1517&language=english&go=Search

# 暴力打表前几项后去oeis上搜索即可
