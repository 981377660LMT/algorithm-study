# 旅店的业务分为两种，入住和退房：

# 旅客入住时，第 i 组旅客需要根据他们的人数 Di，给他们安排 Di 个连续的房间，并且房间号要尽可能的小。如果房间不够，则无法安排。
# 旅客退房时，第 i 组旅客的账单将包含两个参数 Xi 和 Di，你需要将房间号 Xi 到 Xi+Di−1 之间的房间全部清空。

# 现在你需要帮助该旅馆处理 M 单业务。
# 每个入住业务输出一个整数，表示要安排的房间序列中的第一个房间的号码。
# 如果没办法安排，则输出 0。


##################################################################################################
# 维护一个长度为 n 的01串，初始全部为0，进行 m 次操作

# 找到最靠左的长度为 len 的全部为0的子串，输出左端点并把这个子串全部赋值为1，如果不存在就输出0
# 把左端点为 l 长度为 len 的子串全部赋值为0


class Node:
    # 区间内最长的 0 串长度、左端点开始的最长 0 串长度、右端点开始的最长 0 串长度
    __slots__ = ('left', 'right', 'isLazy', 'lazyValue', 'max', 'lMax', 'rMax')

    def __init__(self) -> None:
        self.left = -1
        self.right = -1
        self.isLazy = False
        self.lazyValue = -1
        self.max = 0
        self.lMax = 0
        self.rMax = 0


class SegmentTree:
    def __init__(self, size: int) -> None:
        self.tree = [Node() for _ in range(size << 2)]
        self.build(1, 1, size)

    def build(self, rt: int, left: int, right: int) -> None:
        root = self.tree[rt]
        root.left = left
        root.right = right
        if left == right:
            root.max = root.lMax = root.rMax = 1
            return

        mid = (left + right) >> 1
        self.build(rt << 1, left, mid)
        self.build(rt << 1 | 1, mid + 1, right)
        self.pushUp(rt)

    def query(self, rt: int, length: int) -> int:
        """是否存在长度为 length 的 0 串，不存在返回0"""
        root = self.tree[rt]
        if root.max < length:
            return 0
        if root.max == length and root.right - root.left + 1 == length:
            return root.left

        self.pushDown(rt)
        # mid = (root.left + root.right) >> 1
        left, right = self.tree[rt << 1], self.tree[rt << 1 | 1]

        # 讨论 左 中 右
        if left.max >= length:
            return self.query(rt << 1, length)
        if left.rMax + right.lMax >= length:
            return left.right - left.rMax + 1
        if right.max >= length:
            return self.query(rt << 1 | 1, length)

        return 0

    def update(self, rt: int, left: int, right: int, target: int) -> None:
        """区间染色"""
        root = self.tree[rt]
        if left <= root.left and root.right <= right:
            root.isLazy = True
            root.lazyValue = target
            if target == 0:
                root.max = root.lMax = root.rMax = root.right - root.left + 1
            else:
                root.max = root.lMax = root.rMax = 0
            return

        self.pushDown(rt)
        mid = (root.left + root.right) >> 1
        if left <= mid:
            self.update(rt << 1, left, right, target)
        if mid < right:
            self.update(rt << 1 | 1, left, right, target)
        self.pushUp(rt)

    def pushDown(self, rt: int) -> None:
        root, left, right = self.tree[rt], self.tree[rt << 1], self.tree[rt << 1 | 1]
        if root.isLazy:
            if root.lazyValue == 0:
                left.max = left.lMax = left.rMax = left.right - left.left + 1
                right.max = right.lMax = right.rMax = right.right - right.left + 1
            else:
                left.max = left.lMax = left.rMax = 0
                right.max = right.lMax = right.rMax = 0

            left.isLazy = True
            left.lazyValue = root.lazyValue
            right.isLazy = True
            right.lazyValue = root.lazyValue
            root.isLazy = False
            root.lazyValue = -1

    def pushUp(self, rt: int) -> None:
        root, left, right = self.tree[rt], self.tree[rt << 1], self.tree[rt << 1 | 1]

        root.lMax = left.lMax
        if left.max == left.right - left.left + 1:
            root.lMax = max(root.lMax, left.max + right.lMax)

        root.rMax = right.rMax
        if right.max == right.right - right.left + 1:
            root.rMax = max(root.rMax, right.max + left.rMax)

        root.max = max(left.max, right.max, left.rMax + right.lMax)


n, m = map(int, input().split())
sg = SegmentTree(n)
res = []
for _ in range(m):
    opt, *rest = map(int, input().split())
    if opt == 1:
        # “1 Di”表示这单业务为入住业务。
        size = rest[0]
        pos = sg.query(1, size)
        res.append(str(pos))
        if pos != 0:
            # 染成1
            sg.update(1, pos, pos + size - 1, 1)
    else:
        # “2 Xi Di”表示这单业务为退房业务。
        left, size = rest
        # 染成0
        sg.update(1, left, left + size - 1, 0)


print('\n'.join(res))
