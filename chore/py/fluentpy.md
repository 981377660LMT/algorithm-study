```py
>>> t = (1, 2, [30, 40])
>>> t[2] += [50, 60]
Traceback (most recent call last):
File "<stdin>", line 1, in <module>
TypeError: 'tuple' object does not support item assignment
>>> t
(1, 2, [30, 40, 50, 60])
```

- 不要把可变对象放在元组里面。
- 增量赋值不是一个原子操作。我们刚才也看到了，它虽然抛出了异
  常，但还是完成了操作。
- 查看 Python 的字节码并不难，而且它对我们了解代码背后的运行机
  制很有帮助。

- 当列表不是首选时
  如果我们需要一个只包含数字的列表，那么 array.array 比 list 更节省内存
  **array 相当于 js 中的 TypedArray，但是是可变长的**
  类型数组

  ```py
  arr = array("I", [1, 2, 3, 4, 5, 6, 7, 8, 9, 10])
  #  bad typecode (must be b, B, u, h, H, i, I, l, L, q, Q, f or d)
  # 'b' -> Int8    (byte 表示 8)
  # 'B' -> UInt8
  # 'h' -> Int16   (hex 表示 16)
  # 'H' -> UInt16
  # 'i' -> Int32   (int 表示 32)
  # 'I' -> UInt32
  # 'q' -> Int64   (quad 表示 64)
  # 'Q' -> UInt64
  # 'f' -> Float32 (float)
  # 'd' -> Float64 (double)
  ```

  **bytearray 可以理解为 golang 中的 []byte**

  [bytearray 和 array.array('B') 的区别](https://stackoverflow.com/questions/11882988/python-bytearray-vs-array)

  - 需要注意的是 bytearray 比 array('B') 更快，因为后者在每个元素上要付出额外的存储成本，用来存储类型码,同时可以认为 bytearray 是一个可变的字符串 `mutable str (bytes in Python3)`,因为他有 str 的各种方法
  - 有时候,array 不一定比 list 快
  - py 可以用 bytearray.decode()来解码成 str,但是 golang 中的[]byte 没有这个方法,需要用 string([]byte)来转换成 str,如果要指定解码方式,需要用 utf8/utf16 库中的 DecodeRune 方法

  **memory_view 可以理解为 js 里的**

- key 参数很妙
  key 参数也能让你对一个混有数字字符和数值的列表进行排序

  ```py
  >>> l = [28, 14, '28', 5, '9', '1', 0, 6, '23', 19]
  >>> sorted(l, key=int)
  [0, '1', 5, 6, '9', 14, 19, '23', 28, '28']
  >>> sorted(l, key=str)
  [0, '1', 14, 19, '23', 28, '28', 5, 6, '9']
  ```

- “Python 里所有的不可变类型都是可散列的”。这个说法其实是不
  准确的，比如虽然元组本身是不可变序列，它里面的元素可能是其
  他可变类型的引用。

- defaultdict 里的 default_factory 只会在`__getitem__` 里被调用，在其他的方法里完全不会发挥作用。比如，dd 是个 defaultdict，k 是个找不到的键， dd[k] 这个表达式会调用 default_factory 创造某个默认值，而 dd.get(k) 则会返回 None。

- 不要忘了，如果要创建一个空集合，你必须用不带任何参数的构造方
  法 set()。如果只是写成 {} 的形式，跟以前一样，你创建的其实
  是个空字典。
