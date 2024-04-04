# 对每个固定的右端点`right(0<=right<n)`，找到最小的左端点`minLeft`，
# 使得滑动窗口内的元素满足`predicate(minLeft,right)`成立.
# 如果不存在，`minLeft`为`n`.


from collections import defaultdict
from typing import Callable, List


def getMinLeft(
    n: int,
    append: Callable[[int], None],
    popLeft: Callable[[int], None],
    predicate: Callable[[int, int], bool],
) -> List[int]:
    minLeft = [0] * n
    left = 0
    for right in range(n):
        append(right)
        while left <= right and not predicate(left, right):
            popLeft(left)
            left += 1
        minLeft[right] = n if left > right else left
    return minLeft


if __name__ == "__main__":
    # https://leetcode.cn/problems/longest-substring-without-repeating-characters/description/
    class Solution:
        def lengthOfLongestSubstring(self, s: str) -> int:
            def append(right: int) -> None:
                if counter[s[right]] == 1:
                    nonlocal dupCount
                    dupCount += 1
                counter[s[right]] += 1

            def popLeft(left: int) -> None:
                if counter[s[left]] == 2:
                    nonlocal dupCount
                    dupCount -= 1
                counter[s[left]] -= 1

            def predicate(left: int, right: int) -> bool:
                return dupCount == 0

            n = len(s)
            counter, dupCount = defaultdict(int), 0
            minLeft = getMinLeft(n, append=append, popLeft=popLeft, predicate=predicate)
            res = 0
            for right, left in enumerate(minLeft):
                res = max(res, right - left + 1)
            return res
