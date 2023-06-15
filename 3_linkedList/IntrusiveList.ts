// !侵入式链表/侵入式容器  (intrusive list/container)
// 侵入式容器的一个重要特点是它不负责管理元素的生命周期，
// 仅仅是调整指针的链接，因此没有对象拷贝的运行开销，
// 元素的创建工作是外部做的，最小化了内存使用，提高运行效率。
//
// https://zhuanlan.zhihu.com/p/524894979
// http://quxiao.github.io/blog/2013/07/06/intrusive-list/
// https://blog.csdn.net/weixin_42663840/article/details/81188003
// https://blog.csdn.net/zuolj/article/details/78876591
// https://github.com/ionous/container/blob/master/inlist/list.go
// 一般的链表都是专门有一个链表数据结构，
// 然后这个数据结构中有各种元素以及前后向指针。
//
// 但在Linux内核中，链表数据结构并不包含任何实际元素，
// 就只有简单的前向后向指针。
// !非侵入式链表的实现方式是，链表节点中包含数据。侵入式链表的实现方式与之相反，是业务数据结构中包含链表节点结构：
// 然后在实际的特殊数据结构里包含这个链表数据结构(将链表节点包在数据里面)，
// 大大提高了链表的可复用性。
//
// struct list_node_t
// {
//     list_node_t* prev;
//     list_node_t* next;
// };
//
// struct MyClass
// {
//     int data;
//     list_node_t node;
// };
//
// 使用intrusive list实现，就可以省去了”节点 –> 数据指针 –> 数据”的二次查找。
// !找到intrusive list中的一个节点后，就可以立即找到这个节点对应的数据(地址偏移量)
// !通过从链表对象的内存地址中减去列表成员的偏移量来计算链表对象的基地址
// !C 语言中所有结构体内部数据都是按照约定好的 layout 排布的。我们可以通过 offsetof 计算偏移量来找到用户结构体的起始位置了
//
//
// !不能取地址的语言中如何实现??? -> 加一个owner 字段
// 侵入式就是把所有权交给外面的对象
// c++里可以根据地址偏移取到对象，别的语言只能保存一个外部数据的引用来取值

// !这种侵入式的组件是增强主体的
// 更好的 Data Locality
// 更友好的 API
// !脱离容器进行生命周期管理(同一份对象需要在多个容器中共享,将同一份数据加入多个链表中)

/**
 * 侵入式链表节点.
 * @see {@link https://zhuanlan.zhihu.com/p/524894979}
 */
class IntrusiveListNode<O> {
  prev: IntrusiveListNode<O> | undefined = undefined
  next: IntrusiveListNode<O> | undefined = undefined
  readonly owner: O

  constructor(owner: O) {
    this.owner = owner
  }

  /**
   * 在当前结点之后插入新节点`node`,并返回新节点。
   */
  insertAfter(node: IntrusiveListNode<O>): IntrusiveListNode<O> {
    node.prev = this
    node.next = this.next
    node.prev.next = node
    if (node.next) node.next.prev = node
    return node
  }

  /**
   * 在当前结点之前插入新节点`node`,并返回新节点。
   */
  insertBefore(node: IntrusiveListNode<O>): IntrusiveListNode<O> {
    node.next = this
    node.prev = this.prev
    node.next.prev = node
    if (node.prev) node.prev.next = node
    return node
  }

  /**
   * 从链表里移除自身.
   */
  remove(): IntrusiveListNode<O> {
    if (this.prev) this.prev.next = this.next
    if (this.next) this.next.prev = this.prev
    this.prev = void 0
    this.next = void 0
    return this
  }
}

/**
 * 原有用户数据.
 */
class UserData<V> {
  readonly value: V

  /* 侵入式链表 */
  readonly list: IntrusiveListNode<UserData<V>>

  constructor(value: V) {
    this.value = value
    this.list = new IntrusiveListNode(this)
  }

  doSomething() {
    console.log('do something')
  }

  prev(): UserData<V> | undefined {
    return this.list.prev?.owner
  }

  next(): UserData<V> | undefined {
    return this.list.next?.owner
  }

  insertAfter(node: UserData<V>): UserData<V> {
    return this.list.insertAfter(node.list).owner
  }

  insertBefore(node: UserData<V>): UserData<V> {
    return this.list.insertBefore(node.list).owner
  }
}

export {}
