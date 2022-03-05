# 给定你一个长度为 n 的整数数列。
# 请你使用快速排序对这个数列按照从小到大进行排序。
# 并将排好序的数列按顺序输出。


n = int(input())
nums = list(map(int, input().split()))


def quick_sort(nums):
    if len(nums) <= 1:
        return nums
    pivot = nums[len(nums) // 2]
    left = [x for x in nums if x < pivot]
    mid = [x for x in nums if x == pivot]
    right = [x for x in nums if x > pivot]
    return quick_sort(left) + mid + quick_sort(right)


print(' '.join([str(num) for num in quick_sort(nums)]))
