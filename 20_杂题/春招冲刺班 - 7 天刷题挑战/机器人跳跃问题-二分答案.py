# 1 <= N <= 10^5
# 1 <= H(i) <= 10^5
# 这个数量级可能是贪心排序/线性dp/堆/二分

n = int(input())
nums = list(map(int, input().split()))

MAX = int(1e5) + 10


def check(mid: int) -> bool:
    for num in nums:
        mid += mid - num
        if mid > MAX:
            return True
        elif mid < 0:
            return False
    return True


left, right = 1, MAX
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        right = mid - 1
    else:
        left = mid + 1

print(left)
