# 1 <= left < right <= 1e9


class RangeModule:
    def __init__(self):
        ...

    def addRange(self, left: int, right: int) -> None:
        """添加 半开区间 [left, right)"""
        ...

    def queryRange(self, left: int, right: int) -> bool:
        """ 只有在当前正在跟踪区间 [left, right) 中的每一个实数时，才返回 true"""
        ...

    def removeRange(self, left: int, right: int) -> None:
        """ 停止跟踪 半开区间 [left, right)"""
        ...

