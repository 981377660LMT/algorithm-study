# distinct_permutations
# 无重复元素的排列，一共 `n!/c1!c2!...cn!` 个
#
# distinct_combinations
# 无重复元素的组合
#
# circular_shifts
#
# partitions
# 所有划分
#
# set_partitions
# 将可迭代的集合划分为 k 个部分
#
# product_index
# combination_index
# permutation_index
# combination_with_replacement_index
# 计算笛卡尔积的索引
#
# powerset
# 2^n 个子集
#
# powerset_of_sets
# 2^n 个 set()
#
# random_product
# random_permutation
# random_combination
# random_combination_with_replacement
#
# nth_product
# nth_permutation
# nth_combination
# nth_combination_with_replacement

import math
from random import choice, randrange, sample
from sys import hexversion
from collections import defaultdict, deque
from functools import partial, reduce
from itertools import chain, combinations, cycle, repeat, starmap, zip_longest
from operator import itemgetter, mul
from typing import Any, Iterable, Optional, Sequence, TypeVar

T = TypeVar("T")

_marker: Any = object()


def distinct_permutations(iterable: Iterable[T], r: Optional[int] = None):
    """Yield successive distinct permutations of the elements in *iterable*.

        >>> sorted(distinct_permutations([1, 0, 1]))
        [(0, 1, 1), (1, 0, 1), (1, 1, 0)]

    Equivalent to yielding from ``set(permutations(iterable))``, except
    duplicates are not generated and thrown away. For larger input sequences
    this is much more efficient.

    Duplicate permutations arise when there are duplicated elements in the
    input iterable. The number of items returned is
    `n! / (x_1! * x_2! * ... * x_n!)`, where `n` is the total number of
    items input, and each `x_i` is the count of a distinct item in the input
    sequence.

    If *r* is given, only the *r*-length permutations are yielded.

        >>> sorted(distinct_permutations([1, 0, 1], r=2))
        [(0, 1), (1, 0), (1, 1)]
        >>> sorted(distinct_permutations(range(3), r=2))
        [(0, 1), (0, 2), (1, 0), (1, 2), (2, 0), (2, 1)]

    *iterable* need not be sortable, but note that using equal (``x == y``)
    but non-identical (``id(x) != id(y)``) elements may produce surprising
    behavior. For example, ``1`` and ``True`` are equal but non-identical:

        >>> list(distinct_permutations([1, True, '3']))  # doctest: +SKIP
        [
            (1, True, '3'),
            (1, '3', True),
            ('3', 1, True)
        ]
        >>> list(distinct_permutations([1, 2, '3']))  # doctest: +SKIP
        [
            (1, 2, '3'),
            (1, '3', 2),
            (2, 1, '3'),
            (2, '3', 1),
            ('3', 1, 2),
            ('3', 2, 1)
        ]
    """

    # Algorithm: https://w.wiki/Qai
    def _full(A):
        while True:
            # Yield the permutation we have
            yield tuple(A)

            # Find the largest index i such that A[i] < A[i + 1]
            for i in range(size - 2, -1, -1):
                if A[i] < A[i + 1]:
                    break
            #  If no such index exists, this permutation is the last one
            else:
                return

            # Find the largest index j greater than j such that A[i] < A[j]
            for j in range(size - 1, i, -1):
                if A[i] < A[j]:
                    break

            # Swap the value of A[i] with that of A[j], then reverse the
            # sequence from A[i + 1] to form the new permutation
            A[i], A[j] = A[j], A[i]  # type: ignore
            A[i + 1 :] = A[: i - size : -1]  # A[i + 1:][::-1]

    # Algorithm: modified from the above
    def _partial(A, r):
        # Split A into the first r items and the last r items
        head, tail = A[:r], A[r:]
        right_head_indexes = range(r - 1, -1, -1)
        left_tail_indexes = range(len(tail))

        while True:
            # Yield the permutation we have
            yield tuple(head)

            # Starting from the right, find the first index of the head with
            # value smaller than the maximum value of the tail - call it i.
            pivot = tail[-1]
            for i in right_head_indexes:
                if head[i] < pivot:
                    break
                pivot = head[i]
            else:
                return

            # Starting from the left, find the first value of the tail
            # with a value greater than head[i] and swap.
            for j in left_tail_indexes:
                if tail[j] > head[i]:
                    head[i], tail[j] = tail[j], head[i]
                    break
            # If we didn't find one, start from the right and find the first
            # index of the head with a value greater than head[i] and swap.
            else:
                for j in right_head_indexes:
                    if head[j] > head[i]:
                        head[i], head[j] = head[j], head[i]
                        break

            # Reverse head[i + 1:] and swap it with tail[:r - (i + 1)]
            tail += head[: i - r : -1]  # head[i + 1:][::-1]
            i += 1
            head[i:], tail[:] = tail[: r - i], tail[r - i :]

    items = list(iterable)

    try:
        items.sort()  # type: ignore
        sortable = True
    except TypeError:
        sortable = False

        indices_dict = defaultdict(list)

        for item in items:
            indices_dict[items.index(item)].append(item)

        indices = [items.index(item) for item in items]
        indices.sort()

        equivalent_items = {k: cycle(v) for k, v in indices_dict.items()}

        def permuted_items(permuted_indices):
            return tuple(next(equivalent_items[index]) for index in permuted_indices)

    size = len(items)
    if r is None:
        r = size

    # functools.partial(_partial, ... )
    algorithm = _full if (r == size) else partial(_partial, r=r)

    if 0 < r <= size:
        if sortable:
            return algorithm(items)
        else:
            return (permuted_items(permuted_indices) for permuted_indices in algorithm(indices))  # type: ignore

    return iter(() if r else ((),))


