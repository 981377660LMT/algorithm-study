# https://leetcode.cn/circle/discuss/osaLfj/
# 元素小于0的最少操作数

# 一个数组开始全为正数，每次可以选择让一个数减去X，
# 或是让全部数减1，目标使数组全部数小于0，
# 读入X和数组，问最小操作次数。

# 1<=len(nums)<=1e5
# 1<=x<=1e5
# 1<=nums[i]<=1e5

# 1.!贪心,如果数组中>=0的元素个数>=x,那么就全部减1,否则对最大的元素减x??
# !2.枚举用多少次全减1的操作
# 如果全部减1的个数是a，单个减x的个数为b。
# 那么a从最大值逐个减小的时候，满足与a模x同余的个数，
# 就是下次新加的b个数。


from typing import List


def minOperation(nums: List[int], x: int) -> int:
    nums = sorted(x + 1 for x in nums)  # 所有数变为<=0
    max_ = nums[-1]
    res = max_
    minusOne = 0
    mods = [0] * x
    for minusAll in range(max_, -1, -1):
        while nums and nums[-1] > minusAll:
            cur = nums.pop()
            mods[cur % x] += 1
            delta = cur - minusAll
            minusOne += delta // x + (delta % x > 0)
        res = min(res, minusOne + minusAll)
        minusOne += mods[minusAll % x]
    return res


if __name__ == "__main__":
    arr = [20, 25]
    x = 5
    print(minOperation(arr, x))
