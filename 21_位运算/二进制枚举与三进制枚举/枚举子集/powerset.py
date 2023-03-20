from itertools import chain, combinations
from typing import Any, Collection


def powerset(collection: Collection[Any], isAll=True):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) + 1 if isAll else len(collection)
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


if __name__ == "__main__":
    res = powerset([1, 2, 3, 4])
    print(list(res))


# chain存在的意义是什么?
# 屏蔽可迭代对象的差异 提供了统一的接口
