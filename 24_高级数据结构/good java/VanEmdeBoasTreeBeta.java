package template.datastructure;

import template.binary.Log2;

import java.util.Arrays;
import java.util.function.IntConsumer;

/**
 * VEB-tree, support O(lglg u) insert, delete, member, predecessor, successor operation
 */
public class VanEmdeBoasTreeBeta {
  final static int NIL = -1;

  int[] cluster;
  int[] start;
  int[] summary;
  int[] min;
  int[] max;
  byte[] b;
  byte[] h;
  int clusterTail;
  int nodeTail;

  private int allocNode() {
    return nodeTail++;
  }

  private int allocCluster(int len) {
    int ans = clusterTail;
    clusterTail += len;
    assert clusterTail <= start.length;
    return ans;
  }

  int high(int V, int i) {
    return i >> b[V];
  }

  int low(int V, int i) {
    return i & ((1 << b[V]) - 1);
  }

  int index(int V, int i, int j) {
    return (i << b[V]) + j;
  }

  int clusterLength(int V) {
    return 1 << h[V];
  }

  public void visit(IntConsumer consumer) {
    visit(0, consumer, 0);
  }

  private void visit(int V, IntConsumer consumer, int offset) {
    if (min[V] != NIL) {
      consumer.accept(min[V] + offset);
    } else {
      return;
    }
    if (summary[V] == NIL) {
      if (max[V] != NIL && max[V] != min[V]) {
        consumer.accept(max[V] + offset);
      }
      return;
    }

    for (int i = 0; i < clusterLength(V); i++) {
      visit(cluster[i + start[V]], consumer, offset + (i << b[V]));
    }
  }

  @Override
  public String toString() {
    StringBuilder ans = new StringBuilder("{");
    visit(0, i -> ans.append(i).append(','), 0);
    if (ans.length() > 1) {
      ans.setLength(ans.length() - 1);
    }
    ans.append("}");
    return ans.toString();
  }

  public static VanEmdeBoasTreeBeta newInstance(int u) {
    return new VanEmdeBoasTreeBeta(Log2.ceilLog(u));
  }

  int nodeCnt;
  int clusterCnt;

  private void estimate(int log) {
    nodeCnt++;
    if (log <= 1) {
      return;
    }
    int b = (log / 2);
    int len = 1 << (log - b);
    clusterCnt += len;
    // memory_used += cluster.length;
    for (int i = 0; i < len; i++) {
      estimate(b);
    }
    estimate(log - b);
  }

  private int build(int log) {
    int V = allocNode();
    if (log <= 1) {
      return V;
    }
    b[V] = (byte) (log / 2);
    h[V] = (byte) (log - b[V]);
    int len = 1 << h[V];
    start[V] = allocCluster(len);
    // memory_used += cluster.length;
    for (int i = 0; i < len; i++) {
      cluster[start[V] + i] = build(b[V]);
    }
    summary[V] = build(h[V]);
    return V;
  }

  private VanEmdeBoasTreeBeta(int log) {
    estimate(log);
    cluster = new int[clusterCnt];
    start = new int[nodeCnt];
    summary = new int[nodeCnt];
    min = new int[nodeCnt];
    max = new int[nodeCnt];
    b = new byte[nodeCnt];
    h = new byte[nodeCnt];
    Arrays.fill(min, -1);
    Arrays.fill(max, -1);
    Arrays.fill(summary, -1);

    build(log);
    assert clusterTail == cluster.length;
    assert nodeTail == start.length;
  }


  public int minimum() {
    return min[0];
  }

  public int maximum() {
    return max[0];
  }


  public boolean member(int x) {
    return member(0, x);
  }

  private boolean member(int V, int x) {
    if (x == min[V] || x == max[V]) {
      return true;
    } else if (summary[V] == NIL) {
      return false;
    }
    return member(cluster[high(V, x) + start[V]], low(V, x));
  }

  public int successor(int x) {
    return successor(0, x);
  }

  private int successor(int V, int x) {
    if (summary[V] == NIL) {
      if (x == 0 && max[V] == 1) {
        return 1;
      }
      return NIL;
    } else if (min[V] != NIL && x < min[V]) {
      return min[V];
    } else {
      int maxlow = max[cluster[high(V, x) + start[V]]];
      if (maxlow != NIL && low(V, x) < maxlow) {
        int offset = successor(cluster[high(V, x) + start[V]], low(V, x));
        return index(V, high(V, x), offset);
      } else {
        int succcluster = successor(summary[V], high(V, x));
        if (succcluster == NIL) {
          return NIL;
        } else {
          int offset = min[cluster[succcluster + start[V]]];
          return index(V, succcluster, offset);
        }
      }
    }
  }

  public int predecessor(int x) {
    return predecessor(0, x);
  }

  public int predecessor(int V, int x) {
    if (summary[V] == NIL) {
      if (x == 1 && min[V] == 0) {
        return 0;
      }
      return NIL;
    } else if (max[V] != NIL && x > max[V]) {
      return max[V];
    } else {
      int minlow = min[cluster[high(V, x) + start[V]]];
      if (minlow != NIL && low(V, x) > minlow) {
        int offset = predecessor(cluster[high(V, x) + start[V]], low(V, x));
        return index(V, high(V, x), offset);
      } else {
        int predcluster = predecessor(summary[V], high(V, x));
        if (predcluster == NIL) {
          if (min[V] != NIL && x > min[V]) {
            return min[V];
          } else {
            return NIL;
          }
        } else {
          int offset = max[cluster[predcluster + start[V]]];
          return index(V, predcluster, offset);
        }
      }
    }
  }

  void emptyTreeInsert(int V, int x) {
    min[V] = max[V] = x;
  }

  public void insert(int x) {
    insert(0, x);
  }

  private void insert(int V, int x) {
    if (min[V] == NIL) {
      emptyTreeInsert(V, x);
    } else {
      if (x < min[V]) {
        int tmp = x;
        x = min[V];
        min[V] = tmp;
      }
      if (summary[V] != NIL) {
        if (min[cluster[high(V, x) + start[V]]] == NIL) {
          insert(summary[V], high(V, x));
          emptyTreeInsert(cluster[high(V, x) + start[V]], low(V, x));
        } else {
          insert(cluster[high(V, x) + start[V]], low(V, x));
        }
      }
      if (x > max[V]) {
        max[V] = x;
      }
    }
  }

  public void delete(int x) {
    delete(0, x);
  }

  private void delete(int V, int x) {
    if (min[V] == max[V]) {
      min[V] = max[V] = NIL;
    } else if (summary[V] == NIL) {
      if (x == 0) {
        min[V] = 1;
      } else {
        min[V] = 0;
      }
      max[V] = min[V];
    } else {
      if (x == min[V]) {
        int firstcluster = min[summary[V]];
        x = index(V, firstcluster, min[cluster[firstcluster + start[V]]]);
        min[V] = x;
      }
      delete(cluster[high(V, x) + start[V]], low(V, x));
      if (min[cluster[high(V, x) + start[V]]] == NIL) {
        delete(summary[V], high(V, x));
        if (x == max[V]) {
          int summarymax = max[summary[V]];
          if (summarymax == NIL) {
            max[V] = min[V];
          } else {
            max[V] = index(V, summarymax, max[cluster[summarymax + start[V]]]);
          }
        }
      } else if (x == max[V]) {
        max[V] = index(V, high(V, x), max[cluster[high(V, x) + start[V]]]);
      }
    }
  }
}
