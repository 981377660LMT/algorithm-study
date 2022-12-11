# 货船中心装货物
# 每个查询[left,right] 禁用下标[left,right]的箱子 每个箱子只能放一个货物
# !问每个查询中所有箱子可以装的货物的最大价值
# !n,m,q<=50

# !每次查询O(nlogm)
# 1. 货物按照价值排序,先放价值大的货物
# 2. 箱子按照容量升序，每次找到第一个能放下货物的箱子 (bisect_left)

from typing import List, Tuple
from sortedcontainers import SortedList


def shippingCenter(
    goods: List[Tuple[int, int]], boxes: List[int], queries: List[Tuple[int, int]]
) -> List[int]:
    def ban(left: int, right: int) -> int:
        sl = SortedList(boxes[:left] + boxes[right + 1 :])
        res = 0
        for weight, value in goods:
            pos = sl.bisect_left(weight)
            if pos < len(sl):
                res += value
                sl.pop(pos)
        return res

    goods.sort(key=lambda x: x[1], reverse=True)
    return [ban(left, right) for left, right in queries]


n, m, q = map(int, input().split())
goods = [tuple(map(int, input().split())) for _ in range(n)]  # 重量，价值
boxes = list(map(int, input().split()))  # 容量
queries = []
for _ in range(q):
    left, right = map(int, input().split())
    left, right = left - 1, right - 1
    queries.append((left, right))

res = shippingCenter(goods, boxes, queries)
print(*res, sep="\n")