def distinct_combinations(iterable: Iterable[T], r: int):
    """Yield the distinct combinations of *r* items taken from *iterable*.

        >>> list(distinct_combinations([0, 0, 1], 2))
        [(0, 0), (0, 1)]

    Equivalent to ``set(combinations(iterable))``, except duplicates are not
    generated and thrown away. For larger input sequences this is much more
    efficient.

    """
    if r < 0:
        raise ValueError("r must be non-negative")
    elif r == 0:
        yield ()
        return
    pool = tuple(iterable)
    generators = [_unique_everseen(enumerate(pool), key=itemgetter(1))]
    current_combo = [None] * r
    level = 0
    while generators:
        try:
            cur_idx, p = next(generators[-1])
        except StopIteration:
            generators.pop()
            level -= 1
            continue
        current_combo[level] = p
        if level + 1 == r:
            yield tuple(current_combo)
        else:
            generators.append(
                _unique_everseen(
                    enumerate(pool[cur_idx + 1 :], cur_idx + 1),
                    key=itemgetter(1),
                )
            )
            level += 1


def _unique_everseen(iterable, key=None):
    """
    Yield unique elements, preserving order.

        >>> list(unique_everseen('AAAABBBCCDAABBB'))
        ['A', 'B', 'C', 'D']
        >>> list(unique_everseen('ABBCcAD', str.lower))
        ['A', 'B', 'C', 'D']

    Sequences with a mix of hashable and unhashable items can be used.
    The function will be slower (i.e., `O(n^2)`) for unhashable items.

    Remember that ``list`` objects are unhashable - you can use the *key*
    parameter to transform the list to a tuple (which is hashable) to
    avoid a slowdown.

        >>> iterable = ([1, 2], [2, 3], [1, 2])
        >>> list(unique_everseen(iterable))  # Slow
        [[1, 2], [2, 3]]
        >>> list(unique_everseen(iterable, key=tuple))  # Faster
        [[1, 2], [2, 3]]

    Similarly, you may want to convert unhashable ``set`` objects with
    ``key=frozenset``. For ``dict`` objects,
    ``key=lambda x: frozenset(x.items())`` can be used.

    """
    seenset = set()
    seenset_add = seenset.add
    seenlist = []
    seenlist_add = seenlist.append
    use_key = key is not None

    for element in iterable:
        k = key(element) if use_key else element
        try:
            if k not in seenset:
                seenset_add(k)
                yield element
        except TypeError:
            if k not in seenlist:
                seenlist_add(k)
                yield element


