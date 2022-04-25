from collections import Counter


n, m = map(int, input().split())
cars = input()
likes = input()

# 共有多少个位置连续的子数组，
# 能够满足现有的所有车辆可以在同一时刻把子数组中的乘客同时接上乘客喜爱的颜色的车
counter = Counter(cars)


res, left = 0, 0
# 对每个右端点看车
for right in range(m):
    # 用完了车
    counter[likes[right]] -= 1
    # 不满足，收缩
    while left <= right and counter[likes[right]] < 0:
        counter[likes[left]] += 1
        left += 1
    res += right - left + 1

print(res)

