// def set_partitions(
//   iterable: Iterable[T],
//   k: Optional[int] = None,
//   min_size: Optional[int] = None,
//   max_size: Optional[int] = None,
// ):
//   """
//   Yield the set partitions of *iterable* into *k* parts. Set partitions are
//   not order-preserving.

//   >>> iterable = 'abc'
//   >>> for part in set_partitions(iterable, 2):
//   ...     print([''.join(p) for p in part])
//   ['a', 'bc']
//   ['ab', 'c']
//   ['b', 'ac']

//   If *k* is not given, every set partition is generated.

//   >>> iterable = 'abc'
//   >>> for part in set_partitions(iterable):
//   ...     print([''.join(p) for p in part])
//   ['abc']
//   ['a', 'bc']
//   ['ab', 'c']
//   ['b', 'ac']
//   ['a', 'b', 'c']

//   if *min_size* and/or *max_size* are given, the minimum and/or maximum size
//   per block in partition is set.

//   >>> iterable = 'abc'
//   >>> for part in set_partitions(iterable, min_size=2):
//   ...     print([''.join(p) for p in part])
//   ['abc']
//   >>> for part in set_partitions(iterable, max_size=2):
//   ...     print([''.join(p) for p in part])
//   ['a', 'bc']
//   ['ab', 'c']
//   ['b', 'ac']
//   ['a', 'b', 'c']

//   """
//   L = list(iterable)
//   n = len(L)
//   if k is not None:
//       if k < 1:
//           raise ValueError("Can't partition in a negative or zero number of groups")
//       elif k > n:
//           return

//   min_size = min_size if min_size is not None else 0
//   max_size = max_size if max_size is not None else n
//   if min_size > max_size:
//       return

//   def set_partitions_helper(L, k):
//       n = len(L)
//       if k == 1:
//           yield [L]
//       elif n == k:
//           yield [[s] for s in L]
//       else:
//           e, *M = L
//           for p in set_partitions_helper(M, k - 1):
//               yield [[e], *p]
//           for p in set_partitions_helper(M, k):
//               for i in range(len(p)):
//                   yield p[:i] + [[e] + p[i]] + p[i + 1 :]

//   if k is None:
//       for k in range(1, n + 1):
//           yield from filter(
//               lambda z: all(min_size <= len(bk) <= max_size for bk in z),
//               set_partitions_helper(L, k),
//           )
//   else:
//       yield from filter(
//           lambda z: all(min_size <= len(bk) <= max_size for bk in z),
//           set_partitions_helper(L, k),
//       )

/**
 * 将集合分割成若干个部分.
 */
function setPartitions(
  n: number,
  f: (parts: readonly number[][]) => boolean | void,
  options?: {
    k?: number
    minSize?: number
    maxSize?: number
  }
): void {
  dfs(0, [])
  const dfs = (index: number, parts: number[][]): boolean => {}
}
