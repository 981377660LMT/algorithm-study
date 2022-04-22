# https://docs.scipy.org/doc/scipy/reference/generated/scipy.sparse.csgraph.maximum_flow.html
from scipy.sparse import csr_matrix
from scipy.sparse.csgraph import maximum_flow

# CSR - 压缩稀疏行（Compressed Sparse Row），按行压缩。
adjMatrix = csr_matrix(
    [
        [0, 16, 13, 0, 0, 0],
        [0, 0, 10, 12, 0, 0],
        [0, 4, 0, 0, 14, 0],
        [0, 0, 9, 0, 0, 20],
        [0, 0, 0, 7, 0, 4],
        [0, 0, 0, 0, 0, 0],
    ]
)
print(adjMatrix)
# 稀疏矩阵

print(maximum_flow(adjMatrix, 0, 5).flow_value)
