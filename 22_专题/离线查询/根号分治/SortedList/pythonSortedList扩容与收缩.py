def _expand(self, pos):
    """Split sublists with length greater than double the load-factor.

    Updates the index when the sublist length is less than double the load
    level. This requires incrementing the nodes in a traversal from the
    leaf node to the root. For an example traversal see
    ``SortedList._loc``.

    """
    _load = self._load
    _lists = self._lists
    _index = self._index

    if len(_lists[pos]) > (_load << 1):  # 如果子列表长度大于两倍的负载因子：分割子列表
        _maxes = self._maxes

        _lists_pos = _lists[pos]
        half = _lists_pos[_load:]
        del _lists_pos[_load:]
        _maxes[pos] = _lists_pos[-1]

        _lists.insert(pos + 1, half)
        _maxes.insert(pos + 1, half[-1])

        del _index[:]  # 重置线段树
    else:
        # 如果子列表长度小于等于两倍的负载因子：更新索引中对应节点的计数值，表示增加了一个元素
        # seg.update(pos,1)
        if _index:
            child = self._offset + pos
            while child:
                _index[child] += 1
                child = (child - 1) >> 1
            _index[0] += 1


def _delete(self, pos, idx):
    """Delete value at the given `(pos, idx)`.

    Combines lists that are less than half the load level.

    Updates the index when the sublist length is more than half the load
    level. This requires decrementing the nodes in a traversal from the
    leaf node to the root. For an example traversal see
    ``SortedList._loc``.

    :param int pos: lists index
    :param int idx: sublist index

    """
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


def _expand(self, pos):
    """Split sublists with length greater than double the load-factor.

    Updates the index when the sublist length is less than double the load
    level. This requires incrementing the nodes in a traversal from the
    leaf node to the root. For an example traversal see
    ``SortedList._loc``.

    """
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


def _delete(self, pos, idx):
    """Delete value at the given `(pos, idx)`.

    Combines lists that are less than half the load level.

    Updates the index when the sublist length is more than half the load
    level. This requires decrementing the nodes in a traversal from the
    leaf node to the root. For an example traversal see
    ``SortedList._loc``.

    :param int pos: lists index
    :param int idx: sublist index

    """
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


def discard(self, value):
    """Remove `value` from sorted list if it is a member.

    If `value` is not a member, do nothing.

    Runtime complexity: `O(log(n))` -- approximate.

    >>> sl = SortedList([1, 2, 3, 4, 5])
    >>> sl.discard(5)
    >>> sl.discard(0)
    >>> sl == [1, 2, 3, 4]
    True

    :param value: `value` to discard from sorted list

    """
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


# del _index[:]

nums = [1, 2, 3, 4]
del nums[:]
print(nums)
