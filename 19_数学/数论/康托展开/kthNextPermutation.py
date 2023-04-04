from collections import defaultdict
from typing import Any, List, Tuple
from sortedcontainers import SortedSet


# https://maspypy.github.io/library/seq/kth_next_permutation.hpp
# !有重复的情况：
# kthPrevPermutation: 换成负数比较即可,最后换回来.
def kthNextPermutation(unique: List[Any], k: int, inPlace=False) -> Tuple[bool, List[Any], int]:
    """下k个字典序的排列

    Args:
        unique (List[Any]): `无重复元素`的数组
        k (int): 后续第k个(`本身算第0个`)
        inPlace (bool, optional): 是否原地修改. 默认为False

    Returns:
        Tuple[bool, List[Any], int]: `是否存在, 下k个排列, 需要移动的元素个数`
    """
    if not inPlace:
        unique = unique[:]
    rank, q = [], []
    ss = SortedSet()
    while k and unique:
        n = len(rank) + 1
        p = unique[-1]
        now = ss.bisect_left(p)
        k += now
        r = k % n
        k //= n
        rank.append(r)
        q.append(unique[-1])
        ss.add(unique[-1])
        unique.pop()

    if k:
        return False, unique, len(rank)

    move = len(rank)
    while len(rank):
        r = rank.pop()
        it = ss[r]
        unique.append(it)
        ss.remove(it)
    return True, unique, move


if __name__ == "__main__":
    # num = "11112", k = 4
    mp = defaultdict(int)
    print(kthNextPermutation([1, 2, 3, 4, 2], 1))
