# 对每个固定的左端点`left(0<=left<n)`，找到最大的右端点`maxRight`，
# 使得滑动窗口内的元素满足`predicate(left,maxRight)`成立.
# 如果不存在，`maxRight`为-1.


from collections import defaultdict
from typing import Callable, List


def getMaxRight(
    n: int,
    append: Callable[[int], None],
    popLeft: Callable[[int], None],
    predicate: Callable[[int, int], bool],
) -> List[int]:
    maxRight = [0] * n
    right = 0
    visitedRight = [False] * n
    for left in range(n):
        if right < left:
            right = left
        while right < n:
            if not visitedRight[right]:
                visitedRight[right] = True
                append(right)
            if predicate(left, right):
                right += 1
            else:
                break
        if right == n:
            maxRight[left:] = [n - 1] * (n - left)
        maxRight[left] = right - 1 if right - 1 >= left else -1
        popLeft(left)
    return maxRight


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
            maxRight = getMaxRight(n, append=append, popLeft=popLeft, predicate=predicate)
            res = 0
            for left, right in enumerate(maxRight):
                res = max(res, right - left + 1)
            return res
