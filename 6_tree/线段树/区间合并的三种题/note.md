# 左端点记右端点,右端点记左端点

# 合并左右/左/右/不合并

线段树、区间并查集、区间操纵

统一解法：`map+multiset` 即 `SortedDict+SortedList`

模板参考
`2213. 由单个字符重复的最长子字符串-线段树单点更新结点实现AC`
`2158. Amount of New Area Painted Each Day-node实现`
`2213. 由单个字符重复的最长子字符串-线段树单点更新结点实现`

五个方法:

```py
def build(self, rt: int, left: int, right: int) -> None:
初始化Node准备初始状态的query、要pushUp

def update(self, rt: int, left: int, right: int, target: str) -> None:
要pushDown、pushUp，到了最细粒度的结点更新(也实际的原数组)后打上懒标记

def query(self, rt: int, left: int, right: int) -> int:
要pushDown，到了最细粒度的结点后返回节点值

def pushUp(self, rt: int) -> None:
用子节点更新父节点状态

def pushDown(self, rt: int) -> None:
传递懒标记并更新左右结点
```

不同的线段树一般是 pushDown 和 pushUp 不同，需要想清楚`怎么用左右区间更新整个区间、整个区间更新左右区间`
