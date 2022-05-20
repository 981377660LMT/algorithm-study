# 当 N 为奇数时，我们可以通过以下方法构建一个幻方：

# 首先将 1 写在第一行的中间。

# 之后，按如下方式从小到大依次填写每个数 K(K=2,3,…,N×N)：

# 1. 若 (K−1) 在第一行但不在最后一列，则将 K 填在最后一行，(K−1) 所在列的右一列；
# 2. 若 (K−1) 在最后一列但不在第一行，则将 K 填在第一列，(K−1) 所在行的上一行；
# 3. 若 (K−1) 在第一行最后一列，则将 K 填在 (K−1) 的正下方；
# 4. 若 (K−1) 既不在第一行，也不在最后一列，如果 (K−1) 的右上方还未填数，则将 K 填在 (K−1) 的右上方，否则将 K 填在 (K−1) 的正下方。

# 现给定 N，请按上述方法构造 N×N 的幻方。

# 1≤N≤39,N 为奇数。
n = int(input())

matrix = [[0] * n for _ in range(n)]
matrix[0][n // 2] = 1

row, col = 0, n // 2
for k in range(2, n * n + 1):
    if row == 0 and col != n - 1:
        row, col = n - 1, col + 1
    elif row != 0 and col == n - 1:
        row, col = row - 1, 0
    elif row == 0 and col == n - 1:
        row, col = row + 1, col
    else:
        if matrix[row - 1][col + 1] == 0:
            row, col = row - 1, col + 1
        else:
            row, col = row + 1, col
    matrix[row][col] = k


for row in matrix:
    print(' '.join(map(str, row)))

