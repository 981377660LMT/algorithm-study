import operator
from collections import deque
from itertools import chain, groupby, islice, repeat, tee
from typing import Iterable, List, Optional, Sized, TypeVar

T = TypeVar("T")


def chunked_even(iterable: Iterable[T], n: int):
    """Break *iterable* into lists of approximately length *n*.
    Items are distributed such the lengths of the lists differ by at most
    1 item.

    >>> iterable = [1, 2, 3, 4, 5, 6, 7]
    >>> n = 3
    >>> list(chunked_even(iterable, n))  # List lengths: 3, 2, 2
    [[1, 2, 3], [4, 5], [6, 7]]
    >>> list(chunked(iterable, n))  # List lengths: 3, 3, 1
    [[1, 2, 3], [4, 5, 6], [7]]

    """
    iterable = iter(iterable)

    # Initialize a buffer to process the chunks while keeping
    # some back to fill any underfilled chunks
    min_buffer = (n - 1) * (n - 2)
    buffer = list(islice(iterable, min_buffer))

    # Append items until we have a completed chunk
    for _ in islice(map(buffer.append, iterable), n, None, n):
        yield buffer[:n]
        del buffer[:n]

    # Check if any chunks need addition processing
    if not buffer:
        return
    length = len(buffer)

    # Chunks are either size `full_size <= n` or `partial_size = full_size - 1`
    q, r = divmod(length, n)
    num_lists = q + (1 if r > 0 else 0)
    q, r = divmod(length, num_lists)
    full_size = q + (1 if r > 0 else 0)
    partial_size = full_size - 1
    num_full = length - partial_size * num_lists

    # Yield chunks of full size
    partial_start_idx = num_full * full_size
    if full_size > 0:
        for i in range(0, partial_start_idx, full_size):
            yield buffer[i : i + full_size]

    # Yield chunks of partial size
    if partial_size > 0:
        for i in range(partial_start_idx, length, partial_size):
            yield buffer[i : i + partial_size]


def distribute(n: int, iterable: Iterable[T]):
    """Distribute the items from *iterable* among *n* smaller iterables.

        >>> group_1, group_2 = distribute(2, [1, 2, 3, 4, 5, 6])
        >>> list(group_1)
        [1, 3, 5]
        >>> list(group_2)
        [2, 4, 6]

    If the length of *iterable* is not evenly divisible by *n*, then the
    length of the returned iterables will not be identical:

        >>> children = distribute(3, [1, 2, 3, 4, 5, 6, 7])
        >>> [list(c) for c in children]
        [[1, 4, 7], [2, 5], [3, 6]]

    If the length of *iterable* is smaller than *n*, then the last returned
    iterables will be empty:

        >>> children = distribute(5, [1, 2, 3])
        >>> [list(c) for c in children]
        [[1], [2], [3], [], []]

    This function uses :func:`itertools.tee` and may require significant
    storage.

    If you need the order items in the smaller iterables to match the
    original iterable, see :func:`divide`.

    """
    if n < 1:
        raise ValueError("n must be at least 1")

    children = tee(iterable, n)
    return [islice(it, index, None, n) for index, it in enumerate(children)]


def unzip(iterable):
    """The inverse of :func:`zip`, this function disaggregates the elements
    of the zipped *iterable*.

    The ``i``-th iterable contains the ``i``-th element from each element
    of the zipped iterable. The first element is used to determine the
    length of the remaining elements.

        >>> iterable = [('a', 1), ('b', 2), ('c', 3), ('d', 4)]
        >>> letters, numbers = unzip(iterable)
        >>> list(letters)
        ['a', 'b', 'c', 'd']
        >>> list(numbers)
        [1, 2, 3, 4]

    This is similar to using ``zip(*iterable)``, but it avoids reading
    *iterable* into memory. Note, however, that this function uses
    :func:`itertools.tee` and thus may require significant storage.

    """
    head, iterable = spy(iter(iterable))
    if not head:
        # empty iterable, e.g. zip([], [], [])
        return ()
    # spy returns a one-length iterable as head
    head = head[0]
    iterables = tee(iterable, len(head))

    def itemgetter(i):
        def getter(obj):
            try:
                return obj[i]
            except IndexError:
                # basically if we have an iterable like
                # iter([(1, 2, 3), (4, 5), (6,)])
                # the second unzipped iterable would fail at the third tuple
                # since it would try to access tup[1]
                # same with the third unzipped iterable and the second tuple
                # to support these "improperly zipped" iterables,
                # we create a custom itemgetter
                # which just stops the unzipped iterables
                # at first length mismatch
                raise StopIteration

        return getter

    return tuple(map(itemgetter(i), it) for i, it in enumerate(iterables))


