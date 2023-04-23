from math import lcm


# 2513. 最小化两个数组中的最大值
# 给你两个数组 arr1 和 arr2 ，它们一开始都是空的。你需要往它们中添加正整数，使它们满足以下条件：

# arr1 包含 uniqueCnt1 个 互不相同 的正整数，每个整数都 不能 被 divisor1 整除 。
# arr2 包含 uniqueCnt2 个 互不相同 的正整数，每个整数都 不能 被 divisor2 整除 。
# arr1 和 arr2 中的元素 互不相同 。
# 给你 divisor1 ，divisor2 ，uniqueCnt1 和 uniqueCnt2 ，请你返回两个数组中 最大元素 的 最小值 。


class Solution:
    def minimizeSet(self, divisor1: int, divisor2: int, uniqueCnt1: int, uniqueCnt2: int) -> int:
        def check(mid: int) -> bool:
            """[1,mid]时是否满足条件"""
            count1 = mid - mid // divisor1  # 可以取
            count2 = mid - mid // divisor2  # 可以取
            overlap = mid // lcm_  # 重复计算的个数
            return (
                count1 >= uniqueCnt1
                and count2 >= uniqueCnt2
                and mid - overlap >= uniqueCnt1 + uniqueCnt2
            )

        lcm_ = lcm(divisor1, divisor2)
        left, right = uniqueCnt1 + uniqueCnt2, int(1e18)
        while left <= right:
            mid = (left + right) // 2
            if check(mid):
                right = mid - 1
            else:
                left = mid + 1
        return left


print(Solution().minimizeSet(divisor1=9, divisor2=4, uniqueCnt1=8, uniqueCnt2=3))
# 11
print(Solution().minimizeSet(divisor1=5, divisor2=5, uniqueCnt1=9, uniqueCnt2=3))
# 14
