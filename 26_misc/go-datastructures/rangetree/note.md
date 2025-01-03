Useful to determine if n-dimensional points fall within an n-dimensional range. Not a typical range tree however, as we are actually using an n-dimensional sorted list of points as this proved to be simpler and faster than attempting a traditional range tree while saving space on any dimension greater than one. Inserts are typical BBST times at O(log n^d) where d is the number of dimensions.
有助于确定 n 维点是否落在 n 维范围内。然而，这并不是一个典型的范围树，因为我们实际上使用的是 n 维排序点列表，这被证明比尝试传统范围树更简单、更快速，同时在任何大于一的维度上节省空间。
插入操作的时间复杂度为 O(log n^d)，其中 d 是维度的数量。
