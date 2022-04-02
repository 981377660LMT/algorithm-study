from collections import Counter


class Solution:
    def solve(self, nums, target):
        if not target:
            # every sublist contains every number in target, so the answer is zero
            return 0

        target = set(target)
        # `counter` will count elements in `nums[i:j]` that are in `target`
        counter = Counter()
        n = len(nums)
        left = 0
        res = 0
        for right in range(n):
            if nums[right] in target:
                counter[nums[right]] += 1
            # advance `i` until `nums[i:j+1]` satisfies the condition
            # that not every number in `target` is present
            while len(counter) == len(target):
                if nums[left] in target:
                    counter[nums[left]] -= 1
                    if not counter[nums[left]]:
                        del counter[nums[left]]
                left += 1
            res += right - left + 1
            res %= int(1e9 + 7)
        return res


print(Solution().solve([1, 2, 2], [1, 2]))