def circular_shifts(iterable: Iterable[T], steps: int = 1):
    """Yield the circular shifts of *iterable*.

    >>> list(circular_shifts(range(4)))
    [(0, 1, 2, 3), (1, 2, 3, 0), (2, 3, 0, 1), (3, 0, 1, 2)]

    Set *steps* to the number of places to rotate to the left
    (or to the right if negative).  Defaults to 1.

    >>> list(circular_shifts(range(4), 2))
    [(0, 1, 2, 3), (2, 3, 0, 1)]

    >>> list(circular_shifts(range(4), -1))
    [(0, 1, 2, 3), (3, 0, 1, 2), (2, 3, 0, 1), (1, 2, 3, 0)]

    """
    buffer = deque(iterable)
    if steps == 0:
        raise ValueError("Steps should be a non-zero integer")

    buffer.rotate(steps)
    steps = -steps
    n = len(buffer)
    n //= math.gcd(n, steps)

    for _ in repeat(None, n):
        buffer.rotate(steps)
        yield tuple(buffer)


def partitions(iterable: Iterable[T]):
    """Yield all possible order-preserving partitions of *iterable*.

    >>> iterable = 'abc'
    >>> for part in partitions(iterable):
    ...     print([''.join(p) for p in part])
    ['abc']
    ['a', 'bc']
    ['ab', 'c']
    ['a', 'b', 'c']

    This is unrelated to :func:`partition`.

    """
    sequence = list(iterable)
    n = len(sequence)
    for i in powerset(range(1, n)):
        yield [sequence[i:j] for i, j in zip((0,) + i, i + (n,))]


def set_partitions(
    iterable: Iterable[T],
    k: Optional[int] = None,
    min_size: Optional[int] = None,
    max_size: Optional[int] = None,
):
    """
    Yield the set partitions of *iterable* into *k* parts. Set partitions are
    not order-preserving.

    >>> iterable = 'abc'
    >>> for part in set_partitions(iterable, 2):
    ...     print([''.join(p) for p in part])
    ['a', 'bc']
    ['ab', 'c']
    ['b', 'ac']


    If *k* is not given, every set partition is generated.

    >>> iterable = 'abc'
    >>> for part in set_partitions(iterable):
    ...     print([''.join(p) for p in part])
    ['abc']
    ['a', 'bc']
    ['ab', 'c']
    ['b', 'ac']
    ['a', 'b', 'c']

    if *min_size* and/or *max_size* are given, the minimum and/or maximum size
    per block in partition is set.

    >>> iterable = 'abc'
    >>> for part in set_partitions(iterable, min_size=2):
    ...     print([''.join(p) for p in part])
    ['abc']
    >>> for part in set_partitions(iterable, max_size=2):
    ...     print([''.join(p) for p in part])
    ['a', 'bc']
    ['ab', 'c']
    ['b', 'ac']
    ['a', 'b', 'c']

    """
    L = list(iterable)
    n = len(L)
    if k is not None:
        if k < 1:
            raise ValueError("Can't partition in a negative or zero number of groups")
        elif k > n:
            return

    min_size = min_size if min_size is not None else 0
    max_size = max_size if max_size is not None else n
    if min_size > max_size:
        return

    def set_partitions_helper(L, k):
        n = len(L)
        if k == 1:
            yield [L]
        elif n == k:
            yield [[s] for s in L]
        else:
            e, *M = L
            for p in set_partitions_helper(M, k - 1):
                yield [[e], *p]
            for p in set_partitions_helper(M, k):
                for i in range(len(p)):
                    yield p[:i] + [[e] + p[i]] + p[i + 1 :]

    if k is None:
        for k in range(1, n + 1):
            yield from filter(
                lambda z: all(min_size <= len(bk) <= max_size for bk in z),
                set_partitions_helper(L, k),
            )
    else:
        yield from filter(
            lambda z: all(min_size <= len(bk) <= max_size for bk in z),
            set_partitions_helper(L, k),
        )


def powerset(iterable: Iterable[T]):
    """Yields all possible subsets of the iterable.

        >>> list(powerset([1, 2, 3]))
        [(), (1,), (2,), (3,), (1, 2), (1, 3), (2, 3), (1, 2, 3)]

    :func:`powerset` will operate on iterables that aren't :class:`set`
    instances, so repeated elements in the input will produce repeated elements
    in the output.

        >>> seq = [1, 1, 0]
        >>> list(powerset(seq))
        [(), (1,), (1,), (0,), (1, 1), (1, 0), (1, 0), (1, 1, 0)]

    For a variant that efficiently yields actual :class:`set` instances, see
    :func:`powerset_of_sets`.
    """
    s = list(iterable)
    return chain.from_iterable(combinations(s, r) for r in range(len(s) + 1))


