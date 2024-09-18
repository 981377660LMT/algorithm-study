https://pypi.org/project/more-itertools/
https://docs.python.org/3/library/itertools.html#itertools-recipes
https://more-itertools.readthedocs.io/en/stable/index.html
https://github.com/EndlessCheng/codeforces-go/blob/8d4ab1e456a90853b4d23fe9b79aa881ade82508/copypasta/search.go#L1

TODO: 标记为粗体的

子集、组合、分割、排列、剪枝

---

## itertools

- filterfalse
  过滤掉所有返回 False 的元素
- compress
  只保留 selectors 中为 True 的元素
- cycle
  无限循环迭代器
- islice
  迭代器切片
- tee
  生成 n 个独立的迭代器
  https://gairuo.com/p/python-itertools-tee
- startmap
  将迭代器的元素作为函数的参数

---

## more-itertools

1. Grouping

- chunked
- ichunked
- **chunked_even**
  将迭代器尽可能均等分成`每组 n 个元素`，有些组可能比 n 少一个元素
- sliced
- constrained_batches
- **distribute**
  将元素尽可能均等`分配到 n 个迭代器`中，适用于多线程执行任务
- divide

- **split_at**
- **split_before**
- **split_after**
  根据谓词分割迭代器
- split_into
- split_when
- bucket
- **unzip**
- **batched**
  itertools.batched，将迭代器分成大小为 n 的块
- grouper
- partition
- transpose

2. Lookahead and lookback

- **spy**
  查看迭代器顶部若干个元素
- **peekable**
  创建可查看下一个元素的迭代器
- **seekable**
  即使在消耗了一些元素之后，也可以通过迭代器来回移动

3. Windowing

- windowed
- substrings
- subslices
- substrings_indexes
  所有非空子串
- stagger
- windowed_complete
- pairwise
- **triplewise**
  ```py
  t1, t2, t3 = tee(iterable, 3)
  next(t3, None)
  next(t3, None)
  next(t2, None)
  return zip(t1, t2, t3)
  ```
- **sliding_window**
  对每个窗口应用一个函数

4. Augmenting

- count_cycle
- **intersperse**
  将填充元素 e 散布在 iterable 中的项目之间，在每个填充元素之间留下 n 个项目。
- padded
- **repeat_each**
  重复每个元素 n 次
- mark_ends
- repeat_last
  迭代器耗尽后，重复最后一个元素若干次
- adjacent
- groupby_transform
- pad_none
- **ncycles**
  重复 n 次所有元素

5. Combining

- **collapse**
  展平多层嵌套的迭代器

- sort_together
- **interleave**
  对每个迭代器交错取出一个元素，直到最短的迭代器耗尽
- **interleave_longest**
  对每个迭代器交错取出一个元素，直到最长的迭代器耗尽
- **interleave_evenly**
  交错多个可迭代对象，使其元素均匀分布在整个输出序列中。
- zip_offset
- zip_equal
- zip_broadcast
- flatten
  展开一层嵌套的迭代器
- roundrobin
- **prepend**
  ```py
  chain([x], iterable)
  ```
- value_chain
- partial_product

6. Summarizing

- **ilen**
  统计迭代器元素个数
  ```py
  sum(compress(repeat(1), zip(iterable)))
  ```
- **unique_to_each**
  返回每个输入可迭代对象中不属于其他输入可迭代对象的元素。
- sample
- consecutive_groups
- **run_length**
  实现非常优雅
- **map_reduce**
  MapReduce 实现允许我们指定 3 个函数：键函数（用于分类）、值函数（用于转换）和最后的 reduce 函数（用于规约）。

- join_mappings
- exactly_n
- is_sorted
- all_equal
- all_unique
- minmax
- first_true
- quantify
  统计 True 的个数
- iequals

7. Selecting

- islice_extended
- first
- last
- one
- only
- strictly_n
- **strip**
- **lstrip**
- **rstrip**
  从开头和结尾中删除 pred 返回 True 的元素
- filter_except
  通过检查的元素来过滤输入可迭代的项目
- map_except
- filter_map
- **iter_suppress**
- nth_or_last
- unique_in_window
- before_and_after
- **nth**

  ```py
  next(islice(iterable, n, None), default)
  ```

- **take**
  ```py
  list(islice(iterable, n))
  ```
- tail
  ```py
  iter(collections.deque(iterable, maxlen=n))
  ```
- unique_everseen
- unique_justseen
- **unique**
  迭代器排序去重
- duplicates_everseen
- duplicates_justseen
- classify_unique
- longest_common_prefix
- takewhile_inclusive

8. Math

- dft
- idft
- convolve
- dotproduct
- **factor**
- matmul
  矩阵乘法
- **polynomial_from_roots**
  根据多项式的根生成多项式
- **polynomial_derivative**
  多项式求导
- **polynomial_eval**
  多项式求值
- **sieve**
- sum_of_squares
- **totient**
  不超过 n 的正整数中与 n 互质的数的个数(欧拉函数)

9. Combinatorics

- **distinct_permutations**
  无重复元素的排列，一共 `n!/c1!c2!...cn!` 个
- **distinct_combinations**
  无重复元素的组合
- **circular_shifts**
- **partitions**
  所有划分
- **set_partitions**
  将可迭代的集合划分为 k 个部分
- **product_index**
  计算笛卡尔积的索引
- **combination_index**
- **permutation_index**
- **combination_with_replacement_index**
- gray_product
- outer_product
- **powerset**
  2^n 个子集
- **powerset_of_sets**
  2^n 个 set()
- **random_product**
- **random_permutation**
- **random_combination**
- **random_combination_with_replacement**
- **nth_product**
- **nth_permutation**
- **nth_combination**
- **nth_combination_with_replacement**

10. Wrapping

- always_iterable
- always_reversible
- countable
- consumer
- with_iter
- iter_except

11. Others

- **locate**
- **iter_index**
- **rlocate**
  找到元素的位置(返回数组)
- replace
- numeric_range
- **side_effect**
  每次迭代时调用一个函数
- iterate
- **difference**
  accumulate 的差分版本
- make_decorator
- SequenceView
- time_limited
- map_if
- **consume**
  消耗迭代器中的 n 个元素
- tabulate
  制表工具
  ```py
   map(f, count(start))
  ```
- repeatfunc
  重复调用一个函数
- reshape
- doublestarmap
