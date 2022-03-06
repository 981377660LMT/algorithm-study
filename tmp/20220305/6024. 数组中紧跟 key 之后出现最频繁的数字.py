from typing import Counter, List, Tuple

MOD = int(1e9 + 7)


class Solution:
    def mostFrequent(self, nums: List[int], key: int) -> int:
        counter = Counter()
        for i in range(len(nums)):
            if nums[i] == key and i + 1 < len(nums):
                counter[nums[i + 1]] += 1
        return counter.most_common(1)[0][0]
