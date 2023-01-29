import numpy as np
from scipy.spatial import KDTree

points = [(2, 3), (-2, 3), (2, -3), (1, -1)]
kdtree = KDTree(points)  # 构建 KD 树,默认距离度量是欧氏距离

# 查找到 (1,1) 的 KNN，返回 []
# workers:用于并行处理的工作人员数量。如果给定 -1，则使用所有 CPU 线程
# p : float, 1<=p<=infinity, 使用哪个Minkowski p-norm
# 1:曼哈顿距离 2:欧式距离 无穷:最大坐标差距离
res = kdtree.query(x=(1, 1), k=1, p=1, workers=-1, distance_upper_bound=np.inf)
print(res)

# # 查找点x的距离r内的所有点
print(kdtree.query_ball_point((1, 1), 4, workers=-1))

# # 查找距离内的所有对点
print(kdtree.query_pairs(r=3))


# 查找最近点,kdtree是O(logn)的
# 但是查找knn,复杂度最坏会退化到O(n)
