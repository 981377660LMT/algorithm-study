# 给定一个如下图所示的数字三角形，从顶部出发，
# 在每一结点可以选择移动至其左下方的结点或移动至其右下方的结点，一直走到底层，
# 要求找出一条路径，使路径上的数字的和最大。
#        7
#       3   8
#     8   1   0
#   2   7   4   4
# 4   5   2   6   5
# 自下往上做更方便

n = int(input())
mat = [[0]]
for i in range(n):
    row = list(map(int, input().split()))
    mat.append(row)

# 多开了一些dp的空间，这样可以简化初始化的问题
dp = [[0] * (n + 1) for _ in range(n + 2)]
for i in range(n, 0, -1):
    for j in range(i):
        dp[i][j] = max(dp[i + 1][j], dp[i + 1][j + 1]) + mat[i][j]

print(dp[1][0])

