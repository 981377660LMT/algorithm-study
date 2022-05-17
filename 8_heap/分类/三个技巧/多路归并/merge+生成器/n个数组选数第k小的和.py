# 两个数组的情况
# 0<row≤1000,
# 0<col≤2000


from functools import reduce
from heapq import merge
from itertools import islice
from typing import Generator, List


def mergeTwo(nums1: List[int], nums2: List[int], select: int) -> List[int]:
    """两个数组选前k小的和"""

    def gen(index: int) -> Generator[int, None, None]:
        return (nums1[index] + num for num in nums2)

    allGen = (gen(i) for i in range(len(nums1)))
    iterable = merge(*allGen)  # merge 相当于多路归并
    return list(islice(iterable, select))


def main() -> None:
    row, k = map(int, input().split())
    matrix = []
    for _ in range(row):
        matrix.append(sorted(map(int, input().split())))

    res = reduce(lambda pre, cur: mergeTwo(pre, cur, k), matrix)
    for num in res:
        print(num, end=' ')
    print('')


input()
while True:
    try:
        main()
    except EOFError:
        break

