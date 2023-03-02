# https://hitonanode.github.io/cplib-cpp/utilities/segtree_range_covering_nodes.hpp
# 非递归线段树结点编号对应


from typing import List


def segtreeRangeCoveringNodes(n: int, start: int, end: int) -> List[int]:
    """Enumerate nodes of nonrecursive segtree which cover [start, end)."""
    res, revRes = [], []
    start, end = start + n, end + n
    while start < end:
        if start & 1:
            res.append(start)
            start += 1
        if end & 1:
            end -= 1
            revRes.append(end)
        start >>= 1
        end >>= 1
    res += revRes[::-1]
    return res


def segtreeIndexCoveringNodes(n: int, index: int):
    """枚举包含index的元素的线段树结点编号."""
    index += n
    while index > 0:
        yield index
        index >>= 1


for i in segtreeIndexCoveringNodes(6, 1):
    print(i, end=" ")
