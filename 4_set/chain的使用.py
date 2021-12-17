# 使用 chain() 的一个常见场景是当你想对不同的集合中所有元素执行某些操作的时候。
# 他屏蔽了很多细节，统一了遍历接口
from itertools import chain

aset = set([1, 2, 3])
blis = [12, 3]
cdict = {1: 2, 4: 9}
for key in chain(aset, blis, cdict):
    print(key)


# 还有一个好处是统一遍历矩阵行列
matrix = [[0, 1], [1, 0], [2, 3]]
for line in chain(matrix, zip(*matrix)):
    print(list(line))