def powerset_of_sets(iterable: Iterable[T]):
    """Yields all possible subsets of the iterable.

        >>> list(powerset_of_sets([1, 2, 3]))  # doctest: +SKIP
        [set(), {1}, {2}, {3}, {1, 2}, {1, 3}, {2, 3}, {1, 2, 3}]
        >>> list(powerset_of_sets([1, 1, 0]))  # doctest: +SKIP
        [set(), {1}, {0}, {0, 1}]

    :func:`powerset_of_sets` takes care to minimize the number
    of hash operations performed.
    """
    sets = tuple(map(set, dict.fromkeys(map(frozenset, zip(iterable)))))
    for r in range(len(sets) + 1):
        yield from starmap(set().union, combinations(sets, r))


def product_index(element, *args):
    """Equivalent to ``list(product(*args)).index(element)``

    The products of *args* can be ordered lexicographically.
    :func:`product_index` computes the first index of *element* without
    computing the previous products.

        >>> product_index([8, 2], range(10), range(5))
        42

    ``ValueError`` will be raised if the given *element* isn't in the product
    of *args*.
    """
    index = 0

    for x, pool in zip_longest(element, args, fillvalue=_marker):
        if x is _marker or pool is _marker:
            raise ValueError("element is not a product of args")

        pool = tuple(pool)
        index = index * len(pool) + pool.index(x)

    return index


def combination_index(element, iterable):
    """Equivalent to ``list(combinations(iterable, r)).index(element)``

    The subsequences of *iterable* that are of length *r* can be ordered
    lexicographically. :func:`combination_index` computes the index of the
    first *element*, without computing the previous combinations.

        >>> combination_index('adf', 'abcdefg')
        10

    ``ValueError`` will be raised if the given *element* isn't one of the
    combinations of *iterable*.
    """
    element = enumerate(element)
    k, y = next(element, (None, None))
    if k is None:
        return 0

    indexes = []
    pool = enumerate(iterable)
    for n, x in pool:
        if x == y:
            indexes.append(n)
            tmp, y = next(element, (None, None))
            if tmp is None:
                break
            else:
                k = tmp
    else:
        raise ValueError("element is not a combination of iterable")

    n, _ = _last(pool, default=(n, None))

    # Python versions below 3.8 don't have math.comb
    index = 1
    for i, j in enumerate(reversed(indexes), start=1):
        j = n - j
        if i <= j:
            index += math.comb(j, i)

    return math.comb(n + 1, k + 1) - index


def _last(iterable, default=_marker):
    """Return the last item of *iterable*, or *default* if *iterable* is
    empty.

        >>> last([0, 1, 2, 3])
        3
        >>> last([], 'some default')
        'some default'

    If *default* is not provided and there are no items in the iterable,
    raise ``ValueError``.
    """
    try:
        if isinstance(iterable, Sequence):
            return iterable[-1]
        # Work around https://bugs.python.org/issue38525
        elif hasattr(iterable, "__reversed__") and (hexversion != 0x030800F0):
            return next(reversed(iterable))
        else:
            return deque(iterable, maxlen=1)[-1]
    except (IndexError, TypeError, StopIteration):
        if default is _marker:
            raise ValueError(
                "last() was called on an empty iterable, and no default was " "provided."
            )
        return default


def permutation_index(element, iterable):
    """Equivalent to ``list(permutations(iterable, r)).index(element)```

    The subsequences of *iterable* that are of length *r* where order is
    important can be ordered lexicographically. :func:`permutation_index`
    computes the index of the first *element* directly, without computing
    the previous permutations.

        >>> permutation_index([1, 3, 2], range(5))
        19

    ``ValueError`` will be raised if the given *element* isn't one of the
    permutations of *iterable*.
    """
    index = 0
    pool = list(iterable)
    for i, x in zip(range(len(pool), -1, -1), element):
        r = pool.index(x)
        index = index * i + r
        del pool[r]

    return index


