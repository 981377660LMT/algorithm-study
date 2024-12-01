# D - 1122 Substring
# https://atcoder.jp/contests/abc381/tasks/abc381_d
# 给定一个字符串，问它的最长子串是多少，使得每k位相同，且每个字符出现的次数为k。


from collections import defaultdict
from typing import List


def max2(a: int, b: int) -> int:
    return a if a > b else b


def longest1122Subarray(nums: List[int], k: int) -> int:
    def calc(start: int) -> int:
        pos = []
        for i in range(start, n - (k - 1), k):
            if len(set(nums[i : i + k])) == 1:
                pos.append(nums[i])
            else:
                pos.append(None)

        # 无重复字符的最长子串
        res = 0
        left = 0
        counter = defaultdict(int)
        for right in range(len(pos)):
            if pos[right] is None:
                left = right + 1
                counter.clear()
                continue
            counter[pos[right]] += 1
            while counter[pos[right]] > 1:
                counter[pos[left]] -= 1
                left += 1
            curLen = right - left + 1
            res = max2(res, curLen)
        return res * k

    n = len(nums)
    res = 0
    for i in range(k):
        res = max2(res, calc(i))
    return res


if __name__ == "__main__":
    N = int(input())
    A = list(map(int, input().split()))
    print(longest1122Subarray(A, k=3))
