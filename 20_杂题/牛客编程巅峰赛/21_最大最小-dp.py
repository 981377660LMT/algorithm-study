# 有多少个区间满足区间最大值大于等于区间最小值的两倍。返回满足条件的区间个数。
# 只要找到第一个满足最大值大于等于2倍最小值的子区间，那么包含这个子区间的区间一定都满足要求。
# n<=10^5
class Solution:
    def MaxMin(self, array):
        # write code here
        res = 0
        for i in range(len(array)):
            min_ = array[i]
            max_ = array[i]
            for j in range(i + 1, len(array)):
                if array[j] < min_:
                    min_ = array[j]
                if array[j] > max_:
                    max_ = array[j]
                if max_ >= 2 * min_:
                    res += len(array) - j
                    break
        return res

