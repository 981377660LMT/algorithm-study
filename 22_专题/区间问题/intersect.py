# intersect 函数用于计算两个范围的交集，并返回一个新的范围对象。
# 这个函数可以帮助我们在实现文本编辑功能时，计算出选中范围和操作范围的交集，从而只操作交集部分的文本。

# relativeComplement 函数用于计算一个范围相对于另一个范围的差集，并返回一个范围对象数组。
# 这个函数可以帮助我们在实现文本编辑功能时，计算出要操作的范围和要排除的范围之间的差集，
# 从而只操作差集部分的文本。

# 这两个函数的实现非常简单，但在实际的编辑器操作中非常实用。
# 掌握这两个函数的使用，可以使得我们在理解VSCode关于列表布局、编辑器相关核心代码的时候更加清晰。


"""
vscode Range 模块中的两个函数
"""


from typing import Optional, Tuple, List


Interval = Tuple[int, int]


def intersect(interval1: "Interval", interval2: "Interval") -> Optional["Interval"]:
    """获取两个闭区间交集"""
    if interval1[0] > interval2[1] or interval2[0] > interval1[1]:
        return None
    return (max(interval1[0], interval2[0]), min(interval1[1], interval2[1]))


def relativeComplement(one: "Interval", other: "Interval") -> List["Interval"]:
    """两个闭区间的相对差集.从one中移除other."""
    res = []
    first = [one[0], min(one[1], other[0])]
    second = [max(one[0], other[1]), one[1]]

    # not empty
    if first[0] < first[1]:
        res.append(first)
    if second[0] < second[1]:
        res.append(second)
    return res


if __name__ == "__main__":
    range1, range2 = (1, 5), (3, 7)
    print(intersect(range1, range2))
    print(relativeComplement(range1, range2))
