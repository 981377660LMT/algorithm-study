# SortedList 源码简析 (python)

[toc]

---

原文地址：https://github.com/981377660LMT/ts/issues/267
添加注释后的源码：
https://github.com/981377660LMT/python-sortedcontainers

随手写的一些笔记，~存在大量翻译腔~，希望对理解 SortedList 源码有所帮助

## 前言：S 代表着什么？新一代有序容器 SortedList

与基于树的实现相比，使用列表具有一些基于内存使用情况的优势。

- 大多数插入/删除不需要分配或释放内存。
- 指向元素的指针密集排列。这有利于硬件的内存架构，并更有效地利用缓存。

## 代码整体关系

SortedList 根据传不传 key 分为 SortedList、SortedKeyList，不传 key 时降低了函数调用的开销。这里 [`__new__`](https://github.com/grantjenks/python-sortedcontainers/blob/92ef500158f87f7684023823d689cfd7bef892a1/src/sortedcontainers/sortedlist.py#L174) 为工厂函数。
SortedDict 和 SortedSet 基于 SortedList 封装实现。

## 出发点：短列表的快速插入删除

传统的基于树的设计具有更好的理论复杂度，但这忽略了当今软件和硬件的现实。
python 的 list 很快，有利于 memory management 和 random access.
bisect.insort 很快.(其实 js 的 Array.prototype.splice 也很快).
在不太长(1000 到 2000)的 list 上做插入删除操作，常数极小.
在硬件和软件方面，现代处理器已经花费了大量时间来优化类似 mem-cpy/mem-move 的操作。
`因此，在这里把插入删除当作了 O(1)`
但是长度达到阈值，例如超过 1e4 时，bisect.insort 插入会很慢。
因此需要分块。
插入或删除最常在短的 list 中执行。很少需要添加或删除新列表。

## 实现细节 1:分块+线段树

1. 两种查询：

   - getRankByValue -> 查询块内最大值 -> `_maxes`
   - getValueByRank -> 查询位置 -> `_index+_offset`.维护这个更难一些.
     也可以不要求 O(logn)，而是将块的 size 设置大一点，查找 pos 时直接遍历.(现在 js 第一版的 SortedList 就是这样的;但是不用维护线段树，占用空间少一些.)

2. 三个关键：`_maxes`, `_lists`, `_index`+`_offset`

   - `_lists` 是每个块内部的元素(有序)
   - `_maxes` 是每个块的最大值(用于二分定位).其实时冗余信息，可以通过 `_lists[i][-1]` 得到.这样做的好处是避免二级访问，降低 access 的开销.

   - `_index` 一颗线段树，叶子结点维护每个块的长度.
     这个信息对于插入和删除操作非常有用，因为它提供了一种快速定位某个位置在整个有序列表中的准确位置，加速定位操作。
     作者把块的索引统一命名为`pos`，把块内的索引统一命名为`idx`。
     `_loc(pos,idx)->int`: 树上二分，根据(块的索引，块内的索引)返回元素在整个列表中的索引
     `_pos(idx)->Tuple[int,int]`: 树上二分，返回 idx 对应的(块的索引，块内的索引)
     https://github.com/grantjenks/python-sortedcontainers/blob/92ef500158f87f7684023823d689cfd7bef892a1/src/sortedcontainers/sortedlist.py#L520C19-L520C19
     `_offset` 是第叶子节点的起始索引

3. Q&A:

- Q:为什么不用前缀和+二分来定位
  A:因为是动态的，更新前缀和需要整个更新 O(n)，而更新线段树只需要 O(logn)
- Q:为什么不用树状数组+二分来定位
  A:类似的原因，更新复杂度太大，难以做到 O(logn) 维护

## 实现细节 2：动态扩容与收缩策略

`_expand`函数与`_delete`函数
https://github.com/grantjenks/python-sortedcontainers/blob/92ef500158f87f7684023823d689cfd7bef892a1/src/sortedcontainers/sortedlist.py#L289
https://github.com/grantjenks/python-sortedcontainers/blob/92ef500158f87f7684023823d689cfd7bef892a1/src/sortedcontainers/sortedlist.py#L465
增加和删除元素时还需要根据负载调整`_lists` 的子列表
如果子列表的长度超过负载`(DEFAULT_LOAD_FACTOR)`的两倍，则将其一分为二
如果减少到负载的一半，则与邻居合并
系数为 1000 时，可以支持到 1e7 到 1e8 级别的数据.
作者建议设置立方根级别.

> See :doc:`implementation` and :doc:`performance-scale` for more information.

## 源码中的技巧

- 缓存类变量/类方法来减少属性查找的开销(大量使用)

```python
_lists = self._lists
_maxes = self._maxes
_add = self.add
```

- update 批量添加元素的处理

如果添加的元素太多`if len(values) * 4 >= self._len:`,
全部重构一遍，否则直接添加到`_lists`中.
这里的 4 可能是指线段树最坏时需要的 4 倍空间.

```python
values = sorted(iterable)

if _maxes:
    if len(values) * 4 >= self._len:
        _lists.append(values)
        values = reduce(iadd, _lists, [])
        values.sort()
        self._clear()
    else:
        _add = self.add
        for val in values:
            _add(val)
        return
```

- 重写 list 的几个 方法，抛出错误

```python
    def append(self, value):
    """Raise not-implemented error.

    Implemented to override `MutableSequence.append` which provides an
    erroneous default implementation.

    :raises NotImplementedError: use ``sl.add(value)`` instead

    """
    raise NotImplementedError("use ``sl.add(value)`` instead")
```

- `__getitem__`和`pop` 方法对删除开头和删除结尾做了特殊判断的优化,这两种场合是 O(1)的

- 对边界和负索引的处理

  ```python
  def index(self, value, start=None, stop=None):
     ...

     if start is None:
         start = 0
     if start < 0:
         start += _len
     if start < 0:
         start = 0

     if stop is None:
         stop = _len
     if stop < 0:
         stop += _len
     if stop > _len:
         stop = _len

     if stop <= start:
         raise ValueError("{0!r} is not in list".format(value))
  ```

- View 利用 proxy 实现了延迟求值.

```python

class SortedItemsView(ItemsView, Sequence):
   """Sorted items view is a dynamic view of the sorted dict's items.
   When the sorted dict's items change, the view reflects those changes.
   The items view implements the set and sequence abstract base classes.
   """
   __slots__ = ()

   @classmethod
   def _from_iterable(cls, it):
       return SortedSet(it)


   def __getitem__(self, index):
       _mapping = self._mapping
       _mapping_list = _mapping._list

       if isinstance(index, slice):
           keys = _mapping_list[index]
           return [(key, _mapping[key]) for key in keys]

       key = _mapping_list[index]
       return key, _mapping[key]


   __delitem__ = _view_delitem
```

## 结尾

看到这里，不由感叹 SortedList 的作者 [Grant Jenks](https://github.com/grantjenks) 大叔功力之深厚，创作出了这么优秀的模板。
最后，我想以日本著名游戏制作人加藤·惠（1998 ~) 的一句名言结束本文：
![image.png](https://pic.leetcode.cn/1687014909-IhSDld-image.png){:style="width:200px":align=center}

## 参考

https://www.zhihu.com/question/593450942/answer/2966795898
https://grantjenks.com/docs/sortedcontainers/implementation.html
