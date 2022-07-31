# 蔚来笔试
# 对数组nums的一次操作规定如下：将元素置于第一个位置；将元素置于最后一个位置。
# len(nums)<=3e5
# max(nums)<=1e9
# !求修改为非递减数组的最小操作次数

# 最终肯定是部分数移动过，部分数没移动。
# 只要确定不需要移动的元素的最大数目，就可以得到答案

from typing import List


def minCost(nums: List[int]) -> int:
    arr = [(num, i) for i, num in enumerate(nums)]
    arr.sort()
    print(arr)


if __name__ == "__main__":
    nums = [1, 5, 4, 3, 4, 6]
    print(minCost(nums))


# 先排序，并记录每个数下标，数字相同，下标小的排前面。
# 然后用树状数组记录位子i之前有多少数字已经置换出去了。
# 最后从小到大遍历，如果当前的下标减去之前已经置换的个数不为零，
# 那么答案加一，并更新树状数组。总的时间复杂度是nlogn
