# 分治的迭代写法(down to top)


from functools import reduce
from typing import Callable, List, TypeVar


T = TypeVar("T")


def mergeAll(items: List[T], op: Callable[[T, T], T]) -> "T":
    """会改变items数组"""
    if len(items) == 0:
        raise ValueError("items must not be empty")
    len_ = len(items)
    while True:
        if len_ == 1:
            break
        mid = (len_ + 1) // 2
        for i in range(mid):
            if 2 * i + 1 == len_:
                items[i] = items[2 * i]
            else:
                items[i] = op(items[2 * i], items[2 * i + 1])
        len_ = mid
    return items[0]


lenSum = 0


def concat(a: List[int], b: List[int]) -> List[int]:
    global lenSum
    lenSum += len(a) + len(b)
    return a + b


# concat with reduce
print(reduce(concat, [[1], [2], [3], [4], [5], [6], [7], [8], [9], [10]]))
print("concat with reduce:", lenSum)
lenSum = 0


# concat with mergeAll
print(mergeAll([[1], [2], [3], [4], [5], [6], [7], [8], [9], [10]], concat))
print("concat with mergeAll:", lenSum)
lenSum = 0
