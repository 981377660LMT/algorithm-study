class Solution:
    def validSubarrays(self, nums: list[int], k: int) -> int:
        n, k, res, cands = len(nums), k + 1, 0, [-1]
        for i in range(1, n - 1):
            if nums[i - 1] < nums[i] > nums[i + 1]:
                cands.append(i)
        cands.append(n)
        for i in range(1, len(cands) - 1):
            res += min(k, (cands[i] - cands[i - 1])) * min(k, (cands[i + 1] - cands[i]))
        return res
