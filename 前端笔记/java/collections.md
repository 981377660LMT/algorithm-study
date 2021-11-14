Java 集合， 也叫作容器，主要是由两大接口派生而来：
一个是 Collecton 接口，主要用于存放单一元素；
另一个是 Map 接口，主要用于存放键值对。
对于 Collection 接口，下面又有三个主要的子接口：List、Set 和 Queue。
[collection](<[https://link](https://github.com/Snailclimb/JavaGuide/blob/main/docs/java/collection/java%E9%9B%86%E5%90%88%E6%A1%86%E6%9E%B6%E5%9F%BA%E7%A1%80%E7%9F%A5%E8%AF%86&%E9%9D%A2%E8%AF%95%E9%A2%98%E6%80%BB%E7%BB%93.md)>)

1. HashMap
   HashMap： JDK1.8 之前 HashMap 由`数组+链表`组成的，数组是 HashMap 的主体，链表则是主要为了解决哈希冲突而存在的（“拉链法”解决冲突）。
   JDK1.8 以后在解决哈希冲突时有了较大的变化，当`链表`长度大于阈值（默认为 8）（将链表转换成红黑树前会判断，如果当前`数组的长度小于 64`，那么会选择先进行数组扩容，而不是转换为红黑树）时，将链表转化为红黑树，以减少搜索时间
2. LinkedHashMap 类似于 LRU 结构 哈希表 Key 指向链表结点实现 O(1)
   LinkedHashMap 继承自 HashMap，所以它的底层仍然是基于拉链式散列结构即由数组和链表或红黑树组成。另外，LinkedHashMap 在上面结构的基础上，`增加了一条双向链表`，使得上面的结构可以保持键值对的插入顺序。同时通过对链表进行相应的操作，实现了访问顺序相关逻辑。
3. Hashtable： 数组+链表组成的，数组是 Hashtable 的主体，链表则是主要为了解决哈希冲突而存在的
4. TreeMap： 红黑树（自平衡的排序二叉树
   需要排序时选择 TreeMap,不需要排序时就选择 HashMap,需要保证线程安全就选用 ConcurrentHashMap。

5. RandomAccess 接口

```JAVA
public interface RandomAccess {
}

查看源码我们发现实际上 RandomAccess 接口中什么都没有定义。所以，在我看来 RandomAccess 接口不过是一个标识罢了。标识什么？ 标识实现这个接口的类具有随机访问功能。
```

// TODO

6.  ArrayList 的扩容机制
    ArrayList 的底层是数组队列，相当于动态数组。与 Java 中的数组相比，它的容量能动态增长。在添加大量元素前，应用程序可以使用 ensureCapacity 操作来增加 ArrayList 实例的容量。这可以减少递增式再分配的数量。
7.  HashMap 精华
    HashMap 默认的初始化大小为 16。之后每次扩充，容量变为原来的 2 倍。并且， HashMap 总是使用 2 的幂作为哈希表的大小。
    HashMap 通过 key 的 hashCode 经过扰动函数处理过后得到 hash 值，然后通过 `(n - 1) & hash` 判断当前元素存放的位置（这里的 n 指的是数组的长度），如果当前位置存在元素的话，就判断该元素与要存入的元素的 hash 值以及 key 是否相同，如果相同的话，直接覆盖，不相同就通过拉链法解决冲突。所谓扰动函数指的就是 HashMap 的 hash 方法。使用 hash 方法也就是扰动函数是为了防止一些实现比较差的 hashCode() 方法 换句话说使用扰动函数之后可以减少碰撞。

    ```JAVA
      static final int hash(Object key) {
          int h;
          // key.hashCode()：返回散列值也就是hashcode
          // ^ ：按位异或
          // >>>:无符号右移，忽略符号位，空位都以0补齐
          return (key == null) ? 0 : (h = key.hashCode()) ^ (h >>> 16);
      }


    ```

    `当链表长度大于阈值（默认为 8）时`，会首先调用 treeifyBin()方法。这个方法会根据 HashMap 数组来决定是否转换为红黑树。只有当`数组长度大于或者等于 64 的情况下，才会执行转换红黑树操作`，以减少搜索时间。否则，就是只是执行 resize() 方法对数组扩容。

    - loadFactor 加载因子
      loadFactor 加载因子是控制数组存放数据的疏密程度，loadFactor 越趋近于 1，那么 数组中存放的数据(entry)也就越多，也就越密，也就是会让链表的长度增加，loadFactor 越小，也就是趋近于 0，数组中存放的数据(entry)也就越少，也就越稀疏。loadFactor 太大导致查找元素效率低，太小导致数组的利用率低，存放的数据会很分散。loadFactor 的默认值为 0.75f 是官方给出的一个比较好的临界值。给定的默认容量为 16，负载因子为 0.75。Map 在使用过程中不断的往里面存放数据，当数量达到了 `16 * 0.75 = 12 `就需要将当前 16 的容量进行扩容，而扩容这个过程涉及到 rehash、复制数据等操作，所以非常消耗性能。

8.  ConcurrentHashMap 精华
    ConcurrnetHashMap 由很多个 Segment 组合，`而每一个 Segment 是一个类似于 HashMap 的结构`，所以每一个 HashMap 的内部可以进行扩容。但是 Segment 的个数一旦初始化就不能改变，默认 Segment 的个数是 `16` 个，你也可以认为 `ConcurrentHashMap 默认支持最多 16 个线程并发`。
    Java8 的 ConcurrentHashMap 相对于 Java7 来说变化比较大，不再是之前的 Segment 数组 + HashEntry 数组 + 链表，而是 Node 数组 + 链表 / 红黑树。当冲突链表达到一定长度时，链表会转换成红黑树。ConcurrentHashMap 使用的 Synchronized 锁加 CAS 的机制.

9.  无序性和不可重复性的含义是什么
    1、什么是无序性？无序性不等于随机性 ，`无序性是指存储的数据在底层数组中并非按照数组索引的顺序添加 ，而是根据数据的哈希值决定的`。
    2、什么是不可重复性？不可重复性是指添加的元素按照 equals()判断时 ，返回 false，需要同时重写 equals()方法和 HashCode()方法。

10. ArrayDeque 与 LinkedList 的区别
    ArrayDeque 是基于可变长的数组和双指针来实现，而 LinkedList 则通过链表来实现。
    ArrayDeque 插入时可能存在扩容过程, 不过均摊后的插入操作依然为 O(1)。虽然 LinkedList 不需要扩容，但是每次插入数据时均需要申请新的堆空间，均摊性能相比更慢。
11. 保证了 HashMap 总是使用 2 的幂作为哈希表的大小。

```JAVA
    /**
     * Returns a power of two size for the given target capacity.
     */
    static final int tableSizeFor(int cap) {
        int n = cap - 1;
        n |= n >>> 1;
        n |= n >>> 2;
        n |= n >>> 4;
        n |= n >>> 8;
        n |= n >>> 16;
        return (n < 0) ? 1 : (n >= MAXIMUM_CAPACITY) ? MAXIMUM_CAPACITY : n + 1;
    }

```
