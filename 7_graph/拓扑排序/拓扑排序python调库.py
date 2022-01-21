# https://docs.python.org/zh-cn/3.11/library/graphlib.html#graphlib.TopologicalSorter
from collections import defaultdict
from graphlib import TopologicalSorter

adjMap = defaultdict(list)

adjMap['D'].extend(['B', 'C'])
adjMap['C'].extend(['A'])
adjMap['B'].extend(['A'])

ts = TopologicalSorter(adjMap)
ts.add('A', 'B')
print(*ts.static_order())


# 在通常情况下，对给定的图执行排序所需的步骤如下:
# 通过可选的初始图创建一个 TopologicalSorter 的实例。
# 添加额外的节点到图中。
# 在图上调用 prepare()。
# 当 is_active() 为 True 时，迭代 get_ready() 所返回的节点并加以处理。 完成处理后在每个节点上调用 done()。

# 在只需要对图中的节点进行立即排序并且不涉及并行性的情况下，可以直接使用便捷方法TopologicalSorter.static_order():
