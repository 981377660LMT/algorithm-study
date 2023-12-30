# !完全可持久化数组:
# 每次修改都会生成一个新的数组

# 1 <= length <= 50000
# 题目最多进行50000 次set，snap，和 get的调用 。
# 0 <= index < length
# 0 <= snap_id < 我们调用 snap() 的总次数
# 0 <= val <= 10^9

from PersistentArray1 import PersistentArray as PA1
from PersistentArray2 import PersistentArray as PA2

Q = int(5e4 + 10)


class SnapshotArray1:
    """数组实现的可持久化数组"""

    def __init__(self, length: int):
        """初始化一个与指定长度相等的 类数组 的数据结构。初始时，每个元素都等于 0。"""
        self.snapId = 0
        self.pa = PA1.create(length, updateTimes=Q)
        self.git = dict({0: 0})  # !snap_id: version

    def set(self, index: int, val: int) -> None:
        """会将指定索引 index 处的元素设置为 val"""
        self.pa.update(self.pa.curVersion, index, val)

    def snap(self) -> int:
        """获取该数组的快照，并返回快照的编号 snap_id(快照号是调用 snap() 的总次数减去 1)。"""
        self.git[self.snapId] = self.pa.curVersion
        self.snapId += 1
        return self.snapId - 1

    def get(self, index: int, snap_id: int) -> int:
        """根据指定的 snap_id 调用 snap()，获取该数组在该快照下的指定 index 的值"""
        version = self.git[snap_id]
        return self.pa.query(version, index)


class SnapshotArray2:
    """结点实现的可持久化数组"""

    def __init__(self, length: int):
        """初始化一个与指定长度相等的 类数组 的数据结构。初始时，每个元素都等于 0。"""
        self.snapId = 0
        self.nums = PA2.create(length)
        self.git = dict({0: self.nums})  # !snap_id: nums

    def set(self, index: int, val: int) -> None:
        """会将指定索引 index 处的元素设置为 val"""
        self.nums = self.nums.update(index, val)

    def snap(self) -> int:
        """获取该数组的快照，并返回快照的编号 snap_id(快照号是调用 snap() 的总次数减去 1)。"""
        self.git[self.snapId] = self.nums
        self.snapId += 1
        return self.snapId - 1

    def get(self, index: int, snap_id: int) -> int:
        """根据指定的 snap_id 调用 snap()，获取该数组在该快照下的指定 index 的值"""
        return self.git[snap_id].get(index)


if __name__ == "__main__":
    snapshotArray = SnapshotArray1(3)
    snapshotArray.set(0, 5)
    snapshotArray.snap()
    snapshotArray.set(0, 6)
    print(snapshotArray.get(0, 0))

    snapshotArray = SnapshotArray2(3)
    snapshotArray.set(0, 5)
    snapshotArray.snap()
    snapshotArray.set(0, 6)
    print(snapshotArray.get(0, 0))
