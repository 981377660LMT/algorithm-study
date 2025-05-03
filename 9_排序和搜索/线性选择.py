# https://leetcode.cn/problems/put-marbles-in-bags/solutions/2081183/onxian-xing-xuan-ze-by-hqztrue-ykur/
# https://leetcode.cn/problems/maximum-ice-cream-bars/solutions/732895/c-on-80ms-100-by-hqztrue-x357/
# 用中位数分治思想可以 O(n) 时间解决的类似题目 https://leetcode.cn/problems/minimum-operations-to-halve-array-sum/solutions/1352188/onsuan-fa-by-hqztrue-jalf/

import numpy as np


def demo():
    a = np.array([7, 2, 5, 3, 9])
    k = 2
    b = np.partition(a, k)
    print(b)
    print(b[k])  # 输出: 5 （第3小的元素，索引从0开始）
