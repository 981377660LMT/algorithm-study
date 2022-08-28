# 小红书t2
# # 小明是一名魔法师，他会n种法术，
# # 其中第i种法术的威力为ai
# # 他经常通过双手各自释放一种法术来提力，
# # 能得到的威力值为双手各自释放的法术的威力的乘积，
# # 但是他还不够强大
# # !不能双手释放同一种法术。这天他接到了一个任务，
# # 需要释放威力值至少为K才能完成，他想i他算一算，
# # !在两只手都释放法术的情况下，共有多少方案能达到威力值K。
# # 每种方案可记作(u, v),其威力值为au *av
# # (u,v)和(v,u)会被视为不同的方案。

# # !TODO 改善头尾双指针
from collections import Counter


n, k = map(int, input().split())
nums = list(map(int, input().split()))

nums.sort()
res = 0
right = n - 1
for num in nums:
    while right >= 0 and nums[right] * num >= k:
        right -= 1
    res += n - 1 - right


# TODO 好像不对
counter = Counter(nums)
for key in counter:
    if key * key >= k:
        count = 2 * counter[key]
        res -= count * (count - 1) // 2

print(res)
