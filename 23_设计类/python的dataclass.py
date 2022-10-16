# 方便定义interface
# 例如线段树结点的定义

# 组合的对象不可以使用插槽


# order/slot/fronzen
from dataclasses import dataclass, field
from typing import Any, Callable, List, Optional


Callback = Callable[["TreeNode"], None]


@dataclass(slots=True, unsafe_hash=True, frozen=True)
class TreeNode:
    val: int = 0
    left: Optional["TreeNode"] = None
    right: Optional["TreeNode"] = None
    # metaInfo: List[Any] = field(default_factory=list)

    def walk(self, callback: Callback) -> None:
        if self.left:
            self.left.walk(callback)
        callback(self)
        if self.right:
            self.right.walk(callback)


root = TreeNode(1, TreeNode(2, TreeNode(4), TreeNode(5)), TreeNode(3, TreeNode(6), TreeNode(7)))

# dataclass 会自动添加 __repr__  __eq__  __hash__ 方法
# 相比普通class，dataclass通常不包含私有属性，数据可以直接访问
# dataclass通常情况下是unhashable的，因为默认生成的`__hash__`是`None`，所以不能用来做字典的key，
# 如果有这种需求，那么应该指定你的数据类为frozen或者指定unsafe_hash=True
# Frozen 实例是在初始化对象后无法修改其属性的对象。
# root.val = 2
root.walk(print)
print(root.__hash__())
