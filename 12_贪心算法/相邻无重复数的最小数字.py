# 相邻无重复数的最小数字
# 问好可以替换成任何数
class Solution:
    def solve(self, s):
        nums = list(s)
        for index, num in enumerate(nums):
            if num == "?":
                for d in "123":
                    nums[index] = d
                    if index - 1 >= 0 and nums[index] == nums[index - 1]:
                        continue
                    if index + 1 < len(nums) and nums[index] == nums[index + 1]:
                        continue
                    break

        return "".join(nums)


print(Solution().solve("?1??2"))

# 只需要123即可
