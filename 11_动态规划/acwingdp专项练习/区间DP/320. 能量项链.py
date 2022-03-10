# 如果前一颗能量珠的头标记为 m，尾标记为 r，后一颗能量珠的头标记为 r，尾标记为 n，则聚合后释放的能量为 m×r×n（Mars 单位），新产生的珠子的头标记为 m，尾标记为 n。
# 需要时，Mars 人就用吸盘夹住相邻的两颗珠子，通过聚合得到能量，直到项链上只剩下一颗珠子为止。
# 显然，不同的聚合顺序得到的总能量是不同的，请你设计一个聚合顺序，使一串项链释放出的总能量最大。

# 把合并n颗珠子的问题转化为合并（n+1）个数合并的问题，只不过有一个数是公用的，注意不要间断分割
from functools import lru_cache
import sys

sys.setrecursionlimit(int(1e9))

n = int(input())
nums = list(map(int, input().split()))
nums = nums * 2


@lru_cache(maxsize=None)
def dfs(left: int, right: int) -> int:
    """[left:right]这一段合并的代价最小"""
    """left-right连起来的这条边`不动`，枚举根这条边配对的点的位置，枚举切割位置"""
    if right - left < 2:
        return 0
    if right - left == 2:
        return nums[left] * nums[left + 1] * nums[left + 2]

    res = -int(1e20)
    # 这里是left+1
    for i in range(left + 1, right):
        res = max(res, dfs(left, i) + dfs(i, right) + nums[left] * nums[i] * nums[right])
    return res


res = 0
for i in range(n):
    res = max(res, dfs(i, i + n))
print(res)
