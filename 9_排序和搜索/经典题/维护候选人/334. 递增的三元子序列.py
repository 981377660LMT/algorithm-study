from typing import List

INF = int(1e18)


class Solution:
    def increasingTriplet(self, nums: List[int]) -> bool:
        """
        Use two variables first and second to track the smallest and the second smallest
        values seen so far. As we scan through nums:
          - If num <= first, update first = num
          - Elif num <= second, update second = num
          - Else, we've found num > second > first, so return True
        Time: O(n), Space: O(1)
        """
        first = second = INF
        for num in nums:
            if num <= first:
                first = num
            elif num <= second:
                second = num
            else:
                return True
        return False


if __name__ == "__main__":
    sol = Solution()
    print(sol.increasingTriplet([1, 2, 3, 4, 5]))  # True
    print(sol.increasingTriplet([5, 4, 3, 2, 1]))  # False
    print(sol.increasingTriplet([2, 1, 5, 0, 4, 6]))  # True (0,4,6)
