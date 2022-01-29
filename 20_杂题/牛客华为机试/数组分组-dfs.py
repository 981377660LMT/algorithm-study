# 询问该数组能否分成两组，使得两组中各元素加起来的和相等，并且，
# 所有5的倍数必须在其中一个组中，所有3的倍数在另一个组中（不包括5的倍数），
# 不是5的倍数也不是3的倍数能放在任意一组，可以将数组分为空数组，

# n<=50,abs(ni)<=500

# 不告诉输入，只能使用while 循环 + try/except终止
from typing import List


def canPartition(nums: List[int], target: int) -> bool:
    """选出若干个nums是否能组成容量target；注意nums不全为正数，不能01背包，只能枚举子集/dfs"""
    for state in range(1 << len(nums)):
        curSum = 0
        for i in range(len(nums)):
            if (state >> i) & 1:
                curSum += nums[i]
        if curSum == target:
            return True
    return False


while True:
    try:
        n = int(input())
        nums = list(map(int, input().split()))
        nums5, nums3, spare = [], [], []
        # 不同类数据分组
        for num in nums:
            if num % 5 == 0:
                nums5.append(num)
            elif num % 3 == 0:
                nums3.append(num)
            else:
                spare.append(num)
        diff = abs(sum(nums5) - sum(nums3))
        target = sum(spare) - diff

        if target % 2 != 0:
            print('false')
        elif canPartition(spare, target // 2):
            print('true')
        else:
            print('false')
    except EOFError:
        break

