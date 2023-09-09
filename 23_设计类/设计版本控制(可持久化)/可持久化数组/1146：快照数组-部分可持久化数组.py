# !部分可持久化数组:
# 每个历史版本都可以访问，但只有最新的版本才可以修改

# 1 <= length <= 50000
# 题目最多进行50000 次set，snap，和 get的调用 。
# 0 <= index < length
# 0 <= snap_id < 我们调用 snap() 的总次数
# 0 <= val <= 10^9

# !这种方法又叫做胖节点(FatNode)
# 这种方法（即在每个索引上存储多个值而不擦除旧值）被称为胖节点。
# 虽然易于实现，但胖节点只能部分持久化，这意味着只能修改最新版本的数据结构
# 对于涉及持久数据结构的大多数竞争性编程问题，我们改用路径复制（PathCopying）方法。
# https://usaco.guide/adv/persistent?lang=cpp


from bisect import bisect_right


class SnapshotArray:
    def __init__(self, length: int):
        """初始化一个与指定长度相等的 类数组 的数据结构。初始时，每个元素都等于 0。"""
        self.version = 0
        self.actions = [[(0, 0)] for _ in range(length)]

    def set(self, index: int, val: int) -> None:
        """会将指定索引 index 处的元素设置为 val"""
        self.actions[index].append((self.version, val))  # type: ignore

    def snap(self) -> int:
        """获取该数组的快照，并返回快照的编号 snap_id(快照号是调用 snap() 的总次数减去 1)。"""
        self.version += 1
        return self.version - 1

    def get(self, index: int, snap_id: int) -> int:
        """根据指定的 snap_id 调用 snap()，获取该数组在该快照下的指定 index 的值

        二分action
        """
        pos = bisect_right(self.actions[index], snap_id, key=lambda x: x[0]) - 1
        return self.actions[index][pos][1]


if __name__ == "__main__":
    snapshotArray = SnapshotArray(3)
    snapshotArray.set(0, 5)
    snapshotArray.snap()
    snapshotArray.set(0, 6)
    print(snapshotArray.actions)
    print(snapshotArray.get(0, 0))
