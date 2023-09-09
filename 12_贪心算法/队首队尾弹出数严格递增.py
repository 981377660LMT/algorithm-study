# 队首队尾弹出数严格递增
# https://codeforces.com/problemset/problem/1157/C2

# 输入 n(1≤n≤2e5) 和长为 n 的双端队列 a(1≤a[i]≤2e5)。
# 每次操作，弹出 a 的队首或队尾。
# 从第二次操作开始，弹出的数字必须严格大于上一次弹出的数字。
# 输出最多可以弹出多少个数字，以及操作序列（队首为 L，队尾为 R）。

# https://codeforces.com/problemset/submission/1157/212155329

# 贪心模拟即可。

# 哪边小选哪边（但必须大于上一个数）。
# 如果两边一样，那么后续操作要么都选左，要么都选右，暴力比较选哪边更优（有更长的严格递增）。
# 注意【一样】的情况，在整个模拟过程中，至多发生一次（因为必须严格递增）。


from typing import List

INF = int(1e18)


def solve(nums: List[int]) -> List[str]:
    res = []
    left, right, pre = 0, len(nums) - 1, -INF
    while left <= right and (nums[left] > pre or nums[right] > pre):
        selectRight = nums[left] <= pre or pre < nums[right] < nums[left]
        if nums[left] == nums[right]:  # 最多发生一次
            nextLeft = left + 1
            while nextLeft < right and nums[nextLeft] > nums[nextLeft - 1]:
                nextLeft += 1
            nextRight = right - 1
            while nextRight > left and nums[nextRight] > nums[nextRight + 1]:
                nextRight -= 1
            selectRight = right - nextRight > nextLeft - left

        if selectRight:
            res.append("R")
            pre = nums[right]
            right -= 1
        else:
            res.append("L")
            pre = nums[left]
            left += 1

    return res


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    res = solve(nums)
    print(len(res))
    print("".join(res))
