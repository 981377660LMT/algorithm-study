from scipy.spatial import KDTree

points = [(2, 3), (-2, 3), (2, -3), (1, -1)]
kdtree = KDTree(points)

# 查找到 (1,1) 的 KNN，返回 []
# 用于并行处理的工作人员数量。如果给定 -1，则使用所有 CPU 线程
res = kdtree.query((1, 1), k=1, workers=-1)
print(res)
# # 查找点x的距离r内的所有点
# print(kdtree.query_ball_point((1, 1), 4, workers=-1))
# print(kdtree.query_ball_point((1, 1), 4, workers=-1, return_length=True))

# # 查找距离内的所有对点
# print(kdtree.query_pairs(r=3))

