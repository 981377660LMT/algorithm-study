# 枚举字符串所有的分割方案


from functools import lru_cache
from typing import Generator, Sequence, List, TypeVar


T = TypeVar("T")


def genSplits(arr: Sequence[T]) -> Generator[List[Sequence[T]], None, None]:
    """遍历序列所有的分割方案."""
    if not arr:
        return
    n = len(arr)
    for state in range(1 << (n - 1)):  # 枚举n-1个分割点
        preSplit = 0
        cur = []
        for i in range(n - 1):
            if state & (1 << i):
                cur.append(arr[preSplit : i + 1])
                preSplit = i + 1
        cur.append(arr[preSplit:])
        yield cur


if __name__ == "__main__":
    print(*genSplits([1, 2, 3]))

    # https://leetcode.cn/problems/find-the-punishment-number-of-an-integer/
    @lru_cache(None)
    def check(num: int) -> bool:
        sb = list(str(num * num))
        for splits in genSplits(sb):
            if sum(int("".join(s)) for s in splits) == num:
                return True
        return False

    class Solution:
        def punishmentNumber(self, n: int) -> int:
            return sum(i * i for i in range(1, n + 1) if check(i))
