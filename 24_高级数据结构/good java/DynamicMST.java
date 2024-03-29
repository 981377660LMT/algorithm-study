package template.graph;

public class DynamicMST {
  private LCTNode[] nodes;
  private long totalWeight;
  private int edgeNum;

  public DynamicMST(int n) {
    nodes = new LCTNode[n];
    for (int i = 0; i < n; i++) {
      nodes[i] = new LCTNode();
      nodes[i].id = i;
      nodes[i].weight = 0;
      nodes[i].pushUp();
    }
    for (int i = 1; i < n; i++) {
      LCTNode node = new LCTNode();
      node.weight = Long.MAX_VALUE;
      node.a = nodes[i - 1];
      node.b = nodes[i];
      node.pushUp();
      LCTNode.join(node.a, node);
      LCTNode.join(node.b, node);
    }
  }

  public long getTotalWeight() {
    return totalWeight;
  }

  public int getEdgeNum() {
    return edgeNum;
  }

  public void addEdge(int aId, int bId, long weight) {
    LCTNode a = nodes[aId];
    LCTNode b = nodes[bId];
    LCTNode.findRoute(a, b);
    LCTNode.splay(a);
    if (a.largest.weight <= weight) {
      return;
    }
    LCTNode largest = a.largest;
    LCTNode.splay(largest);
    LCTNode.cut(largest.a, largest);
    LCTNode.cut(largest.b, largest);
    if (largest.weight < Long.MAX_VALUE) {
      edgeNum--;
      totalWeight -= largest.weight;
    }

    LCTNode node = new LCTNode();
    node.weight = weight;
    node.a = a;
    node.b = b;
    node.pushUp();
    LCTNode.join(node.a, node);
    LCTNode.join(node.b, node);
    edgeNum++;
    totalWeight += node.weight;
  }

  /**
   * 检查两个顶点之间是否存在一条路径
   */
  public boolean check(int aId, int bId) {
    return maxWeightBetween(aId, bId) < Long.MAX_VALUE;
  }

  public long maxWeightBetween(int aId, int bId) {
    LCTNode a = nodes[aId];
    LCTNode b = nodes[bId];
    LCTNode.findRoute(a, b);
    LCTNode.splay(b);
    return b.largest.weight;
  }

  private static class LCTNode {
    public static final LCTNode NIL = new LCTNode();

    static {
      NIL.left = NIL;
      NIL.right = NIL;
      NIL.father = NIL;
      NIL.treeFather = NIL;
      NIL.weight = 0;
      NIL.largest = NIL;
    }

    LCTNode left = NIL;
    LCTNode right = NIL;
    LCTNode father = NIL;
    LCTNode treeFather = NIL;
    boolean reverse;
    int id;

    LCTNode a;
    LCTNode b;
    LCTNode largest;
    long weight;

    public static LCTNode larger(LCTNode a, LCTNode b) {
      return a.weight >= b.weight ? a : b;
    }

    public static void access(LCTNode x) {
      LCTNode last = NIL;
      while (x != NIL) {
        splay(x);
        x.right.father = NIL;
        x.right.treeFather = x;
        x.setRight(last);
        x.pushUp();

        last = x;
        x = x.treeFather;
      }
    }

    public static void makeRoot(LCTNode x) {
      access(x);
      splay(x);
      x.reverse();
    }

    public static void cut(LCTNode y, LCTNode x) {
      makeRoot(y);
      access(x);
      splay(y);
      y.right.treeFather = NIL;
      y.right.father = NIL;
      y.setRight(NIL);
      y.pushUp();
    }

    public static void join(LCTNode y, LCTNode x) {
      makeRoot(x);
      x.treeFather = y;
    }

    public static void findRoute(LCTNode x, LCTNode y) {
      makeRoot(y);
      access(x);
    }

    public static void splay(LCTNode x) {
      if (x == NIL) {
        return;
      }
      LCTNode y, z;
      while ((y = x.father) != NIL) {
        if ((z = y.father) == NIL) {
          y.pushDown();
          x.pushDown();
          if (x == y.left) {
            zig(x);
          } else {
            zag(x);
          }
        } else {
          z.pushDown();
          y.pushDown();
          x.pushDown();
          if (x == y.left) {
            if (y == z.left) {
              zig(y);
              zig(x);
            } else {
              zig(x);
              zag(x);
            }
          } else {
            if (y == z.left) {
              zag(x);
              zig(x);
            } else {
              zag(y);
              zag(x);
            }
          }
        }
      }

      x.pushDown();
      x.pushUp();
    }

    public static void zig(LCTNode x) {
      LCTNode y = x.father;
      LCTNode z = y.father;
      LCTNode b = x.right;

      y.setLeft(b);
      x.setRight(y);
      z.changeChild(y, x);

      y.pushUp();
    }

    public static void zag(LCTNode x) {
      LCTNode y = x.father;
      LCTNode z = y.father;
      LCTNode b = x.left;

      y.setRight(b);
      x.setLeft(y);
      z.changeChild(y, x);

      y.pushUp();
    }

    public static LCTNode findRoot(LCTNode x) {
      x.pushDown();
      while (x.left != NIL) {
        x = x.left;
        x.pushDown();
      }
      splay(x);
      return x;
    }

    @Override
    public String toString() {
      return "" + id;
    }

    public void pushDown() {
      if (reverse) {
        reverse = false;

        LCTNode tmpNode = left;
        left = right;
        right = tmpNode;

        left.reverse();
        right.reverse();
      }

      left.treeFather = treeFather;
      right.treeFather = treeFather;
    }

    public void reverse() {
      reverse = !reverse;
    }

    public void setLeft(LCTNode x) {
      left = x;
      x.father = this;
    }

    public void setRight(LCTNode x) {
      right = x;
      x.father = this;
    }

    public void changeChild(LCTNode y, LCTNode x) {
      if (left == y) {
        setLeft(x);
      } else {
        setRight(x);
      }
    }

    public void pushUp() {
      largest = larger(this, left.largest);
      largest = larger(largest, right.largest);
    }
  }
}
