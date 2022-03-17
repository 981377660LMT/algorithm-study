from itertools import chain, combinations
from typing import Any, Collection, List


def main(nums: List[int]):
    n = len(nums)
    res = []
    for state in range(1 << n):
        group1, group2 = state, 0
        while group1:
            res.append((group1, group2))
            # 关键，不断减一+与运算跳数
            group1 = state & (group1 - 1)
            group2 = state ^ group1
    return res


# chain存在的意义是什么?
# 用于简化yielf from 写法 配合生成式使用
def powerset(collection: Collection[Any]):
    """求真子集"""
    return chain.from_iterable(combinations(collection, n) for n in range(len(collection)))


if __name__ == '__main__':
    print(main([1, 2, 3, 4]))

