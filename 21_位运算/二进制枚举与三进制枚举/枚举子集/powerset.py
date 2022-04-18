from itertools import chain, combinations
from typing import Collection


def powerset(collection: Collection, isProperSubset=False):
    """求(真)子集,时间复杂度O(n*2^n)

    默认求所有子集
    """
    upper = len(collection) if isProperSubset else len(collection) + 1
    return chain.from_iterable(combinations(collection, n) for n in range(upper))


if __name__ == '__main__':
    print(*powerset([1, 2, 3, 4]))


# chain存在的意义是什么?
# 屏蔽可迭代对象的差异 提供了统一的接口
