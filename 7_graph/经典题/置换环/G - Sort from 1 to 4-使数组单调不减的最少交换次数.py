# https://atcoder.jp/contests/abc302/tasks/abc302_g
# G - Sort from 1 to 4 - 使数组单调不减的最少交换次数
# 给定一个仅包含1,2,3,4 的数组，一次操作可以交换任意两个数。
# 问最少进行多少次交换，使得数组单调不减(不降数组).


# 由于大小为x的置换环至少需要x-1次交换，
# !所以我们可以贪心地先考虑长度为1的置换环，然后长度为2的置换环，以此类推。
# 考虑长度从 1 到 4 的置换环即可。


from collections import defaultdict
from itertools import permutations
from typing import List


def sortFrom1to4(nums: List[int]) -> int:
    target = sorted(nums)
    mp = defaultdict(lambda: defaultdict(int))
    for a, b in zip(nums, target):
        mp[a][b] += 1

    res = 0
    for size in range(1, 5):  # 置换环大小为size
        for perm in permutations(range(1, 5), size):
            perm = perm + (perm[0],)  # 生成可能的置换环
            min_ = min(mp[pre][cur] for pre, cur in zip(perm, perm[1:]))
            res += min_ * (size - 1)
            for pre, cur in zip(perm, perm[1:]):
                mp[pre][cur] -= min_
    return res


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    print(sortFrom1to4(nums))
