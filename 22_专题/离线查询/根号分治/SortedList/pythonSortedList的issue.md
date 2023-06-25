- [If SortedList is given a sequence, figure out load automatically based on length # 4](https://github.com/grantjenks/python-sortedcontainers/issues/4)
  有人建议负载因子应该随长度变化动态调整
  在作者的本地测试中，以列表长度的立方根作为加载因子似乎效果最好

  ```python
  import timeit

  def benchmark(func):
      def wrapper(*args, **kwargs):
          times = [func(*args, **kwargs) for repeat in xrange(5)]
          return sum(times) / len(times)
      return wrapper

  @benchmark
  def insert(size, load, repeat=100):
      command = 'sl.add(random.randrange({0}))'.format(size)
      setup = ';'.join(
          ['import random',
           'from sortedcontainers import SortedList',
           'sl = SortedList(xrange({0}), load={1})'.format(size, load)])
      return timeit.timeit(command, setup=setup, number=repeat)

  def test():
      for size in [int(1e5), int(1e6), int(5e6), int(1e7), int(2e7)]:
          times = [(insert(size, load), load)
                   for load in [50, 100, 150, 200, 1000]]
          times.sort()
          print 'Size:', size, 'Load:', times[0][1]

  test()
  ```

  答案可能是“视情况而定”，因此基准测试是最好的方法。

- [Implement key argument to all `__init__` methods. #5](https://github.com/grantjenks/python-sortedcontainers/issues/5)
  有人建议增加 key 参数
  作者认为作为一组新的类型会更好。对于那些不需要该功能的用户，它会降低性能

- [Have .add(value) return the index for value #6](https://github.com/grantjenks/python-sortedcontainers/issues/6)
  有人建议 add 返回插入的下标，这样可以避免再次查找

- [Make SortedDict and SortedSet inherit from dict/set for speed improvements #17](https://github.com/grantjenks/python-sortedcontainers/issues/17)
  作者发现继承 dict 和 set 而不是组合，会有性能提升

- [Support thread-safety #105](https://github.com/grantjenks/python-sortedcontainers/issues/105)
  有人希望支持线程安全
  但作者表示不方便。reduce 就不是线程安全的(需要一个防止删除的锁)

- [The performance of add method could be improved #162](https://github.com/grantjenks/python-sortedcontainers/issues/162)
  有人表明不用`bisect.insort`，而是切片插入的方式会更快
  ```python
  p = bisect_right(_lists[pos], value)
  _lists[pos][p:p] = [value]
  ```
  经过更多测试，**这个技巧对小列表没有帮助**。由于 SortedList 中的 DEFAULT_LOAD_FACTOR 是 1000，因此无法提高 add 方法的性能。
- [Can not distinguish two different items with the same value #161](<[https://](https://github.com/grantjenks/python-sortedcontainers/issues/161)>)
  SortedList 无法区分两个具有相同值的不同对象
  作者在 https://grantjenks.com/docs/sortedcontainers/introduction.html#caveats 中提到了这个问题
  排序容器数据类型有三个要求：

  1. The comparison value or key must have a total ordering.
     比较值或键必须具有总排序。

  2. The comparison value or key must not change while the value is stored in the sorted container.
     当比较值或键存储在已排序的容器中时，该值不得更改。

  3. If the key-function parameter is used, then equal values must have equal keys.
     **如果使用键函数参数，则相等的值必须具有相等的键。**

  如果违反了这三个要求中的任何一个，则分类容器的保修无效，并且将无法正常运行。

- ["lower_bound()" function same as that of C++ STL's "map" library #172](https://github.com/grantjenks/python-sortedcontainers/issues/172)
  `SortedDict.irange`类似于 C++ STL 的 `map` 的 `lower_bound` 函数
- [Feature Proposal: Introduce Higher Level APIs like ceil / floor #207](https://github.com/grantjenks/python-sortedcontainers/issues/207)
  有人建议：引入更高级别的 API，如 ceil / floor
  这之前，在[How to do floor and ceiling lookups? #87](https://github.com/grantjenks/python-sortedcontainers/issues/87)中作者提到了性能更好的 floor 和 ceiling 方法的实现
  作者认为`irange`可以替代这两个功能
  但是，`irange`使用生成器，效率并不高。特别是对于仅选择一个值的场合，它将引入额外的计算开销
