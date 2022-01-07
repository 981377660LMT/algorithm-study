class Solution:
    def areNumbersAscending(self, s: str) -> bool:
        nums = [int(char) for char in s.split() if char.isdigit()]
        return all([a < b for a, b in zip((nums), (nums[1:]))])


print(
    Solution().areNumbersAscending(
        "sunset is at 7 51 pm overnight lows will be in the low 50 and 60 s"
    )
)
