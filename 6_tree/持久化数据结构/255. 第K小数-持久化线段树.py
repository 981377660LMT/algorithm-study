# 给定长度为 N 的整数序列 A，下标为 1∼N。
# 现在要执行 M 次操作，其中第 i 次操作为给出三个整数 li,ri,ki，
# 求 即 A 的下标区间 [li,ri]中第 ki 小的数是多少。
# 静态问题，原数组不变

# 255. 第K小数-持久化线段树
# 1.离散化
# 2.建立一颗权值线段树，每个点存储的信息为该值域区间存在的数的个数。


class Node:
    __slots__ = ('count', 'lower', 'upper', 'left', 'right')

    def __init__(self, lower: int, upper: int, count=0) -> None:
        self.count = count
        self.lower = lower
        self.upper = upper
        self.left = None
        self.right = None

    def build(self) -> 'Node':
        if self.lower < self.upper:
            mid = (self.lower + self.upper) >> 1
            self.left = Node(self.lower, mid).build()
            self.right = Node(mid + 1, self.upper).build()
        return self

    def update(self, k: int) -> 'Node':
        newNode = Node(self.lower, self.upper, self.count + 1)
        # 非叶子结点
        if self.lower < self.upper:
            if k <= self.left.upper:
                newNode.right = self.right
                newNode.left = self.left.update(k)
            else:
                newNode.left = self.left
                newNode.right = self.right.update(k)
        return newNode


# 根据需要修改
def query(tree1: Node, tree2: Node, k: int) -> int:
    # 叶子结点
    if tree1.lower == tree1.upper:
        return tree1.lower
    leftCount = tree2.left.count - tree1.left.count
    return (
        query(tree1.left, tree2.left, k)
        if k <= leftCount
        else query(tree1.right, tree2.right, k - leftCount)
    )


n, m = map(int, input().split())
nums = list(map(int, input().split()))
sortedNums = sorted(nums)  # sorted(set(nums)) 会超时
mapping = {sortedNums[i]: i for i in range(len(sortedNums))}

# 默认查[0, N - 1]区间
trees = [Node(0, n - 1).build()]
for num in nums:
    trees.append(trees[-1].update(mapping[num]))

for _ in range(m):
    left, right, k = map(int, input().split())
    # 查询这个范围里的第k小的数(把每个值看做修改历史版本，即历史版本left-1到right里的第k小的数)
    tmp = query(trees[left - 1], trees[right], k)
    print(sortedNums[tmp])

