from typing import List

# 旋转数组
# A[i] 的取值范围是 [0, A.length]。
# 任何值小于或等于其索引的项都可以记作一分。
# 返回我们所能得到的最高分数对应的轮调索引 K。


# 如果把初始分数看成一个常数S,那么每一次调度都会导致其变化
# 那么就只考虑调度次数不同情况下的变化即可
# 那么先记录每个元素调度多少次后开始让分数较少 (什么时候下车)
# 即从某一步，需要扣一分
# 所以除了第一个元素移动到最后面会使score+1外，其他的数的左移会使score减小。
# 那么，最后只需要从标记数组中找到最大的下标即可
class Solution:
    def bestRotation(self, nums: List[int]) -> int:
        # score[k]:表示移动K步后，当前分数应该加几分，正数为加，负数为扣
        diff = [0] * len(nums)
        for i in range(len(nums)):
            diff[(i + 1 - nums[i]) % len(nums)] -= 1
        for i in range(1, len(nums)):
            diff[i] += diff[i - 1] + 1

        return diff.index(max(diff))


print(Solution().bestRotation([2, 3, 1, 4, 0]))
# 输出：3
# 解释：
# 下面列出了每个 K 的得分：
# K = 0,  A = [2,3,1,4,0],    score 2
# K = 1,  A = [3,1,4,0,2],    score 3
# K = 2,  A = [1,4,0,2,3],    score 3
# K = 3,  A = [4,0,2,3,1],    score 4
# K = 4,  A = [0,2,3,1,4],    score 3
# 所以我们应当选择 K = 3，得分最高。
