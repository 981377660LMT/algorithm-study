from itertools import accumulate

# 给出一个长度为n的01串，现在请你找到两个区间，
# 使得这两个区间中0和1的个数相同。两个区间可以相交，
# 但是不可以完全重叠，即两个区间的左右端点不可以完全相同。
# 现在请你找到两个最长的区间,满足以上要求。

# n<=1e6


def check(mid: int) -> bool:
    """限定区间长度为mid，是否存在两个子区间和相等"""
    global resValue, res
    visited = dict({})
    for left in range(len(string) - mid):
        right = left + mid
        curSum = preSum[right] - preSum[left]
        if curSum in visited:
            if mid > resValue:
                s1, t1 = visited[curSum]
                s2, t2 = left, right
                resValue = mid
                res = [s1, t1, s2, t2]
            return True
        visited[curSum] = (left, right)
    return False


string = input()
nums = [1 if char == '1' else 0 for char in string]
preSum = [0] + list(accumulate(nums))

res = [-1, -1, -1, -1]
resValue = -1
left, right = 1, len(string) + 5
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        left = mid + 1
    else:
        right = mid - 1
for num in res:
    print(num, end='')

