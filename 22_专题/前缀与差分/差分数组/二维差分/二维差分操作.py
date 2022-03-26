# 一个n行m列的字母矩阵，矩阵中仅包含大小写字母。
# 有q次操作，每次操作选择一个子矩阵，将其中所有字母大小写转换(大写变小写，小写变大写)。
# 请你输出所有操作结束后的矩阵。


class DiffMatrix:
    def __init__(self, A):
        self.m, self.n = len(A), len(A[0])
        self.matrix = [[0] * self.n for _ in range(self.m)]
        for i in range(self.m):
            for j in range(self.n):
                self.matrix[i][j] = A[i][j]
        self.diff = [[0] * (self.n + 2) for _ in range(self.m + 2)]

    def add(self, r1, c1, r2, c2, k):
        self.diff[r1 + 1][c1 + 1] += k
        self.diff[r1 + 1][c2 + 2] -= k
        self.diff[r2 + 2][c1 + 1] -= k
        self.diff[r2 + 2][c2 + 2] += k

    def update(self):
        for i in range(self.m):
            for j in range(self.n):
                self.diff[i + 1][j + 1] += (
                    self.diff[i + 1][j] + self.diff[i][j + 1] - self.diff[i][j]
                )
                self.matrix[i][j] += self.diff[i + 1][j + 1]


n, m, q = map(int, input().split())
matrix = []
for _ in range(n):
    matrix.append(list(input()))

diff = DiffMatrix([[0] * m for _ in range(n)])
for _ in range(q):
    x1, y1, x2, y2 = map(int, input().split())
    x1, y1, x2, y2 = x1 - 1, y1 - 1, x2 - 1, y2 - 1
    diff.add(x1, y1, x2, y2, 1)
diff.update()

newMatrix = diff.matrix
for i in range(n):
    for j in range(m):
        if newMatrix[i][j] & 1:
            raw = ord(matrix[i][j])
            matrix[i][j] = chr(raw ^ 32)

for row in matrix:
    print(''.join(row))

