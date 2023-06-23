# 执行此操作的算法的详细信息包含在SortedList._loc和SortedList._pos的文档字符串中。

# 这两段代码是 Sorted List 数据结构中的两个辅助方法，用于在有序列表中进行索引转换和定位。

# _loc(self, pos, idx) 方法用于将索引位置转换为索引对 (lists index, sublist index)，
# 以便访问相应的子列表中的元素。该方法通过在索引树中向上遍历，
# 累计左子节点的数量，最终得到子列表中的具体位置。

# 首先，根据给定的 pos 和 idx，在索引树中定位到叶节点的位置，
# 然后根据节点的类型（左子节点或右子节点），累计左子节点的数量。
# 最后，将累计的数量和 idx 相加，得到在有序列表中的实际位置。
# !本质是线段树上二分

# _pos(self, idx) 方法用于将索引位置转换为索引对 (lists index, sublist index)，
# 以便定位到相应的子列表和具体位置。该方法通过在索引树中向下遍历，
# 根据索引的大小选择左子节点或右子节点，直到达到叶节点。

# 首先，根据给定的 idx，判断它是否超出了有序列表的范围。如果超出范围，则抛出异常。
# 然后，在索引树中向下遍历，根据当前节点的值和索引的比较结果选择左子节点或右子节点，
# 更新位置 pos 和剩余索引 idx。最后，返回计算得到的索引对 (lists index, sublist index)。


# 这两个方法在 Sorted List 中起到了索引转换和定位的作用，用于支持对有序列表的高效操作。


def _loc(self, pos, idx):
    """Convert an index pair (lists index, sublist index) into a single
    index number that corresponds to the position of the value in the
    sorted list.

    Many queries require the index be built. Details of the index are
    described in ``SortedList._build_index``.

    Indexing requires traversing the tree from a leaf node to the root. The
    parent of each node is easily computable at ``(pos - 1) // 2``.

    Left-child nodes are always at odd indices and right-child nodes are
    always at even indices.

    When traversing up from a right-child node, increment the total by the
    left-child node.

    The final index is the sum from traversal and the index in the sublist.

    For example, using the index from ``SortedList._build_index``::

        _index = 14 5 9 3 2 4 5
        _offset = 3

    Tree::

             14
          5      9
        3   2  4   5

    Converting an index pair (2, 3) into a single index involves iterating
    like so:

    1. Starting at the leaf node: offset + alpha = 3 + 2 = 5. We identify
       the node as a left-child node. At such nodes, we simply traverse to
       the parent.

    2. At node 9, position 2, we recognize the node as a right-child node
       and accumulate the left-child in our total. Total is now 5 and we
       traverse to the parent at position 0.

    3. Iteration ends at the root.

    The index is then the sum of the total and sublist index: 5 + 3 = 8.

    :param int pos: lists index
    :param int idx: sublist index
    :return: index in sorted list

    """
    if not pos:
        return idx

    _index = self._index  # 缓存

    if not _index:
        self._build_index()  # 动态开点

    total = 0

    # Increment pos to point in the index to len(self._lists[pos]).

    pos += self._offset

    # Iterate until reaching the root of the index tree at pos = 0.

    while pos:
        # Right-child nodes are at odd indices. At such indices
        # account the total below the left child node.

        if not pos & 1:
            total += _index[pos - 1]

        # Advance pos to the parent node.

        pos = (pos - 1) >> 1

    return total + idx


def _pos(self, idx):
    """Convert an index into an index pair (lists index, sublist index)
    that can be used to access the corresponding lists position.

    Many queries require the index be built. Details of the index are
    described in ``SortedList._build_index``.

    Indexing requires traversing the tree to a leaf node. Each node has two
    children which are easily computable. Given an index, pos, the
    left-child is at ``pos * 2 + 1`` and the right-child is at ``pos * 2 +
    2``.

    When the index is less than the left-child, traversal moves to the
    left sub-tree. Otherwise, the index is decremented by the left-child
    and traversal moves to the right sub-tree.

    At a child node, the indexing pair is computed from the relative
    position of the child node as compared with the offset and the remaining
    index.

    For example, using the index from ``SortedList._build_index``::

        _index = 14 5 9 3 2 4 5
        _offset = 3

    Tree::

             14
          5      9
        3   2  4   5

    Indexing position 8 involves iterating like so:

    1. Starting at the root, position 0, 8 is compared with the left-child
       node (5) which it is greater than. When greater the index is
       decremented and the position is updated to the right child node.

    2. At node 9 with index 3, we again compare the index to the left-child
       node with value 4. Because the index is the less than the left-child
       node, we simply traverse to the left.

    3. At node 4 with index 3, we recognize that we are at a leaf node and
       stop iterating.

    4. To compute the sublist index, we subtract the offset from the index
       of the leaf node: 5 - 3 = 2. To compute the index in the sublist, we
       simply use the index remaining from iteration. In this case, 3.

    The final index pair from our example is (2, 3) which corresponds to
    index 8 in the sorted list.

    :param int idx: index in sorted list
    :return: (lists index, sublist index) pair

    """
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
