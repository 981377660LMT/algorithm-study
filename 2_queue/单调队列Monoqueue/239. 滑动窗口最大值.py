from typing import List
from MaxQueue import MaxQueue


def maxSlidingWindow(nums: List[int], k: int) -> List[int]:
    """滑动窗口最大值"""
    queue = MaxQueue()
    res = []
    n = len(nums)
    for right in range(n):
        queue.append(nums[right])
        if right >= k:
            queue.popleft()
        if right >= k - 1:
            res.append(queue.max)
    return res


if __name__ == "__main__":
    print(maxSlidingWindow([1, 3, -1, -3, 5, 3, 6, 7], 3))
