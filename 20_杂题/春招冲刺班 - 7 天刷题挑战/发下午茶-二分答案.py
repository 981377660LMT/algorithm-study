# https://leetcode.cn/problems/OMrszv/
# 0< K, N <= 1000

# !发下午茶
# 现在告诉你字节君的数量以及每个工具需要的下午茶个数，
# 请问所有的字节君最少花费多长时间才能送完所有的下午茶？
people, area = list(map(int, input().split()))
need = list(map(int, input().split()))


def check(mid: int) -> bool:
    """
    mid时间以内能否送完
    先考虑最远的，由远到近考虑
    """
    curNeed = need[:]
    nextArea = area
    for _ in range(people):
        power = mid - nextArea
        while power > 0:
            if power >= curNeed[nextArea - 1]:
                power -= curNeed[nextArea - 1]
                nextArea -= 1
                if nextArea == 0:
                    return True
            else:
                curNeed[nextArea - 1] -= power
                break

    return False


left, right = 0, int(1e18)
while left <= right:
    mid = (left + right) // 2
    if check(mid):
        right = mid - 1
    else:
        left = mid + 1
print(left)
