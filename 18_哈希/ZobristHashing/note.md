```python
pool = defaultdict(lambda: randint(1, (1 << 61) - 1))

随机数+异或来生成集合哈希值
随机数+和来生成区间哈希值
```

- 作用

1. 用于 O(1)判断区间是否相同
2. 用于 O(1)判断集合/多重集是否相同
3. 增/删一个元素时 O(1)重新计算哈希值
4. 计算两个集合的对称差(symmetric difference) 时 O(1)

- 实现

1. 将每个`状态`分配任意的随机数
2. 集合的哈希值为所有`状态`的随机数的异或和

---

https://hackmd.io/@tatyam-prime/r1dg9Q389
本质是`改变量只有一个`的数据结构的哈希

还可以对棋盘哈希
怎么表示?
(棋子种类,棋子颜色,棋子位置) => 6x2x64 = 768
**一个棋盘是所有棋子的哈希值.**
对这 768 个元素进行哈希即可
棋子移动时,只需要异或(删除)上棋子的旧位置和(添加)新位置的哈希值即可

---

FastHash 增强 的接口类:

```python
class FastHashxxx:
    """快速计算哈希值的集合."""

    _poolSingleton = defaultdict(lambda: randint(1, (1 << 61) - 1))

    __slots__ = ("_xxx", "_hash")

    def __init__(self) -> None:
        self._set = xxx
        self._hash = 0

    def add(self, x: int) -> None:
        ...

    def discard(self, x: int) -> bool:
        ...

    def getHash(self) -> int:
        return self._hash

    def symmetricDifference(self, other: "FastHashxxx") -> int:
        return self._hash ^ other._hash

    def __hash__(self) -> int:
        return self._hash
```
