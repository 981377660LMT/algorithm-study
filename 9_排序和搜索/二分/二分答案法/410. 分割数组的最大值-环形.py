# No.1211 円環はお断り(圆环)
# https://yukicoder.me/problems/no/1211
# !给定一个环形数组,分成k个非空连续子数组,使得这k个子数组的和的最小值最大,求出最大值
# 1<=k<=n<=10^5 1<=nums[i]<=10^9

# 0.断环成链
# 1.二分mid
# 2.预处理出从每个位置出发,最远可以走多远,使得和不超过mid


from bisect import bisect_left, bisect_right
from itertools import accumulate


if __name__ == "__main__":
    n, k = map(int, input().split())
    nums = list(map(int, input().split()))

    def check(mid: int) -> bool:
        """每段>=mid,是否能分成>=k段"""
        max_ = bisect_right(preSum, mid)  # 第一段的起点
        for start in range(max_):
            if start > n:
                break
            cur = start
            for _ in range(k):
                cur = bisect_left(preSum, preSum[cur] + mid, lo=cur)
                if cur > start + n:
                    break
            else:
                return True
        return False

    preSum = [0] + list(accumulate(nums + nums))
    left, right = 1, (sum(nums) // k) + 1
    while left <= right:
        mid = (left + right) // 2
        if check(mid):
            left = mid + 1
        else:
            right = mid - 1
    print(right)
