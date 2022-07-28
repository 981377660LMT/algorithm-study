"""Sorted List Used in atcoder

Dont restrict type of `T` by using `SupportsComparison Protocol`
which is only available when python version is 3.8 or later and not supported in pypy3

为了保持性能,原版的SortedList在初始化时会判断是否传了key
如果传了key,`__new__`工厂会返回 `SortedKeyList(SortedList)` 每个方法里都会调用这个key
否则直接使用 SortedList 减小调用开销

这里的SortedList不支持传key 需要自定义大小时需要自己维护传入的value的比较方法
"""

from bisect import bisect_left, bisect_right, insort
from collections.abc import Sequence
from functools import reduce
from itertools import chain, repeat, starmap
from math import log
from operator import add, eq, ge, gt, iadd, le, lt, ne
from textwrap import dedent
from typing import Generic, Iterable, Optional, Tuple, TypeVar

T = TypeVar("T")


class SortedList(Generic[T]):
    DEFAULT_LOAD_FACTOR = 1000

    def __init__(self, iterable: Optional[Iterable[T]] = None):
        self._len = 0
        self._load = self.DEFAULT_LOAD_FACTOR
        self._lists = []  # 分块
        self._maxes = []  # 各个分块的最大值
        self._index = []
        self._offset = 0
        if iterable is not None:
            self._update(iterable)

    def clear(self) -> None:
        self._len = 0
        del self._lists[:]
        del self._maxes[:]
        del self._index[:]
        self._offset = 0

    _clear = clear

    def add(self, value: T) -> None:
        _lists = self._lists
        _maxes = self._maxes
        if _maxes:
            pos = bisect_right(_maxes, value)
            if pos == len(_maxes):
                pos -= 1
                _lists[pos].append(value)
                _maxes[pos] = value
            else:
                insort(_lists[pos], value)
            self._expand(pos)
        else:
            _lists.append([value])
            _maxes.append(value)
        self._len += 1

    def _expand(self, pos: int) -> None:
        _load = self._load
        _lists = self._lists
        _index = self._index
        if len(_lists[pos]) > (_load << 1):
            _maxes = self._maxes
            _lists_pos = _lists[pos]
            half = _lists_pos[_load:]
            del _lists_pos[_load:]
            _maxes[pos] = _lists_pos[-1]
            _lists.insert(pos + 1, half)
            _maxes.insert(pos + 1, half[-1])
            del _index[:]
        else:
            if _index:
                child = self._offset + pos
                while child:
                    _index[child] += 1
                    child = (child - 1) >> 1
                _index[0] += 1

    def update(self, iterable: Iterable[T]) -> None:
        _lists = self._lists
        _maxes = self._maxes
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
        _load = self._load
        _lists.extend(values[pos : (pos + _load)] for pos in range(0, len(values), _load))
        _maxes.extend(sublist[-1] for sublist in _lists)
        self._len = len(values)
        del self._index[:]

    _update = update

    def __contains__(self, value: T) -> bool:
        _maxes = self._maxes
        if not _maxes:
            return False
        pos = bisect_left(_maxes, value)
        if pos == len(_maxes):
            return False
        _lists = self._lists
        idx = bisect_left(_lists[pos], value)
        return _lists[pos][idx] == value

    def remove(self, value: T) -> None:
        _maxes = self._maxes
        if not _maxes:
            raise ValueError("{0!r} not in list".format(value))
        pos = bisect_left(_maxes, value)
        if pos == len(_maxes):
            raise ValueError("{0!r} not in list".format(value))
        _lists = self._lists
        idx = bisect_left(_lists[pos], value)
        if _lists[pos][idx] == value:
            self._delete(pos, idx)
        else:
            raise ValueError("{0!r} not in list".format(value))

    def discard(self, value: T) -> None:
        _maxes = self._maxes
        if not _maxes:
            return
        pos = bisect_left(_maxes, value)
        if pos == len(_maxes):
            return
        _lists = self._lists
        idx = bisect_left(_lists[pos], value)
        if _lists[pos][idx] == value:
            self._delete(pos, idx)

    def _delete(self, pos: int, idx: int) -> None:
        _lists = self._lists
        _maxes = self._maxes
        _index = self._index
        _lists_pos = _lists[pos]
        del _lists_pos[idx]
        self._len -= 1
        len_lists_pos = len(_lists_pos)
        if len_lists_pos > (self._load >> 1):
            _maxes[pos] = _lists_pos[-1]
            if _index:
                child = self._offset + pos
                while child > 0:
                    _index[child] -= 1
                    child = (child - 1) >> 1
                _index[0] -= 1
        elif len(_lists) > 1:
            if not pos:
                pos += 1
            prev = pos - 1
            _lists[prev].extend(_lists[pos])
            _maxes[prev] = _lists[prev][-1]
            del _lists[pos]
            del _maxes[pos]
            del _index[:]
            self._expand(prev)
        elif len_lists_pos:
            _maxes[pos] = _lists_pos[-1]
        else:
            del _lists[pos]
            del _maxes[pos]
            del _index[:]

    def _loc(self, pos: int, idx: int) -> int:
        if not pos:
            return idx
        _index = self._index
        if not _index:
            self._build_index()
        total = 0
        pos += self._offset
        while pos:
            if not pos & 1:
                total += _index[pos - 1]
            pos = (pos - 1) >> 1
        return total + idx

    def _pos(self, idx: int) -> Tuple[int, int]:
        if idx < 0:
            last_len = len(self._lists[-1])
            if (-idx) <= last_len:
                return len(self._lists) - 1, last_len + idx
            idx += self._len
            if idx < 0:
                raise IndexError("list index out of range")
        elif idx >= self._len:
            raise IndexError("list index out of range")
        if idx < len(self._lists[0]):
            return 0, idx
        _index = self._index
        if not _index:
            self._build_index()
        pos = 0
        child = 1
        len_index = len(_index)
        while child < len_index:
            index_child = _index[child]
            if idx < index_child:
                pos = child
            else:
                idx -= index_child
                pos = child + 1
            child = (pos << 1) + 1
        return (pos - self._offset, idx)

    def _build_index(self) -> None:
        row0 = list(map(len, self._lists))
        if len(row0) == 1:
            self._index[:] = row0
            self._offset = 0
            return
        head = iter(row0)
        tail = iter(head)
        row1 = list(starmap(add, zip(head, tail)))
        if len(row0) & 1:
            row1.append(row0[-1])
        if len(row1) == 1:
            self._index[:] = row1 + row0
            self._offset = 1
            return
        size = 2 ** (int(log(len(row1) - 1, 2)) + 1)
        row1.extend(repeat(0, size - len(row1)))
        tree = [row0, row1]
        while len(tree[-1]) > 1:
            head = iter(tree[-1])
            tail = iter(head)
            row = list(starmap(add, zip(head, tail)))
            tree.append(row)
        reduce(iadd, reversed(tree), self._index)
        self._offset = size * 2 - 1

    def __delitem__(self, index: int) -> None:
        if isinstance(index, slice):
            start, stop, step = index.indices(self._len)
            if step == 1 and start < stop:
                if start == 0 and stop == self._len:
                    return self._clear()
                elif self._len <= 8 * (stop - start):
                    values = self._getitem(slice(None, start))
                    if stop < self._len:
                        values += self._getitem(slice(stop, None))
                    self._clear()
                    return self._update(values)
            indices = range(start, stop, step)
            if step > 0:
                indices = reversed(indices)
            _pos, _delete = self._pos, self._delete
            for index in indices:
                pos, idx = _pos(index)
                _delete(pos, idx)
        else:
            pos, idx = self._pos(index)
            self._delete(pos, idx)

    def __getitem__(self, index: int):
        _lists = self._lists
        if isinstance(index, slice):
            start, stop, step = index.indices(self._len)
            if step == 1 and start < stop:
                if start == 0 and stop == self._len:
                    return reduce(iadd, self._lists, [])
                start_pos, start_idx = self._pos(start)
                start_list = _lists[start_pos]
                stop_idx = start_idx + stop - start
                if len(start_list) >= stop_idx:
                    return start_list[start_idx:stop_idx]
                if stop == self._len:
                    stop_pos = len(_lists) - 1
                    stop_idx = len(_lists[stop_pos])
                else:
                    stop_pos, stop_idx = self._pos(stop)
                prefix = _lists[start_pos][start_idx:]
                middle = _lists[(start_pos + 1) : stop_pos]
                result = reduce(iadd, middle, prefix)
                result += _lists[stop_pos][:stop_idx]
                return result
            if step == -1 and start > stop:
                result = self._getitem(slice(stop + 1, start + 1))
                result.reverse()
                return result
            indices = range(start, stop, step)
            return list(self._getitem(index) for index in indices)
        else:
            if self._len:
                if index == 0:
                    return _lists[0][0]
                elif index == -1:
                    return _lists[-1][-1]
            else:
                raise IndexError("list index out of range")
            if 0 <= index < len(_lists[0]):
                return _lists[0][index]
            len_last = len(_lists[-1])
            if -len_last < index < 0:
                return _lists[-1][len_last + index]
            pos, idx = self._pos(index)
            return _lists[pos][idx]

    _getitem = __getitem__

    def __setitem__(self, index, value):
        message = "use ``del sl[index]`` and ``sl.add(value)`` instead"
        raise NotImplementedError(message)

    def __iter__(self):
        return chain.from_iterable(self._lists)

    def __reversed__(self):
        return chain.from_iterable(map(reversed, reversed(self._lists)))

    def islice(self, start: Optional[T] = None, stop: Optional[T] = None, reverse=False):
        _len = self._len
        if not _len:
            return iter(())
        start, stop, _ = slice(start, stop).indices(self._len)
        if start >= stop:
            return iter(())
        _pos = self._pos
        min_pos, min_idx = _pos(start)
        if stop == _len:
            max_pos = len(self._lists) - 1
            max_idx = len(self._lists[-1])
        else:
            max_pos, max_idx = _pos(stop)
        return self._islice(min_pos, min_idx, max_pos, max_idx, reverse)

    def _islice(self, min_pos, min_idx, max_pos, max_idx, reverse):
        _lists = self._lists
        if min_pos > max_pos:
            return iter(())
        if min_pos == max_pos:
            if reverse:
                indices = reversed(range(min_idx, max_idx))
                return map(_lists[min_pos].__getitem__, indices)
            indices = range(min_idx, max_idx)
            return map(_lists[min_pos].__getitem__, indices)
        next_pos = min_pos + 1
        if next_pos == max_pos:
            if reverse:
                min_indices = range(min_idx, len(_lists[min_pos]))
                max_indices = range(max_idx)
                return chain(
                    map(_lists[max_pos].__getitem__, reversed(max_indices)),
                    map(_lists[min_pos].__getitem__, reversed(min_indices)),
                )
            min_indices = range(min_idx, len(_lists[min_pos]))
            max_indices = range(max_idx)
            return chain(
                map(_lists[min_pos].__getitem__, min_indices),
                map(_lists[max_pos].__getitem__, max_indices),
            )
        if reverse:
            min_indices = range(min_idx, len(_lists[min_pos]))
            sublist_indices = range(next_pos, max_pos)
            sublists = map(_lists.__getitem__, reversed(sublist_indices))
            max_indices = range(max_idx)
            return chain(
                map(_lists[max_pos].__getitem__, reversed(max_indices)),
                chain.from_iterable(map(reversed, sublists)),
                map(_lists[min_pos].__getitem__, reversed(min_indices)),
            )
        min_indices = range(min_idx, len(_lists[min_pos]))
        sublist_indices = range(next_pos, max_pos)
        sublists = map(_lists.__getitem__, sublist_indices)
        max_indices = range(max_idx)
        return chain(
            map(_lists[min_pos].__getitem__, min_indices),
            chain.from_iterable(sublists),
            map(_lists[max_pos].__getitem__, max_indices),
        )

    def irange(
        self,
        minimum: Optional[T] = None,
        maximum: Optional[T] = None,
        inclusive=(True, True),
        reverse=False,
    ):
        """Create an iterator of values between `minimum` and `maximum`.
        >>> sl = SortedList('abcdefghij')
        >>> it = sl.irange('c', 'f')
        >>> list(it)
        ['c', 'd', 'e', 'f']
        """
        _maxes = self._maxes
        if not _maxes:
            return iter(())
        _lists = self._lists
        if minimum is None:
            min_pos = 0
            min_idx = 0
        else:
            if inclusive[0]:
                min_pos = bisect_left(_maxes, minimum)
                if min_pos == len(_maxes):
                    return iter(())
                min_idx = bisect_left(_lists[min_pos], minimum)
            else:
                min_pos = bisect_right(_maxes, minimum)
                if min_pos == len(_maxes):
                    return iter(())
                min_idx = bisect_right(_lists[min_pos], minimum)
        if maximum is None:
            max_pos = len(_maxes) - 1
            max_idx = len(_lists[max_pos])
        else:
            if inclusive[1]:
                max_pos = bisect_right(_maxes, maximum)
                if max_pos == len(_maxes):
                    max_pos -= 1
                    max_idx = len(_lists[max_pos])
                else:
                    max_idx = bisect_right(_lists[max_pos], maximum)
            else:
                max_pos = bisect_left(_maxes, maximum)
                if max_pos == len(_maxes):
                    max_pos -= 1
                    max_idx = len(_lists[max_pos])
                else:
                    max_idx = bisect_left(_lists[max_pos], maximum)
        return self._islice(min_pos, min_idx, max_pos, max_idx, reverse)

    def __len__(self):
        return self._len

    def bisect_left(self, value: T) -> int:
        _maxes = self._maxes
        if not _maxes:
            return 0
        pos = bisect_left(_maxes, value)
        if pos == len(_maxes):
            return self._len
        idx = bisect_left(self._lists[pos], value)
        return self._loc(pos, idx)

    def bisect_right(self, value: T) -> int:
        _maxes = self._maxes
        if not _maxes:
            return 0
        pos = bisect_right(_maxes, value)
        if pos == len(_maxes):
            return self._len
        idx = bisect_right(self._lists[pos], value)
        return self._loc(pos, idx)

    bisect = bisect_right
    _bisect_right = bisect_right

    def count(self, value: T) -> int:
        _maxes = self._maxes
        if not _maxes:
            return 0
        pos_left = bisect_left(_maxes, value)
        if pos_left == len(_maxes):
            return 0
        _lists = self._lists
        idx_left = bisect_left(_lists[pos_left], value)
        pos_right = bisect_right(_maxes, value)
        if pos_right == len(_maxes):
            return self._len - self._loc(pos_left, idx_left)
        idx_right = bisect_right(_lists[pos_right], value)
        if pos_left == pos_right:
            return idx_right - idx_left
        right = self._loc(pos_right, idx_right)
        left = self._loc(pos_left, idx_left)
        return right - left

    def copy(self):
        return self.__class__(self)

    __copy__ = copy

    def pop(self, index=-1) -> T:
        if not self._len:
            raise IndexError("pop index out of range")
        _lists = self._lists
        if index == 0:
            val = _lists[0][0]
            self._delete(0, 0)
            return val
        if index == -1:
            pos = len(_lists) - 1
            loc = len(_lists[pos]) - 1
            val = _lists[pos][loc]
            self._delete(pos, loc)
            return val
        if 0 <= index < len(_lists[0]):
            val = _lists[0][index]
            self._delete(0, index)
            return val
        len_last = len(_lists[-1])
        if -len_last < index < 0:
            pos = len(_lists) - 1
            loc = len_last + index
            val = _lists[pos][loc]
            self._delete(pos, loc)
            return val
        pos, idx = self._pos(index)
        val = _lists[pos][idx]
        self._delete(pos, idx)
        return val

    def index(self, value: T, start: Optional[int] = None, stop: Optional[int] = None) -> int:
        _len = self._len
        if not _len:
            raise ValueError("{0!r} is not in list".format(value))
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
        _maxes = self._maxes
        pos_left = bisect_left(_maxes, value)
        if pos_left == len(_maxes):
            raise ValueError("{0!r} is not in list".format(value))
        _lists = self._lists
        idx_left = bisect_left(_lists[pos_left], value)
        if _lists[pos_left][idx_left] != value:
            raise ValueError("{0!r} is not in list".format(value))
        stop -= 1
        left = self._loc(pos_left, idx_left)
        if start <= left:
            if left <= stop:
                return left
        else:
            right = self._bisect_right(value) - 1
            if start <= right:
                return start
        raise ValueError("{0!r} is not in list".format(value))

    def __add__(self, other):
        values = reduce(iadd, self._lists, [])
        values.extend(other)
        return self.__class__(values)

    __radd__ = __add__

    def __iadd__(self, other):
        self._update(other)
        return self

    def __mul__(self, num):
        values = reduce(iadd, self._lists, []) * num
        return self.__class__(values)

    __rmul__ = __mul__

    def __imul__(self, num):
        values = reduce(iadd, self._lists, []) * num
        self._clear()
        self._update(values)
        return self

    def __make_cmp(seq_op, symbol, doc):
        "Make comparator method."

        def comparer(self, other):
            "zip逐项比较大小,如果大小一样,再比较长度"
            if not isinstance(other, Sequence):
                return NotImplemented
            self_len = self._len
            len_other = len(other)
            if self_len != len_other:
                if seq_op is eq:
                    return False
                if seq_op is ne:
                    return True
            for alpha, beta in zip(self, other):
                if alpha != beta:
                    return seq_op(alpha, beta)
            return seq_op(self_len, len_other)

        seq_op_name = seq_op.__name__
        comparer.__name__ = "__{0}__".format(seq_op_name)
        doc_str = """Return true if and only if sorted list is {0} `other`.
        ``sl.__{1}__(other)`` <==> ``sl {2} other``
        Comparisons use lexicographical order as with sequences.
        Runtime complexity: `O(n)`
        :param other: `other` sequence
        :return: true if sorted list is {0} `other`
        """
        comparer.__doc__ = dedent(doc_str.format(doc, seq_op_name, symbol))
        return comparer

    __eq__ = __make_cmp(eq, "==", "equal to")
    __ne__ = __make_cmp(ne, "!=", "not equal to")
    __lt__ = __make_cmp(lt, "<", "less than")
    __gt__ = __make_cmp(gt, ">", "greater than")
    __le__ = __make_cmp(le, "<=", "less than or equal to")
    __ge__ = __make_cmp(ge, ">=", "greater than or equal to")
    __make_cmp = staticmethod(__make_cmp)

    def __reduce__(self):
        values = reduce(iadd, self._lists, [])
        return (type(self), (values,))

    def __repr__(self) -> str:
        """一般不会嵌套使用,所以这里没有用 `recursive_repr`"""
        return "{0}({1!r})".format(self.__class__.__name__, list(self))


if __name__ == "__main__":
    sl = SortedList([1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
    assert sl.bisect_left(5) == 4
    assert sl.bisect_right(5) == 5
    assert sl[4] == 5
    sl.add(11)
    sl.remove(11)
    assert 11 not in sl
    sl.discard(11)
    assert sl.index(5) == 4
    assert list(sl.islice(2, 6)) == [3, 4, 5, 6]
    sl.pop(2)

    print(sl)
