# m, n = list(map(int, input().split()))
# nums = list(map(int, input().split()))
m, n = (5, 5)
nums = [4, 1, 4, 1, 2]
# 两人默契的二元组 <l,r> 一共有多少种
# 枚举小美选的数，看小团选的数最小可以是多少


def check(lower: int, upper: int) -> bool:
    pre = -1
    for num in nums:
        if num < lower or upper < num:
            if pre > num:
                return False
            pre = num
    return True


res = 0

for lower in range(1, m + 1):
    left, right = lower, m
    while left <= right:
        mid = (left + right) >> 1
        # 如果一个r可以，那么比r大的都可以
        if check(lower, mid):
            right = mid - 1
        else:
            left = mid + 1
    res += m - right
    # 关键是这一句
    if m == right:
        break


print(res)

