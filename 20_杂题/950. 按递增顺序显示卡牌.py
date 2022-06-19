# 现在，重复执行以下步骤，直到显示所有卡牌为止：

# 从牌组顶部抽一张牌，显示它，然后将其从牌组中移出。
# 如果牌组中仍有牌，则将下一张处于牌组顶部的牌放在牌组的底部。
# 如果仍有未显示的牌，那么返回步骤 1。否则，停止行动。


# !返回能以递增顺序显示卡牌的牌组顺序。
from collections import deque


class Solution:
    def deckRevealedIncreasing(self, deck: List[int]) -> List[int]:
        n = len(deck)
        res = [0] * n
        indexes = deque(range(n))
        nums = sorted(deck)
        for num in nums:
            res[indexes.popleft()] = num
            if indexes:
                indexes.append(indexes.popleft())

        return res

