package template.datastructure;

import template.binary.Log2;

/**
 * VEB-tree, support O(lglg u) insert, delete, member, predecessor, successor operation
 */
public class VanEmdeBoasTree {
  final static int NIL = -1;

  VanEmdeBoasTree[] cluster;
  VanEmdeBoasTree summary;


  int min = NIL;
  int max = NIL;
  byte b;

  int high(int i) {
    return i >> b;
  }

  int low(int i) {
    return i & ((1 << b) - 1);
  }

  int index(int i, int j) {
    return (i << b) + j;
  }

  private void toString(StringBuilder builder, int offset) {
    if (summary == null) {
      if (min != NIL) {
        builder.append(min + offset).append(',');
      }
      if (max != NIL && max != min) {
        builder.append(max + offset).append(',');
      }
      return;
    }
    for (int i = 0; i < cluster.length; i++) {
      cluster[i].toString(builder, offset + (i << b));
    }
  }

  @Override
  public String toString() {
    StringBuilder ans = new StringBuilder("{");
    toString(ans, 0);
    if (ans.length() > 1) {
      ans.setLength(ans.length() - 1);
    }
    ans.append("}");
    return ans.toString();
  }

  public static VanEmdeBoasTree newInstance(int u) {
    return new VanEmdeBoasTree(Log2.ceilLog(u));
  }

  private VanEmdeBoasTree(int log) {
    // memory_used++;
    if (log <= 1) {
      return;
    }
    b = (byte) (log / 2);
    cluster = new VanEmdeBoasTree[1 << (log - b)];
    // memory_used += cluster.length;
    for (int i = 0; i < cluster.length; i++) {
      cluster[i] = new VanEmdeBoasTree(b);
    }
    summary = new VanEmdeBoasTree(log - b);
  }


  public int minimum() {
    return min;
  }

  public int maximum() {
    return max;
  }

  public boolean member(int x) {
    if (x == min || x == max) {
      return true;
    } else if (summary == null) {
      return false;
    }
    return cluster[high(x)].member(low(x));
  }


  public int successor(int x) {
    if (summary == null) {
      if (x == 0 && max == 1) {
        return 1;
      }
      return NIL;
    } else if (min != NIL && x < min) {
      return min;
    } else {
      int maxlow = cluster[high(x)].max;
      if (maxlow != NIL && low(x) < maxlow) {
        int offset = cluster[high(x)].successor(low(x));
        return index(high(x), offset);
      } else {
        int succcluster = summary.successor(high(x));
        if (succcluster == NIL) {
          return NIL;
        } else {
          int offset = cluster[succcluster].min;
          return index(succcluster, offset);
        }
      }
    }
  }

  public int predecessor(int x) {
    if (summary == null) {
      if (x == 1 && min == 0) {
        return 0;
      }
      return NIL;
    } else if (max != NIL && x > max) {
      return max;
    } else {
      int minlow = cluster[high(x)].min;
      if (minlow != NIL && low(x) > minlow) {
        int offset = cluster[high(x)].predecessor(low(x));
        return index(high(x), offset);
      } else {
        int predcluster = summary.predecessor(high(x));
        if (predcluster == NIL) {
          if (min != NIL && x > min) {
            return min;
          } else {
            return NIL;
          }
        } else {
          int offset = cluster[predcluster].max;
          return index(predcluster, offset);
        }
      }
    }
  }

  void emptyTreeInsert(int x) {
    min = max = x;
  }

  public void insert(int x) {
    if (min == NIL) {
      emptyTreeInsert(x);
    } else {
      if (x < min) {
        int tmp = x;
        x = min;
        min = tmp;
      }
      if (summary != null) {
        if (cluster[high(x)].min == NIL) {
          summary.insert(high(x));
          cluster[high(x)].emptyTreeInsert(low(x));
        } else {
          cluster[high(x)].insert(low(x));
        }
      }
      if (x > max) {
        max = x;
      }
    }
  }

  public void delete(int x) {
    if (min == max) {
      min = max = NIL;
    } else if (summary == null) {
      if (x == 0) {
        min = 1;
      } else {
        min = 0;
      }
      max = min;
    } else {
      if (x == min) {
        int firstcluster = summary.min;
        x = index(firstcluster, cluster[firstcluster].min);
        min = x;
      }
      cluster[high(x)].delete(low(x));
      if (cluster[high(x)].min == NIL) {
        summary.delete(high(x));
        if (x == max) {
          int summarymax = summary.max;
          if (summarymax == NIL) {
            max = min;
          } else {
            max = index(summarymax, cluster[summarymax].max);
          }
        }
      } else if (x == max) {
        max = index(high(x), cluster[high(x)].max);
      }
    }
  }
}
