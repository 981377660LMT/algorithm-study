from typing import List, Sequence, TypeVar


T = TypeVar("T")


def chunk(arr: Sequence[T], maxSize: int) -> List[List[T]]:
    if maxSize <= 1:
        return [[item] for item in arr]
    res = []
    ptr = 0
    while ptr < len(arr):
        res.append(arr[ptr : ptr + maxSize])
        ptr += maxSize
    return res


if __name__ == "__main__":
    arr = [1, 2, 3, 4, 5, 6, 7, 8, 9]
    print(chunk(arr, 2))
