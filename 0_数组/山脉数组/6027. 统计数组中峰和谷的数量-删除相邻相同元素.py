from typing import List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def countHillValley(self, nums: List[int]) -> int:
        arr = []
        for num in nums:
            if not arr or num != arr[-1]:
                arr.append(num)
        res = 0
        for i in range(1, len(arr) - 1):
            if arr[i - 1] < arr[i] > arr[i + 1]:
                res += 1
            elif arr[i - 1] > arr[i] < arr[i + 1]:
                res += 1

        return res


print(Solution().countHillValley(nums=[2, 4, 1, 1, 6, 5]))
