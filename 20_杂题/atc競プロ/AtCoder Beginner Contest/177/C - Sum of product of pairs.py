# 所有二元组的乘积之和
# 公式变形
# (∑nums[i])^2=∑nums[i]^2+2∑nums[i]nums[j]


from typing import List

MOD = int(1e9 + 7)


def sumOfProductOfPairs(nums: List[int]) -> int:
    sum2 = sum(i * i for i in nums)
    sum_ = sum(nums)
    return (sum_ * sum_ - sum2) // 2 % MOD


n = int(input())
nums = list(map(int, input().split()))
print(sumOfProductOfPairs(nums))
