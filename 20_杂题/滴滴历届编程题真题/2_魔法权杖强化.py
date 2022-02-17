# 有一把魔法权杖，权杖上有 n 颗并排的法术石(编号为 1 到 n)。
# 每颗法术石具有一个能量值，权杖的法术强度等同于法术石的最小能量值。
# 权杖可以强化，一次强化可以将`两颗相邻`的法术石融合为一颗，
# 融合后的能量值为这两颗法术石能量值之和。现在有 m 次强化的机会，
# 请问权杖能 强化到的最大法术强度是多少？

n, m = map(int, input().split())
nums = list(map(int, input().split()))


def check(mid: int) -> bool:
    """贪心,遇到不行的就向右合并"""
    chance = m
    left = 0
    while left < len(nums):
        cur = nums[left]
        right = left + 1
        while right < len(nums) and cur < mid:
            if chance <= 0:
                return False
            cur += nums[right]
            right += 1
            chance -= 1
        left = right
    if cur >= mid:
        return True
    # 最后一块仍可合并
    return chance >= 1


l, r = min(nums), sum(nums)

while l <= r:
    mid = (l + r) >> 1
    if check(mid):
        l = mid + 1
    else:
        r = mid - 1

print(r)
