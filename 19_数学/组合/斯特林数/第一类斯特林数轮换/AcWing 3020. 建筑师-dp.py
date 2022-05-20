# 1—n个建筑物让他们排列起来，左边与右边分别可以看见A,B个建筑物，问建筑物排列的方案数？
# n<=5e4 A,B<=100
# https://www.acwing.com/solution/content/47769/

# 将最高点n当作分割点 左边看到 A-1个建筑物 右边B-1个建筑物
# 每个小组对应一个圆排列
# 即从n-1个数中选出(A+B-2)个圆排列 再选(A-1)个放在左边
MOD = int(1e9 + 7)

ROW, COL = 50005, 200
dp1 = [[0] * COL for _ in range(ROW)]
dp1[0][0] = 1
for i in range(1, ROW):
    for j in range(1, COL):
        dp1[i][j] = (dp1[i - 1][j - 1] + (i - 1) * dp1[i - 1][j]) % MOD

C = [[0] * 200 for _ in range(200)]
for i in range(200):
    C[i][0] = 1
    for j in range(1, i + 1):
        C[i][j] = C[i - 1][j - 1] + C[i - 1][j]


def main(n: int, A: int, B: int) -> int:
    stirling1 = dp1[n - 1][A + B - 2]
    comb = C[A + B - 2][A - 1]
    return stirling1 * comb % MOD


T = int(input())
for _ in range(T):
    n, A, B = map(int, input().split())
    print(main(n, A, B))

