// !复杂度有问题

import java.util.function.Supplier;
import java.util.Objects;

interface Int2ToObjectFunction<T> {
  T apply(int a, int b);
}

class Pair<A, B> implements Cloneable {
  public A a;
  public B b;

  public Pair(A a, B b) {
    this.a = a;
    this.b = b;
  }

  @Override
  public String toString() {
    return "a=" + a + ",b=" + b;
  }

  @Override
  public int hashCode() {
    return a.hashCode() * 31 + b.hashCode();
  }

  @Override
  public boolean equals(Object obj) {
    Pair<A, B> casted = (Pair<A, B>) obj;
    return Objects.equals(casted.a, a) && Objects.equals(casted.b, b);
  }

  @Override
  public Pair<A, B> clone() {
    try {
      return (Pair<A, B>) super.clone();
    } catch (CloneNotSupportedException e) {
      throw new RuntimeException(e);
    }
  }
}

class SegmentUtils {
  public static boolean enter(int L, int R, int l, int r) {
    return L <= l && R >= r;
  }

  public static boolean leave(int L, int R, int l, int r) {
    return L > r || R < l;
  }
}

/**
 * all operation in this framework take O(N/B) while N is the maximum size and B
 * is the block size
 *
 * @param <S>
 * @param <U>
 * @param <E>
 */
class BlockChain<S, U, E, B extends BlockChain.Block<S, U, E, B>> {
  LinkedNode<B> head = new LinkedNode<>();
  LinkedNode<B> tail = new LinkedNode<>();
  int B;
  int size;

  public BlockChain(int B, Supplier<B> supplier) {
    // add an empty node
    B block = supplier.get();
    LinkedNode<B> node = new LinkedNode<>();
    node.data = block;
    node.size = 0;
    this.B = B;
    LinkedNode.link(tail.prev, node);
    LinkedNode.link(node, tail);
  }

  private BlockChain(int B, LinkedNode<B> begin, LinkedNode<B> end) {
    // add an empty node
    this.B = B;
    LinkedNode.link(head, begin);
    LinkedNode.link(end, tail);
    maintain();
  }

  public BlockChain(int n, int B, Int2ToObjectFunction<B> supplier) {
    assert n > 0;
    this.B = B;
    size = n;
    LinkedNode.link(head, tail);
    for (int i = 0; i < n; i += B) {
      int l = i;
      int r = i + B - 1;
      r = Math.min(r, n - 1);
      B block = supplier.apply(l, r);
      LinkedNode<B> node = new LinkedNode<>();
      node.data = block;
      node.size = r - l + 1;
      LinkedNode.link(tail.prev, node);
      LinkedNode.link(node, tail);
    }
  }

  void split(LinkedNode<B> node, int n) {
    // split if necessary
    LinkedNode<B> post = new LinkedNode<>();
    Pair<B, B> pair = node.data.split(n);
    LinkedNode.link(post, node.next);
    LinkedNode.link(node, post);
    post.data = pair.b;
    post.size = node.size - n;
    node.data = pair.a;
    node.size = n;
  }

  void afterInsertion(LinkedNode<B> node) {
    if (node.size >= 2 * B) {
      split(node, B);
    }
  }

  void mergeNode(LinkedNode<B> a, LinkedNode<B> b) {
    LinkedNode.link(a, b.next);
    a.data = a.data.merge(b.data);
    a.size += b.size;
  }

