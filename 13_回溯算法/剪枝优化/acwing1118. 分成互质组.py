# 1≤n≤10
# 给定 n 个正整数，将它们分组，使得每组中任意两个数互质。
# 至少要分成多少个组？

# 搜索：
# 把某个数加到最后一组
# 新开一个组

from math import gcd
from typing import List


def check(group: List[int], num: int) -> bool:
    if not group:
        return True
    return all(gcd(group[i], num) == 1 for i in range(len(group)))


def bt(numIndex: int, groupIndex: int) -> None:
    global res

    # 不是答案，剪枝
    if groupIndex + 1 >= res:
        return
    if numIndex == n:
        res = min(res, groupIndex + 1)
        return

    # 优先先在已经使用的组里面找，看看能不能插进去
    for i in range(groupIndex + 1):
        if check(groups[i], nums[numIndex]):
            groups[i].append(nums[numIndex])
            bt(numIndex + 1, groupIndex)
            groups[i].pop()

    # 都不行，开一个新组
    if groupIndex + 1 <= n:
        groups[groupIndex + 1].append(nums[numIndex])
        bt(numIndex + 1, groupIndex + 1)
        groups[groupIndex + 1].pop()


n = int(input())
nums = list(map(int, input().split()))
nums = sorted(nums, reverse=True)
groups = [[] for _ in range(n + 1)]
res = n
bt(0, 0)
print(res)
