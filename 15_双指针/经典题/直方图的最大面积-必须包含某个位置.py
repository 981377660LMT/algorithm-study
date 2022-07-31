# `min(nums) * nums.length`最大的子数组 且子数组必须包含pos位上的元素
# 注意到这个乘积是直方图的面积

# 注意到滑动单向性，可以用滑动窗口来解决。
INF = int(1e20)


class Solution:
    def solve(self, nums, pos):
        res = min_ = nums[pos]
        i = j = pos

        # 从中心向两边扩展
        for _ in range(len(nums) - 1):
            left = nums[i - 1] if i - 1 >= 0 else -INF
            right = nums[j + 1] if j + 1 < len(nums) else -INF

            if left >= right:
                i -= 1
                min_ = min(min_, nums[i])
            else:
                j += 1
                min_ = min(min_, nums[j])

            res = max(res, min_ * (j - i + 1))

        return res


print(Solution().solve(nums=[-1, 1, 4, 3], pos=3))
# The best sublist is [4, 3]. Since min(4, 3) = 3 and its length is 2 we have 3 * 2 = 6.
