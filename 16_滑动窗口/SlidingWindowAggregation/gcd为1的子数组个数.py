from math import gcd

from SlidingWindowAggregation import SlidingWindowAggregation


if __name__ == "__main__":
    n = int(input())
    nums = list(map(int, input().split()))
    S = SlidingWindowAggregation(lambda: 0, gcd)
    res = 0
    right = 0
    for left in range(n):
        right = max(right, left)
        while right < n and gcd(S.query(), nums[right]) != 1:
            S.append(nums[right])
            right += 1
        res += n - right
        S.popleft()
    print(res)