  void maintain() {
    size = 0;
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      size += node.size;
      if (node.size >= 2 * B) {
        split(node, B);
      } else if (node.prev != head && node.size + node.prev.size <= B) {
        mergeNode(node.prev, node);
      }
    }
  }

  private LinkedNode<B> findKth(int k) {
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.size < k) {
        k -= node.size;
        continue;
      }
      return node;
    }
    return tail;
  }

  private LinkedNode<B> splitKth(int k) {
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.size < k) {
        k -= node.size;
        continue;
      }
      if (k != 1) {
        split(node, k - 1);
        node = node.next;
      }
      return node;
    }
    return tail;
  }

  public E get(int index) {
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.size <= index) {
        index -= node.size;
        continue;
      }
      node.data.beforePartialQuery();
      return node.data.get(index);
    }
    throw new IndexOutOfBoundsException();
  }

  public int prefixSize(B block, boolean include) {
    int ans = 0;
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.data == block) {
        if (include) {
          ans += node.size;
        }
        break;
      }
      ans += node.size;
    }
    return ans;
  }

  /**
   * k-th in res.b
   *
   * @param k
   * @return
   */
  public Pair<BlockChain<S, U, E, B>, BlockChain<S, U, E, B>> split(int k, Supplier<B> supplier) {
    if (k == 1) {
      return new Pair<>(new BlockChain<>(B, supplier), this);
    }
    if (k > size) {
      return new Pair<>(this, new BlockChain<>(B, supplier));
    }
    LinkedNode<B> head = splitKth(k);
    LinkedNode<B> end = tail.prev;
    LinkedNode.link(head.prev, tail);
    BlockChain<S, U, E, B> b = new BlockChain<>(B, head, end);
    maintain();
    return new Pair<>(this, b);
  }

  /**
   * insert after the index-th elements, 0 means first
   *
   * @param index
   * @param e
   */
  public void insert(int index, E e) {
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.size < index) {
        index -= node.size;
        continue;
      }
      node.data.insert(index, e);
      node.size++;
      break;
    }
    maintain();
  }

  public void delete(int index) {
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      if (node.size <= index) {
        index -= node.size;
        continue;
      }
      node.data.delete(index);
      node.size--;
      break;
    }
    maintain();
  }

  @Override
  public String toString() {
    StringBuilder ans = new StringBuilder("[");
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      ans.append("<").append(node.data).append(">,");
    }
    if (ans.charAt(ans.length() - 1) == ',') {
      ans.setLength(ans.length() - 1);
    }
    ans.append("]");
    return ans.toString();
  }

  public void update(int L, int R, U upd) {
    int offset = 0;
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      int l = offset;
      int r = offset + node.size - 1;
      offset += node.size;
      if (SegmentUtils.enter(L, R, l, r)) {
        node.data.fullyUpdate(upd);
      } else if (SegmentUtils.leave(L, R, l, r)) {
        continue;
      } else {
        for (int i = Math.max(l, L), to = Math.min(r, R); i <= to; i++) {
          node.data.partialUpdate(i - l, upd);
        }
        node.data.afterPartialUpdate();
      }
    }
  }

  public void query(int L, int R, S sum) {
    int offset = 0;
    for (LinkedNode<B> node = head.next; node != tail; node = node.next) {
      int l = offset;
      int r = offset + node.size - 1;
      offset += node.size;
      if (SegmentUtils.enter(L, R, l, r)) {
        node.data.fullyQuery(sum);
      } else if (SegmentUtils.leave(L, R, l, r)) {
        continue;
      } else {
        node.data.beforePartialQuery();
        for (int i = Math.max(l, L), to = Math.min(r, R); i <= to; i++) {
          node.data.partialQuery(i - l, sum);
        }
      }
    }
  }

  /**
   * 将第k个元素旋转到第一的位置
   */
  public void rotate(int k) {
    LinkedNode<B> node = splitKth(k);

    // rotate
    LinkedNode<B> h1 = head.next;
    LinkedNode<B> e1 = node.prev;
    LinkedNode<B> h2 = node;
    LinkedNode<B> e2 = node.next;

    LinkedNode.link(head, h2);
    LinkedNode.link(e2, h1);
    LinkedNode.link(e1, tail);

    maintain();
  }

  private void reverse(LinkedNode<B> root, LinkedNode<B> p) {
    if (root == null) {
      return;
    }
    reverse(root.next, root);
    root.data.reverse();
    root.prev = root.next;
    root.next = p;
  }

  public void reverse(int l, int r) {
    LinkedNode<B> left = splitKth(l + 1);
    LinkedNode<B> right = splitKth(r + 2).prev;

    LinkedNode<B> begin = left.prev;
    LinkedNode<B> end = right.next;

    right.next = null;
    reverse(left, null);

    LinkedNode.link(begin, right);
    LinkedNode.link(left, end);

    maintain();
  }

  public int size() {
    return size;
  }

  public static interface Block<S, U, E, B extends Block<S, U, E, B>> {
    /**
     * insert after the index-th elements, 0 means first.
     *
     * factor 1
     */
    default void insert(int index, E e) {
      throw new UnsupportedOperationException();
    }

    /**
     * factor 1
     */
    default void delete(int index) {
      throw new UnsupportedOperationException();
    }

    /**
     * factor 1
     */
    default E get(int index) {
      throw new UnsupportedOperationException();
    }

    /**
     * res.a contain first n elements, and res.b contains others
     * after split, you can release this as you like
     *
     * factor 1
     */
    Pair<B, B> split(int n);

    /**
     * after merge, you can release block as you like
     *
     * factor 1
     * 
     * @param block
     * @return
     */
    B merge(B block);

    /**
     * factor 1
     */
    default void reverse() {
      throw new UnsupportedOperationException();
    }

    /**
     * factor n/B
     */
    void fullyQuery(S sum);

    /**
     * factor B
     */
    void partialQuery(int index, S sum);

    /**
     * factor n/B
     */
    void fullyUpdate(U upd);

    /**
     * factor B
     */
    void partialUpdate(int index, U upd);

    /**
     * factor 1
     */
    default void afterPartialUpdate() {

    }

    /**
     * factor 1
     */
    default void beforePartialQuery() {

    }

  }

  private static class LinkedNode<E> {
    LinkedNode<E> prev;
    LinkedNode<E> next;
    int size;
    E data;

    static <T> void link(LinkedNode<T> a, LinkedNode<T> b) {
      b.prev = a;
      a.next = b;
    }

    static <T> void cut(LinkedNode<T> a, LinkedNode<T> b) {
      a.next = null;
      b.prev = null;
    }

    @Override
    public String toString() {
      return "" + data;
    }
  }
}