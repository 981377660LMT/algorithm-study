package template.algo;

import template.datastructure.Treap;
import template.rand.RandomWrapper;

import java.util.function.IntFunction;

public class RandomWrapper {
  private Random random;

  public RandomWrapper() {
    this(new Random());
  }

  public RandomWrapper(Random random) {
    this.random = random;
  }

  public RandomWrapper(long seed) {
    this(new Random(seed));
  }

  public int nextInt(int l, int r) {
    return random.nextInt(r - l + 1) + l;
  }

  public int nextInt(int n) {
    return random.nextInt(n);
  }


  public double nextDouble(double l, double r) {
    return random.nextDouble() * (r - l) + l;
  }

  public double nextDouble() {
    return random.nextDouble();
  }

  public long nextLong(long l, long r) {
    return nextLong(r - l + 1) + l;
  }

  public long nextLong(long n) {
    return Math.round(random.nextDouble() * (n - 1));
  }

  public long nextLong() {
    return random.nextLong();
  }

  public String nextString(char l, char r, int len) {
    StringBuilder builder = new StringBuilder(len);
    for (int i = 0; i < len; i++) {
      builder.append((char) nextInt(l, r));
    }
    return builder.toString();
  }

  public String nextString(char[] s, int len) {
    StringBuilder builder = new StringBuilder(len);
    for (int i = 0; i < len; i++) {
      builder.append(s[nextInt(0, s.length - 1)]);
    }
    return builder.toString();
  }


  public Random getRandom() {
    return random;
  }

  public int range(int... x) {
    return x[nextInt(0, x.length - 1)];
  }

  public char range(char... x) {
    return x[nextInt(0, x.length - 1)];
  }

  public long range(long... x) {
    return x[nextInt(0, x.length - 1)];
  }

  public <T> T rangeT(T... x) {
    return x[nextInt(0, x.length - 1)];
  }

  public static final RandomWrapper INSTANCE = new RandomWrapper();

  public static void main(String[] args) {
    RandomWrapper random = new RandomWrapper();
    System.out.println(random.nextInt(1, 10));
    System.out.println(random.nextDouble(1.0, 10.0));
    System.out.println(random.nextLong(1, 10));
    System.out.println(random.nextString('a', 'z', 10));
    System.out.println(random.nextString(new char[] {'a', 'b', 'c'}, 10));
  }
}


/**
 * It's a special tree, which can maintain a set and perform below actions<br/>
 * <ul>
 * <li>insert multiple elements</li>
 * <li>pop top k largest elements, and reduce them by 1, then insert them back into set</li>
 * </ul>
 * All actions works in O(lg n) tc and O(1) mc<br/>
 * It's quite useful for some types of question
 *
 */
public class ReorderTree {
  Treap root = Treap.NIL;

  public void add(int x) {
    var s0 = Treap.splitByKey(root, x);
    Treap node = new Treap();
    node.key = x;
    node.pushUp();
    root = Treap.merge(s0[0], node);
    root = Treap.merge(root, s0[1]);
  }

  public boolean reduceTopK(int k) {
    if (root.size < k) {
      return false;
    }
    var s0 = Treap.splitByRank(root, root.size - k);
    var minKey = Treap.getKeyByRank(s0[1], 1);
    if (minKey <= 0) {
      root = Treap.merge(s0[0], s0[1]);
      return false;
    }
    var s1 = Treap.splitByKey(s0[0], minKey - 1);
    s0[1].modify(-1);
    var s2 = Treap.splitByKey(s0[1], minKey - 1);
    s0[0] = Treap.merge(s1[0], s2[0]);
    s0[1] = Treap.merge(s1[1], s2[1]);
    root = Treap.merge(s0[0], s0[1]);
    return true;
  }

  public String toString() {
    return root.toString();
  }

  static class Treap implements Cloneable {

    public static final Treap NIL = new Treap();

    static {
      NIL.left = NIL.right = NIL;
      NIL.size = 0;
    }

    public Treap left = NIL;
    public Treap right = NIL;
    public int size = 1;
    public int key;
    public int u;

    public static Treap build(IntFunction<Treap> func, int l, int r) {
      if (l > r) {
        return Treap.NIL;
      }
      int mid = (l + r) / 2;
      Treap root = func.apply(mid);
      root.left = build(func, l, mid - 1);
      root.right = build(func, mid + 1, r);
      root.pushUp();
      return root;
    }

    @Override
    public Treap clone() {
      try {
        return (Treap) super.clone();
      } catch (CloneNotSupportedException e) {
        throw new RuntimeException(e);
      }
    }

    public void modify(int x) {
      u += x;
      key += x;
    }


    public void pushDown() {
      if (this == NIL) {
        return;
      }
      if (u != 0) {
        left.modify(u);
        right.modify(u);
        u = 0;
      }
    }

    public void pushUp() {
      if (this == NIL) {
        return;
      }
      size = left.size + right.size + 1;
    }

    /**
     * split by rank and the node whose rank is argument will stored at result[0]
     */
    public static Treap[] splitByRank(Treap root, int rank) {
      if (root == NIL) {
        return new Treap[] {NIL, NIL};
      }
      root.pushDown();
      Treap[] result;
      if (root.left.size >= rank) {
        result = splitByRank(root.left, rank);
        root.left = result[1];
        result[1] = root;
      } else {
        result = splitByRank(root.right, rank - (root.size - root.right.size));
        root.right = result[0];
        result[0] = root;
      }
      root.pushUp();
      return result;
    }

    public static Treap merge(Treap a, Treap b) {
      if (a == NIL) {
        return b;
      }
      if (b == NIL) {
        return a;
      }
      if (RandomWrapper.INSTANCE.nextInt(a.size + b.size) < a.size) {
        a.pushDown();
        a.right = merge(a.right, b);
        a.pushUp();
        return a;
      } else {
        b.pushDown();
        b.left = merge(a, b.left);
        b.pushUp();
        return b;
      }
    }

    public static void toString(Treap root, StringBuilder builder) {
      if (root == NIL) {
        return;
      }
      root.pushDown();
      toString(root.left, builder);
      builder.append(root.key).append(',');
      toString(root.right, builder);
    }

    public static Treap clone(Treap root) {
      if (root == NIL) {
        return NIL;
      }
      Treap clone = root.clone();
      clone.left = clone(root.left);
      clone.right = clone(root.right);
      return clone;
    }

    @Override
    public String toString() {
      StringBuilder builder = new StringBuilder().append(key).append(":");
      toString(clone(this), builder);
      return builder.toString();
    }

    /**
     * nodes with key <= arguments will stored at result[0]
     */
    public static Treap[] splitByKey(Treap root, int key) {
      if (root == NIL) {
        return new Treap[] {NIL, NIL};
      }
      root.pushDown();
      Treap[] result;
      if (root.key > key) {
        result = splitByKey(root.left, key);
        root.left = result[1];
        result[1] = root;
      } else {
        result = splitByKey(root.right, key);
        root.right = result[0];
        result[0] = root;
      }
      root.pushUp();
      return result;
    }

    public static int getKeyByRank(Treap treap, int k) {
      while (treap.size > 1) {
        treap.pushDown();
        if (treap.left.size >= k) {
          treap = treap.left;
        } else {
          k -= treap.left.size;
          if (k == 1) {
            break;
          }
          k--;
          treap = treap.right;
        }
      }
      return treap.key;
    }

  }
}
