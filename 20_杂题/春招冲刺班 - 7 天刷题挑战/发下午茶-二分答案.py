# 0< K, N <= 1000
people, area = list(map(int, input().split()))
needs = list(map(int, input().split()))


def check(mid: int) -> bool:
    """先考虑最远的，由远到近考虑"""
    needs_ = needs[:]
    curArea = area
    for _ in range(people):
        power = mid - curArea
        while power > 0:
            if power >= needs_[curArea - 1]:
                power -= needs_[curArea - 1]
                curArea -= 1
                if curArea == 0:
                    return True
            else:
                needs_[curArea - 1] -= power
                power = 0

    return False


left, right = area, int(1e10)
while left <= right:
    mid = (left + right) >> 1
    if check(mid):
        right = mid - 1
    else:
        left = mid + 1

print(left)
