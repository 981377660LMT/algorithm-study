# 给定两个长度为n的数组A[i], B[i] (0<=A[i], B[i]<=1e9,n <=2e5)。
# !问有多少对(i,j)满足1<=i,j <=n,A[i] >= A[j],B[i] <= B[j].
# !注意(i,j)间没有大小的关系，所以(1,2)!=(2,1)
# !二维偏序 => 一个维度排序，另一个维度用数据结构维护
# 注意等号要取到 所以要一次处理完相同的数


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

    def queryRange(self, left: int, right: int) -> int:
        return self.query(right) - self.query(left - 1)


if __name__ == "__main__":
    n = int(input())
    nums1 = list(map(int, input().split()))
    nums2 = list(map(int, input().split()))

    events = [(a, b, i) for i, (a, b) in enumerate(zip(nums1, nums2))]
    events.sort(key=lambda x: (x[0], -x[1]))  # 注意是A[i] >= A[j],B[i] <= B[j]
    res, ei = [0] * n, 0
    bit = BIT1(int(1e9) + 10)
    while ei < n:
        group = [events[ei][2]]
        a, b = events[ei][0], events[ei][1]
        while ei + 1 < n and (events[ei + 1][0] == a and events[ei + 1][1] == b):
            group.append(events[ei + 1][2])
            ei += 1
        bit.add(b + 1, len(group))
        tmp = bit.queryRange(b + 1, int(1e9) + 10)
        for qi in group:
            res[qi] = tmp
        ei += 1

    print(sum(res))
