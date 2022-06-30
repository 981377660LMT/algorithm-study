from math import ceil
from typing import List


class Solution:
    def minimumTime(self, jobs: List[int], workers: List[int]) -> int:
        """
        Each job should be assigned to exactly one worker, such that each worker completes exactly one job.
        Return the minimum number of days needed to complete all the jobs after assignment.
        """
        jobs, workers = sorted(jobs), sorted(workers)
        return max(ceil(j / w) for j, w in zip(jobs, workers))


print(Solution().minimumTime(jobs=[5, 2, 4], workers=[1, 7, 5]))
print(Solution().minimumTime(jobs=[3, 18, 15, 9], workers=[6, 5, 1, 3]))
