# 是否能反转子数组的一段得到递增数组


class Solution:
    def solve(self, nums):
        if not nums:
            return True
        s, e, sn = -1, -1, sorted(nums)
        for num in range(len(nums)):
            if nums[num] != sn[num]:
                if s == -1:
                    s = num
                else:
                    e = num
        if s == -1 and e == -1:
            return True
        return (nums[:s] + nums[s : e + 1][::-1] + nums[e + 1 :]) == sn


print(Solution().solve(nums=[1, 3, 3, 7, 6, 9]))
# If we reverse the sublist [7, 6], then we can sort the list in ascending order: [1, 3, 3, 6, 7, 9].
