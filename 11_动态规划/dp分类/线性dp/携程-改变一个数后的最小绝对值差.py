# 携程t4
# 游游得到了一个有n个数字的数列。
# 游游定义了"平滑值"的概念∶
# !平滑值指任意两个相邻的数的差的绝对值的最大值。
# 例如[1,2,5,7,8]的平滑值是3。
# 游游现在想知道，
# !在只修改一个位置的数字（可以修改为任意值)或者不修改的情况下，
# !数列的平滑值最小是多少?
# n <= 2e5
# -1e9 <= nums[i] <= 1e9

# !注意不能用dp 因为调整数字可以调左边 会导致dp是有后效性的
# 想dp做却碰到环的不适感 应该及时收手 考虑诸如贪心等解法

# !考虑贪心
# !先找到最大平滑值的两个索引
# 在数组左端点就把第一个数置为第二个数再求一遍平滑值；
# 在数组右端点就把最后一个数置为倒数第二个数再求一遍平滑值；
# 在数组中间就分别把平滑值左右端点置为相邻两个数的中值，取最优解

from typing import List


def solve(nums: List[int]) -> int:
    def calMaxDiff(changedNums: List[int]) -> int:
        res = 0
        for pre, cur in zip(changedNums, changedNums[1:]):
            res = max(res, abs(pre - cur))
        return res

    n = len(nums)
    maxDiff, maxIndex = 0, 0
    for index, (pre, cur) in enumerate(zip(nums, nums[1:])):
        diff = abs(pre - cur)
        if diff > maxDiff:
            maxDiff = diff
            maxIndex = index

    if maxIndex == 0:
        return calMaxDiff(nums[1:])
    if maxIndex == n - 2:
        return calMaxDiff(nums[:-1])

    mid1 = (nums[maxIndex - 1] + nums[maxIndex + 1]) // 2
    nums1 = nums[:maxIndex] + [mid1] + nums[maxIndex + 1 :]
    res1 = calMaxDiff(nums1)
    mid2 = (nums[maxIndex] + nums[maxIndex + 2]) // 2
    nums2 = nums[: maxIndex + 1] + [mid2] + nums[maxIndex + 2 :]
    res2 = calMaxDiff(nums2)
    return min(res1, res2)


if __name__ == "__main__":
    nums1 = [1, 3, 4]
    print(solve(nums1))  # 1
    nums2 = [-1, 1, 2, 5, 7]
    print(solve(nums2))  # 2
    nums3 = [1, 5, 18, 9, 9]
    print(solve(nums3))  # 4