def combination_with_replacement_index(element, iterable):
    """Equivalent to
    ``list(combinations_with_replacement(iterable, r)).index(element)``

    The subsequences with repetition of *iterable* that are of length *r* can
    be ordered lexicographically. :func:`combination_with_replacement_index`
    computes the index of the first *element*, without computing the previous
    combinations with replacement.

        >>> combination_with_replacement_index('adf', 'abcdefg')
        20

    ``ValueError`` will be raised if the given *element* isn't one of the
    combinations with replacement of *iterable*.
    """
    element = tuple(element)
    l = len(element)
    element = enumerate(element)

    k, y = next(element, (None, None))
    if k is None:
        return 0

    indexes = []
    pool = tuple(iterable)
    for n, x in enumerate(pool):
        while x == y:
            indexes.append(n)
            tmp, y = next(element, (None, None))
            if tmp is None:
                break
            else:
                k = tmp
        if y is None:
            break
    else:
        raise ValueError("element is not a combination with replacement of iterable")

    n = len(pool)
    occupations = [0] * n
    for p in indexes:
        occupations[p] += 1

    index = 0
    cumulative_sum = 0
    for k in range(1, n):
        cumulative_sum += occupations[k - 1]
        j = l + n - 1 - k - cumulative_sum
        i = n - k
        if i <= j:
            index += math.comb(j, i)

    return index


def random_product(*args, repeat=1):
    """Draw an item at random from each of the input iterables.

        >>> random_product('abc', range(4), 'XYZ')  # doctest:+SKIP
        ('c', 3, 'Z')

    If *repeat* is provided as a keyword argument, that many items will be
    drawn from each iterable.

        >>> random_product('abcd', range(4), repeat=2)  # doctest:+SKIP
        ('a', 2, 'd', 3)

    This equivalent to taking a random selection from
    ``itertools.product(*args, **kwarg)``.

    """
    pools = [tuple(pool) for pool in args] * repeat
    return tuple(choice(pool) for pool in pools)


def random_permutation(iterable, r=None):
    """Return a random *r* length permutation of the elements in *iterable*.

    If *r* is not specified or is ``None``, then *r* defaults to the length of
    *iterable*.

        >>> random_permutation(range(5))  # doctest:+SKIP
        (3, 4, 0, 1, 2)

    This equivalent to taking a random selection from
    ``itertools.permutations(iterable, r)``.

    """
    pool = tuple(iterable)
    r = len(pool) if r is None else r
    return tuple(sample(pool, r))


def random_combination(iterable, r):
    """Return a random *r* length subsequence of the elements in *iterable*.

        >>> random_combination(range(5), 3)  # doctest:+SKIP
        (2, 3, 4)

    This equivalent to taking a random selection from
    ``itertools.combinations(iterable, r)``.

    """
    pool = tuple(iterable)
    n = len(pool)
    indices = sorted(sample(range(n), r))
    return tuple(pool[i] for i in indices)


def random_combination_with_replacement(iterable, r):
    """Return a random *r* length subsequence of elements in *iterable*,
    allowing individual elements to be repeated.

        >>> random_combination_with_replacement(range(3), 5) # doctest:+SKIP
        (0, 0, 1, 2, 2)

    This equivalent to taking a random selection from
    ``itertools.combinations_with_replacement(iterable, r)``.

    """
    pool = tuple(iterable)
    n = len(pool)
    indices = sorted(randrange(n) for _ in range(r))
    return tuple(pool[i] for i in indices)


def nth_product(index, *args):
    """Equivalent to ``list(product(*args))[index]``.

    The products of *args* can be ordered lexicographically.
    :func:`nth_product` computes the product at sort position *index* without
    computing the previous products.

        >>> nth_product(8, range(2), range(2), range(2), range(2))
        (1, 0, 0, 0)

    ``IndexError`` will be raised if the given *index* is invalid.
    """
    pools = list(map(tuple, reversed(args)))
    ns = list(map(len, pools))

    c = reduce(mul, ns)

    if index < 0:
        index += c

    if not 0 <= index < c:
        raise IndexError

    result = []
    for pool, n in zip(pools, ns):
        result.append(pool[index % n])
        index //= n

    return tuple(reversed(result))


