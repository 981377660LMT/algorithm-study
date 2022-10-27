from typing import Any, List, TypeVar
from functools import reduce
from operator import iconcat

T = TypeVar("T", Any, str, bytes, int, float, complex, bool, tuple, list, dict, set)


def flat(arr: List[List[T]]) -> List[T]:
    """二维数组flat

    todo : Nested list type
    """
    return reduce(iconcat, arr, [])
    return [item for pair in arr for item in pair]


if __name__ == "__main__":
    arr = [[1, 2], [3, 4], [5, 6]]
    res = flat(arr)
    print(res)
