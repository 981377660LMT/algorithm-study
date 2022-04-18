# Number of rectangles in N*M grid
# 矩阵中的矩形数量


def count(n, m):
    return (m * n * (n + 1) * (m + 1)) // 4


# Driver code
n, m = 5, 4
print(count(n, m))