def nth_permutation(iterable: Iterable[T], r: Optional[int], index: int):
    """Equivalent to ``list(permutations(iterable, r))[index]```

    The subsequences of *iterable* that are of length *r* where order is
    important can be ordered lexicographically. :func:`nth_permutation`
    computes the subsequence at sort position *index* directly, without
    computing the previous subsequences.

        >>> nth_permutation('ghijk', 2, 5)
        ('h', 'i')

    ``ValueError`` will be raised If *r* is negative or greater than the length
    of *iterable*.
    ``IndexError`` will be raised if the given *index* is invalid.
    """
    pool = list(iterable)
    n = len(pool)

    if r is None or r == n:
        r, c = n, math.factorial(n)
    elif not 0 <= r < n:
        raise ValueError
    else:
        c = math.perm(n, r)
    assert c > 0  # factortial(n)>0, and r<n so perm(n,r) is never zero

    if index < 0:
        index += c

    if not 0 <= index < c:
        raise IndexError

    result = [0] * r
    q = index * math.factorial(n) // c if r < n else index
    for d in range(1, n + 1):
        q, i = divmod(q, d)
        if 0 <= n - d < r:
            result[n - d] = i
        if q == 0:
            break

    return tuple(map(pool.pop, result))


def nth_combination(iterable: Iterable[T], r: int, index: int):
    """Equivalent to ``list(combinations(iterable, r))[index]``.

    The subsequences of *iterable* that are of length *r* can be ordered
    lexicographically. :func:`nth_combination` computes the subsequence at
    sort position *index* directly, without computing the previous
    subsequences.

        >>> nth_combination(range(5), 3, 5)
        (0, 3, 4)

    ``ValueError`` will be raised If *r* is negative or greater than the length
    of *iterable*.
    ``IndexError`` will be raised if the given *index* is invalid.
    """
    pool = tuple(iterable)
    n = len(pool)
    if (r < 0) or (r > n):
        raise ValueError

    c = 1
    k = min(r, n - r)
    for i in range(1, k + 1):
        c = c * (n - k + i) // i

    if index < 0:
        index += c

    if (index < 0) or (index >= c):
        raise IndexError

    result = []
    while r:
        c, n, r = c * r // n, n - 1, r - 1
        while index >= c:
            index -= c
            c, n = c * (n - r) // n, n - 1
        result.append(pool[-1 - n])

    return tuple(result)


def nth_combination_with_replacement(iterable: Iterable[T], r: int, index: int):
    """Equivalent to
    ``list(combinations_with_replacement(iterable, r))[index]``.


    The subsequences with repetition of *iterable* that are of length *r* can
    be ordered lexicographically. :func:`nth_combination_with_replacement`
    computes the subsequence at sort position *index* directly, without
    computing the previous subsequences with replacement.

        >>> nth_combination_with_replacement(range(5), 3, 5)
        (0, 1, 1)

    ``ValueError`` will be raised If *r* is negative or greater than the length
    of *iterable*.
    ``IndexError`` will be raised if the given *index* is invalid.
    """
    pool = tuple(iterable)
    n = len(pool)
    if (r < 0) or (r > n):
        raise ValueError

    c = math.comb(n + r - 1, r)

    if index < 0:
        index += c

    if (index < 0) or (index >= c):
        raise IndexError

    result = []
    i = 0
    while r:
        r -= 1
        while n >= 0:
            num_combs = math.comb(n + r - 1, r)
            if index < num_combs:
                break
            n -= 1
            i += 1
            index -= num_combs
        result.append(pool[i])

    return tuple(result)


if __name__ == "__main__":
    print(sorted(distinct_permutations([1, 0, 1])))
    print(sorted(distinct_permutations([1, 0, 1], r=2)))
    print(sorted(distinct_combinations([0, 0, 1], 2)))

    print(list(circular_shifts(range(4))))

    print(*partitions("abc"))
    print(*set_partitions("abc", 2))

    print(*powerset([1, 2, 3]))
    print(*powerset_of_sets([1, 1, 2, 3]))

    print(product_index([8, 2], range(10), range(5)))
    print(combination_index("adf", "abcdefg"))
    print(permutation_index([1, 3, 2], range(5)))
    print(combination_with_replacement_index("adf", "abcdefg"))

    print(random_product("abc", range(4), "XYZ"))
    print(random_permutation(range(5)))
    print(random_combination(range(5), 3))
    print(random_combination_with_replacement(range(3), 5))

    print(nth_product(8, range(2), range(2), range(2), range(2)))
    print(nth_permutation("ghijk", 2, 5))
    print(nth_combination(range(5), 3, 5))
    print(nth_combination_with_replacement(range(5), 3, 5))
