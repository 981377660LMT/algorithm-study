# 输入参数的立方根。保留一位小数。

n = float(input())
left, right = -30, 30
while left <= right:
    mid = (left + right) / 2
    if mid ** 3 < n:
        left = mid + 0.01
    else:
        right = mid - 0.01

print(round(left, 1))