def spy(iterable: Iterable[T], n: int = 1):
    """Return a 2-tuple with a list containing the first *n* elements of
    *iterable*, and an iterator with the same items as *iterable*.
    This allows you to "look ahead" at the items in the iterable without
    advancing it.

    There is one item in the list by default:

        >>> iterable = 'abcdefg'
        >>> head, iterable = spy(iterable)
        >>> head
        ['a']
        >>> list(iterable)
        ['a', 'b', 'c', 'd', 'e', 'f', 'g']

    You may use unpacking to retrieve items instead of lists:

        >>> (head,), iterable = spy('abcdefg')
        >>> head
        'a'
        >>> (first, second), iterable = spy('abcdefg', 2)
        >>> first
        'a'
        >>> second
        'b'

    The number of items requested can be larger than the number of items in
    the iterable:

        >>> iterable = [1, 2, 3, 4, 5]
        >>> head, iterable = spy(iterable, 10)
        >>> head
        [1, 2, 3, 4, 5]
        >>> list(iterable)
        [1, 2, 3, 4, 5]

    """
    it = iter(iterable)
    head = take(n, it)

    return head.copy(), chain(head, it)


def take(n: int, iterable: Iterable[T]):
    """Return first *n* items of the iterable as a list.

        >>> take(3, range(10))
        [0, 1, 2]

    If there are fewer than *n* items in the iterable, all of them are
    returned.

        >>> take(10, range(3))
        [0, 1, 2]

    """
    return list(islice(iterable, n))


def batched(iterable: List[T], n: int, *, strict: bool = False):
    """Batch data into tuples of length *n*. If the number of items in
    *iterable* is not divisible by *n*:
    * The last batch will be shorter if *strict* is ``False``.
    * :exc:`ValueError` will be raised if *strict* is ``True``.

    >>> list(batched('ABCDEFG', 3))
    [('A', 'B', 'C'), ('D', 'E', 'F'), ('G',)]

    On Python 3.13 and above, this is an alias for :func:`itertools.batched`.
    """
    if n < 1:
        raise ValueError("n must be at least one")
    it = iter(iterable)
    while batch := tuple(islice(it, n)):
        if strict and len(batch) != n:
            raise ValueError("batched(): incomplete batch")
        yield batch


def triplewise(iterable: Iterable[T]):
    """Return overlapping triplets from *iterable*.

    >>> list(triplewise('ABCDE'))
    [('A', 'B', 'C'), ('B', 'C', 'D'), ('C', 'D', 'E')]

    """
    # This deviates from the itertools documentation reciple - see
    # https://github.com/more-itertools/more-itertools/issues/889
    t1, t2, t3 = tee(iterable, 3)
    next(t3, None)
    next(t3, None)
    next(t2, None)
    return zip(t1, t2, t3)


def ncycles(iterable: Iterable[T], n: int):
    """Returns the sequence elements *n* times

    >>> list(ncycles(["a", "b"], 3))
    ['a', 'b', 'a', 'b', 'a', 'b']

    """
    return chain.from_iterable(repeat(tuple(iterable), n))


def collapse(
    iterable: Iterable[T], *, levels: Optional[int] = None, base_type: Optional[type] = None
):
    """Flatten an iterable with multiple levels of nesting (e.g., a list of
    lists of tuples) into non-iterable types.

        >>> iterable = [(1, 2), ([3, 4], [[5], [6]])]
        >>> list(collapse(iterable))
        [1, 2, 3, 4, 5, 6]

    Binary and text strings are not considered iterable and
    will not be collapsed.

    To avoid collapsing other types, specify *base_type*:

        >>> iterable = ['ab', ('cd', 'ef'), ['gh', 'ij']]
        >>> list(collapse(iterable, base_type=tuple))
        ['ab', ('cd', 'ef'), 'gh', 'ij']

    Specify *levels* to stop flattening after a certain level:

    >>> iterable = [('a', ['b']), ('c', ['d'])]
    >>> list(collapse(iterable))  # Fully flattened
    ['a', 'b', 'c', 'd']
    >>> list(collapse(iterable, levels=1))  # Only one level flattened
    ['a', ['b'], 'c', ['d']]

    """
    stack = deque()
    # Add our first node group, treat the iterable as a single node
    stack.appendleft((0, repeat(iterable, 1)))

    while stack:
        node_group = stack.popleft()
        level, nodes = node_group

        # Check if beyond max level
        if levels is not None and level > levels:
            yield from nodes
            continue

        for node in nodes:
            # Check if done iterating
            if isinstance(node, (str, bytes)) or (
                (base_type is not None) and isinstance(node, base_type)
            ):
                yield node
            # Otherwise try to create child nodes
            else:
                try:
                    tree = iter(node)
                except TypeError:
                    yield node
                else:
                    # Save our current location
                    stack.appendleft(node_group)
                    # Append the new child node
                    stack.appendleft((level + 1, tree))
                    # Break to process child node
                    break


def prepend(value: T, iterator: Iterable[T]):
    """Yield *value*, followed by the elements in *iterator*.

        >>> value = '0'
        >>> iterator = ['1', '2', '3']
        >>> list(prepend(value, iterator))
        ['0', '1', '2', '3']

    To prepend multiple values, see :func:`itertools.chain`
    or :func:`value_chain`.

    """
    return chain([value], iterator)


def nth(iterable: Iterable[T], n: int, default=None):
    """Returns the nth item or a default value.

    >>> l = range(10)
    >>> nth(l, 3)
    3
    >>> nth(l, 20, "zebra")
    'zebra'

    """
    return next(islice(iterable, n, None), default)


