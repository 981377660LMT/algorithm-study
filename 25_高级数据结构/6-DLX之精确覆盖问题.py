# 能否在矩阵中找到一个行的集合，使得这些行中，每一列都有且仅有一个数字 1。
# n,m<=500
# 数据保证矩阵中 1 的数量不超过 5000。 (稀疏矩阵)
# https://www.acwing.com/problem/content/1069/
row, col = map(int, input().split())
matrix = []
for _ in range(row):
    matrix.append(list(map(int, input().split())))

# todo
