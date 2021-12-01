from typing import List

# https://leetcode-cn.com/problems/er94lq/solution/mo-ni-xi-pai-guo-cheng-xiang-xi-fen-xi-li-jie-kde-/
# 取第一次 计算公共前缀长度len k必须小于等于len:
# 你第一次洗牌后取前k个数，其中前Len个数相同，第Len+1个数不同，那得到的结果肯定不同。
# 那么k能不能小于Len呢？
# 如图我们取出前Len-1个数，那么第Len个数就肯定被保留下来，参与下一轮洗牌;firstSort数组中的第一个数肯定会因为洗牌与target中相同的数错开.两个数组中的第一个数肯定不同。
class Solution:
    def isMagic(self, target: List[int]) -> bool:
        if target[0] != 2:
            return False
        n = len(target)

        # 求最大的k，即第一次的公共前缀长
        arr = [*range(2, n + 1, 2), *range(1, n + 1, 2)]
        k = 0
        i = 0
        while i < n and target[i] == arr[i]:
            k += 1
            i += 1

        while i < n:
            arr[i:] = arr[i + 1 :: 2] + arr[i::2]
            if arr[i : i + k] != target[i : i + k]:
                return False
            i += k
        return True

    # 暴力模拟9000ms
    def isMagic2(self, target: List[int]) -> bool:
        if target[0] != 2:
            return False
        n = len(target)
        init = [i for i in range(1, n + 1)]

        # 模拟
        for k in range(1, n + 1):
            cardsdraw = []
            cardsremain = init[:]  # 可以过
            # cardsremain = [i for i in range(1, n + 1)]  # 超时
            while True:
                cardsremain = cardsremain[1::2] + cardsremain[::2]
                cardsdraw += cardsremain[:k]
                cardsremain = cardsremain[k:]
                if cardsremain == []:
                    break
            if cardsdraw == target:
                return True
        return False


print(Solution().isMagic([2, 4, 3, 1, 5]))
