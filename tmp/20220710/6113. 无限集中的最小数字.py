from sortedcontainers import SortedSet

# 6113. 无限集中的最小数字
# 假的无限集


class SmallestInfiniteSet:
    """
    1 <= num <= 1000
    最多调用 popSmallest 和 addBack 方法 共计 1000 次
    """

    def __init__(self):
        """初始化 SmallestInfiniteSet 对象以包含 所有 正整数。"""
        self.ss = SortedSet(range(1, int(1e4)))

    def popSmallest(self) -> int:
        """移除 并返回该无限集中的最小整数"""
        return self.ss.pop(0)

    def addBack(self, num: int) -> None:
        """如果正整数 num 不 存在于无限集中，则将一个 num 添加 到该无限集中"""
        self.ss.add(num)