def tail(n: int, iterable: Iterable[T]):
    """Return an iterator over the last *n* items of *iterable*.

    >>> t = tail(3, 'ABCDEFG')
    >>> list(t)
    ['E', 'F', 'G']

    """
    # If the given iterable has a length, then we can use islice to get its
    # final elements. Note that if the iterable is not actually Iterable,
    # either islice or deque will throw a TypeError. This is why we don't
    # check if it is Iterable.
    if isinstance(iterable, Sized):
        yield from islice(iterable, max(0, len(iterable) - n), None)
    else:
        yield from iter(deque(iterable, maxlen=n))


def unique(iterable, key=None, reverse=False):
    """Yields unique elements in sorted order.

    >>> list(unique([[1, 2], [3, 4], [1, 2]]))
    [[1, 2], [3, 4]]

    *key* and *reverse* are passed to :func:`sorted`.

    >>> list(unique('ABBcCAD', str.casefold))
    ['A', 'B', 'c', 'D']
    >>> list(unique('ABBcCAD', str.casefold, reverse=True))
    ['D', 'c', 'B', 'A']

    The elements in *iterable* need not be hashable, but they must be
    comparable for sorting to work.
    """
    return _unique_justseen(sorted(iterable, key=key, reverse=reverse), key=key)


def _unique_justseen(iterable, key=None):
    """Yields elements in order, ignoring serial duplicates

    >>> list(unique_justseen('AAAABBBCCDAABBB'))
    ['A', 'B', 'C', 'D', 'A', 'B']
    >>> list(unique_justseen('ABBCcAD', str.lower))
    ['A', 'B', 'C', 'A', 'D']

    """
    if key is None:
        return map(operator.itemgetter(0), groupby(iterable))

    return map(next, map(operator.itemgetter(1), groupby(iterable, key)))


def iter_index(iterable, value, start=0, stop=None):
    """Yield the index of each place in *iterable* that *value* occurs,
    beginning with index *start* and ending before index *stop*.


    >>> list(iter_index('AABCADEAF', 'A'))
    [0, 1, 4, 7]
    >>> list(iter_index('AABCADEAF', 'A', 1))  # start index is inclusive
    [1, 4, 7]
    >>> list(iter_index('AABCADEAF', 'A', 1, 7))  # stop index is not inclusive
    [1, 4]

    The behavior for non-scalar *values* matches the built-in Python types.

    >>> list(iter_index('ABCDABCD', 'AB'))
    [0, 4]
    >>> list(iter_index([0, 1, 2, 3, 0, 1, 2, 3], [0, 1]))
    []
    >>> list(iter_index([[0, 1], [2, 3], [0, 1], [2, 3]], [0, 1]))
    [0, 2]

    See :func:`locate` for a more general means of finding the indexes
    associated with particular values.

    """
    seq_index = getattr(iterable, "index", None)
    if seq_index is None:
        # Slow path for general iterables
        it = islice(iterable, start, stop)
        for i, element in enumerate(it, start):
            if element is value or element == value:
                yield i
    else:
        # Fast path for sequences
        stop = len(iterable) if stop is None else stop
        i = start - 1
        try:
            while True:
                yield (i := seq_index(value, i + 1, stop))
        except ValueError:
            pass


def consume(iterator: Iterable[T], n: Optional[int] = None):
    """Advance *iterable* by *n* steps. If *n* is ``None``, consume it
    entirely.

    Efficiently exhausts an iterator without returning values. Defaults to
    consuming the whole iterator, but an optional second argument may be
    provided to limit consumption.

        >>> i = (x for x in range(10))
        >>> next(i)
        0
        >>> consume(i, 3)
        >>> next(i)
        4
        >>> consume(i)
        >>> next(i)
        Traceback (most recent call last):
          File "<stdin>", line 1, in <module>
        StopIteration

    If the iterator has fewer items remaining than the provided limit, the
    whole iterator will be consumed.

        >>> i = (x for x in range(3))
        >>> consume(i, 5)
        >>> next(i)
        Traceback (most recent call last):
          File "<stdin>", line 1, in <module>
        StopIteration

    """
    # Use functions that consume iterators at C speed.
    if n is None:
        # feed the entire iterator into a zero-length deque
        deque(iterator, maxlen=0)
    else:
        # advance to the empty slice starting at position n
        next(islice(iterator, n, n), None)


if __name__ == "__main__":

    def chunked_even_demo():
        iterable = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
        n = 3
        print(list(chunked_even(iterable, n)))

    def distribute_demo():
        n = 3
        iterable = [1, 2, 3, 4, 5, 6, 7]
        children = distribute(n, iterable)
        print([list(c) for c in children])

    def triplewise_demo():
        print(list(triplewise("ABCDE")))

    def collapse_demo():
        iterable = [(1, 2), ([3, 4], [[5], [6]])]
        print(list(collapse(iterable)))

    chunked_even_demo()
    distribute_demo()
    triplewise_demo()
    collapse_demo()
