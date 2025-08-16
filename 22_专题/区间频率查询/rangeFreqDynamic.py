from typing import Generic, List, TypeVar


from sortedcontainers import SortedList

T = TypeVar("T")


class RangeFreqDynamic(Generic[T]):
    __slots__ = "_data", "_valueToIndexes"

    def __init__(self, data: List[T]) -> None:
        self._data = data[:]
        self._valueToIndexes = dict()
        mp = dict()
        for i, v in enumerate(data):
            mp.setdefault(v, []).append(i)
        for v, indexes in mp.items():
            self._valueToIndexes[v] = SortedList(indexes)

    def query(self, start: int, end: int, value: T) -> int:
        """[start,end)区间内value的频率."""
        if start >= end:
            return 0
        if value not in self._valueToIndexes:
            return 0
        pos = self._valueToIndexes[value]
        return pos.bisect_left(end) - pos.bisect_left(start)

    def set(self, index: int, value: T) -> bool:
        pre = self._data[index]
        if pre == value:
            return False

        if pre in self._valueToIndexes:
            self._valueToIndexes[pre].remove(index)
            if len(self._valueToIndexes[pre]) == 0:
                del self._valueToIndexes[pre]

        self._data[index] = value

        if value not in self._valueToIndexes:
            self._valueToIndexes[value] = SortedList()
        self._valueToIndexes[value].add(index)

        return True

    def get(self, index: int) -> T:
        return self._data[index]

    def findFirst(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return -1
        if value not in self._valueToIndexes:
            return -1
        pos = self._valueToIndexes[value]
        idx = pos.bisect_left(start)
        if idx < len(pos):
            tmp = pos[idx]
            if tmp < end:  # type: ignore
                return tmp  # type: ignore
        return -1

    def findLast(self, start: int, end: int, value: T) -> int:
        if start >= end:
            return -1
        if value not in self._valueToIndexes:
            return -1
        pos = self._valueToIndexes[value]
        idx = pos.bisect_left(end)
        if idx > 0:
            tmp = pos[idx - 1]
            if tmp >= start:  # type: ignore
                return tmp  # type: ignore
        return -1


if __name__ == "__main__":

    def test_range_freq_dynamic():
        """测试动态区间频率查询类"""
        print("开始测试 RangeFreqDynamic...")

        # 测试1: 基本功能测试
        print("测试1: 基本功能")
        nums = [1, 2, 3, 2, 4, 2, 5]
        rf = RangeFreqDynamic(nums)

        # 测试初始状态的查询
        assert rf.query(0, 7, 2) == 3, "元素2在整个数组中出现3次"
        assert rf.query(1, 5, 2) == 2, "元素2在区间[1,5)中出现2次"
        assert rf.query(0, 3, 2) == 1, "元素2在区间[0,3)中出现1次"
        assert rf.query(0, 7, 6) == 0, "元素6不存在"

        # 测试get方法
        assert rf.get(0) == 1, "索引0的值应该是1"
        assert rf.get(3) == 2, "索引3的值应该是2"

        print("✓ 基本功能测试通过")

        # 测试2: 更新操作
        print("测试2: 更新操作")

        # 更新索引1的值从2改为10
        assert rf.set(1, 10) == True, "应该成功更新"
        assert rf.get(1) == 10, "索引1的值应该是10"
        assert rf.query(0, 7, 2) == 2, "更新后元素2只出现2次"
        assert rf.query(0, 7, 10) == 1, "元素10出现1次"

        # 尝试设置相同的值
        assert rf.set(1, 10) == False, "设置相同值应该返回False"

        # 恢复原值
        rf.set(1, 2)
        assert rf.query(0, 7, 2) == 3, "恢复后元素2又出现3次"

        print("✓ 更新操作测试通过")

        # 测试3: findFirst和findLast方法
        print("测试3: findFirst和findLast方法")

        # 测试findFirst
        assert rf.findFirst(0, 7, 2) == 1, "元素2的第一个位置是1"
        assert rf.findFirst(2, 7, 2) == 3, "在区间[2,7)中元素2的第一个位置是3"
        assert rf.findFirst(4, 7, 2) == 5, "在区间[4,7)中元素2的第一个位置是5"
        assert rf.findFirst(0, 7, 6) == -1, "元素6不存在"
        assert rf.findFirst(6, 7, 2) == -1, "区间[6,7)中没有元素2"

        # 测试findLast
        assert rf.findLast(0, 7, 2) == 5, "元素2的最后一个位置是5"
        assert rf.findLast(0, 4, 2) == 3, "在区间[0,4)中元素2的最后一个位置是3"
        assert rf.findLast(0, 2, 2) == 1, "在区间[0,2)中元素2的最后一个位置是1"
        assert rf.findLast(0, 7, 6) == -1, "元素6不存在"
        assert rf.findLast(0, 1, 2) == -1, "区间[0,1)中没有元素2"

        print("✓ findFirst和findLast测试通过")

        # 测试4: 边界情况
        print("测试4: 边界情况")

        # 空区间
        assert rf.query(3, 3, 2) == 0, "空区间"
        assert rf.findFirst(3, 3, 2) == -1, "空区间"
        assert rf.findLast(3, 3, 2) == -1, "空区间"

        # 单元素区间
        assert rf.query(1, 2, 2) == 1, "单元素区间包含目标值"
        assert rf.query(0, 1, 2) == 0, "单元素区间不包含目标值"

        print("✓ 边界情况测试通过")

        # 测试5: 复杂更新场景
        print("测试5: 复杂更新场景")

        # 创建新的测试数据
        rf2 = RangeFreqDynamic([1, 1, 1, 2, 2, 3])

        # 初始状态
        assert rf2.query(0, 6, 1) == 3, "初始状态元素1出现3次"
        assert rf2.query(0, 6, 2) == 2, "初始状态元素2出现2次"

        # 将所有1改为4
        rf2.set(0, 4)
        rf2.set(1, 4)
        rf2.set(2, 4)

        assert rf2.query(0, 6, 1) == 0, "元素1不再存在"
        assert rf2.query(0, 6, 4) == 3, "元素4出现3次"
        assert rf2.findFirst(0, 6, 4) == 0, "元素4的第一个位置是0"
        assert rf2.findLast(0, 6, 4) == 2, "元素4的最后一个位置是2"

        print("✓ 复杂更新场景测试通过")

        # 测试6: 字符串类型
        print("测试6: 字符串类型")

        str_nums = ["a", "b", "c", "b", "d", "b"]
        str_rf = RangeFreqDynamic(str_nums)

        assert str_rf.query(0, 6, "b") == 3, "字符串b出现3次"
        assert str_rf.findFirst(0, 6, "b") == 1, "字符串b的第一个位置是1"
        assert str_rf.findLast(0, 6, "b") == 5, "字符串b的最后一个位置是5"

        # 更新字符串
        str_rf.set(1, "x")
        assert str_rf.query(0, 6, "b") == 2, "更新后字符串b出现2次"
        assert str_rf.query(0, 6, "x") == 1, "字符串x出现1次"

        print("✓ 字符串类型测试通过")

        # 测试7: 大规模随机测试
        print("测试7: 大规模测试")

        import random

        random.seed(42)

        # 创建随机数据
        size = 100
        data = [random.randint(1, 10) for _ in range(size)]
        rf_large = RangeFreqDynamic(data)

        # 执行随机查询和更新
        for _ in range(50):
            if random.random() < 0.7:  # 70% 查询
                start = random.randint(0, size - 1)
                end = random.randint(start + 1, size + 1)
                value = random.randint(1, 10)
                count = rf_large.query(start, end, value)
                assert count >= 0, "查询结果应该非负"
            else:  # 30% 更新
                index = random.randint(0, size - 1)
                new_value = random.randint(1, 10)
                rf_large.set(index, new_value)

        print("✓ 大规模测试通过")

        print("所有测试通过! ✅")

    test_range_freq_dynamic()

    # 此外，给定一个二维整数数组 queries，其中每个 queries[i] 可以是以下两种类型之一：
    # [1, l, r, k] - 计算在区间 [l, r] 中，满足 nums[j] 的 popcount-depth 等于 k 的索引 j 的数量。
    # [2, idx, val] - 将 nums[idx] 更新为 val。
    # 返回一个整数数组 answer，其中 answer[i] 表示第 i 个类型为 [1, l, r, k] 的查询的结果。
    class Solution:
        # 3624. 位计数深度为 K 的整数数目 II
        # https://leetcode.cn/problems/number-of-integers-with-popcount-depth-equal-to-k-ii/
        def popcountDepth(self, nums: List[int], queries: List[List[int]]) -> List[int]:
            def depth(x: int) -> int:
                res = 0
                while x > 1:
                    res += 1
                    x = x.bit_count()
                return res

            R = RangeFreqDynamic([depth(v) for v in nums])
            res = []
            for q in queries:
                if q[0] == 1:
                    res.append(R.query(q[1], q[2] + 1, q[3]))
                else:
                    R.set(q[1], depth(q[2]))
            return res
