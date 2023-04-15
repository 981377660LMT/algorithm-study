template <class T> class partially_retroactive_queue {
  class node_type;

public:
  using time_type = typename std::list<node_type>::const_iterator;

private:
  class node_type {
    friend partially_retroactive_queue;

    T value;
    bool in_queue;
    node_type(const T value) : value(value), in_queue() {}
  };

  std::list<node_type> list;
  typename std::list<node_type>::iterator front_itr;

public:
  partially_retroactive_queue() : list(), front_itr(list.end()) {}

  time_type now() const { return list.cend(); }

  bool empty() const { return front_itr == list.end(); }

  T front() const {
    assert(!empty());
    return front_itr->value;
  }

  time_type insert_push(const time_type time, const T x) {
    // 将新节点插入到队列中
    const auto itr = list.insert(time, node_type(x));
    // 如果新节点是队列的第一个元素或者前一个元素不在队列中，那么新节点就不在队列中
    if (itr == list.begin() || !std::prev(itr)->in_queue) {
      itr->in_queue = false;
      // 将新节点的前一个节点作为队列的第一个节点
      --front_itr;
      front_itr->in_queue = true;
    } else {
      itr->in_queue = true;
    }
    return time_type(itr);
  }
  void erase_push(const time_type time) {
    assert(time != now());
    // 如果要删除的节点是队列的第一个元素或者前一个节点不在队列中，那么要删除的节点就不在队列中
    if (time == list.cbegin() || !std::prev(time)->in_queue) {
      front_itr->in_queue = false;
      ++front_itr;
    }
    // 将要删除的节点从队列中删除
    list.erase(time);
  }
  void insert_pop() {
    assert(!empty());
    front_itr->in_queue = false;
    ++front_itr;
  }
  void erase_pop() {
    assert(front_itr != list.begin());
    --front_itr;
    front_itr->in_queue = true;
  }
};