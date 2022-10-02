# video 中所有值 互不相同 。
# upload 和 longest 总调用 次数至多不超过 2 * 1e5 次。

# !请你实现一个数据结构，在上传的过程中计算 最长上传前缀 (mex)
# 如果 闭区间 1 到 i 之间的视频全部都已经被上传到服务器，
# 那么我们称 i 是上传前缀。
# !最长上传前缀指的是符合定义的 i 中的 最大值 。
# !求mex


class LUPrefix:
    def __init__(self, n: int):
        self.mex = 1
        self.visited = set()

    def upload(self, video: int) -> None:
        self.visited.add(video)

    def longest(self) -> int:
        while self.mex in self.visited:
            self.mex += 1
        return self.mex - 1


#############################################
# !如果增加了删除的操作,可以使用树状数组维护01序列,
# !然后树状数组树上二分


class BIT1:
    """单点修改"""

    __slots__ = "size", "bit", "tree"

    def __init__(self, n: int):
        self.size = n
        self.bit = n.bit_length()
        self.tree = dict()

    def add(self, index: int, delta: int) -> None:
        # assert index >= 1, 'index must be greater than 0'
        while index <= self.size:
            self.tree[index] = self.tree.get(index, 0) + delta
            index += index & -index

    def query(self, index: int) -> int:
        if index > self.size:
            index = self.size
        res = 0
        while index > 0:
            res += self.tree.get(index, 0)
            index -= index & -index
        return res

    def longest(self) -> int:
        """求最长上传前缀,即以1为起点的最长连续1区间的长度"""
        res = 0
        for i in range(self.bit, -1, -1):
            nextPos = res + (1 << i)
            if nextPos <= self.size and self.tree.get(nextPos, 0) == nextPos & -nextPos:
                res = nextPos
        return res


class LUPrefix:
    def __init__(self, n: int):
        self.bit = BIT1(n + 10)

    def upload(self, video: int) -> None:
        self.bit.add(video, 1)

    def longest(self) -> int:
        return self.bit.longest()


# 输入：
# ["LUPrefix", "upload", "longest", "upload", "longest", "upload", "longest"]
# [[4], [3], [], [1], [], [2], []]
# 输出：
# [null, null, 0, null, 1, null, 3]

# 解释：
# LUPrefix server = new LUPrefix(4);   // 初始化 4个视频的上传流
# server.upload(3);                    // 上传视频 3 。
# server.longest();                    // 由于视频 1 还没有被上传，最长上传前缀是 0 。
# server.upload(1);                    // 上传视频 1 。
# server.longest();                    // 前缀 [1] 是最长上传前缀，所以我们返回 1 。
# server.upload(2);                    // 上传视频 2 。
# server.longest();                    // 前缀 [1,2,3] 是最长上传前缀，所以我们返回 3 。

lru = LUPrefix(4)
lru.upload(3)
print(lru.longest())
lru.upload(1)
print(lru.longest())
lru.upload(2)
print(lru.longest())
