# 给定你一个长度为 n 的整数数列。
# 请你使用快速排序对这个数列按照从小到大进行排序。
# 并将排好序的数列按顺序输出。


from typing import List


n, k = list(map(int, input().split()))
nums = list(map(int, input().split()))

# 给定一个长度为 n 的整数数列，以及一个整数 k，请用快速选择算法求出数列从小到大排序后的第 k 个数。


# 根据与pivot的关系, 将元素分成三类, 三个区间: small, same, big
# 然后判断k落子哪个区间里面, 然后在哪个区间里面继续查找.
def kth(nums: List[int], k: int) -> int:
    pivot = nums[len(nums) // 2]
    left = [x for x in nums if x < pivot]
    mid = [x for x in nums if x == pivot]
    right = [x for x in nums if x > pivot]
    if k <= len(left):
        return kth(left, k)
    elif k <= len(left) + len(mid):
        return pivot
    else:
        return kth(right, k - len(left) - len(mid))


print(kth(nums, k))
