package template.string;

import template.binary.Log2;


/**
 * It's a tree, each node represents a string, they support O(log_2 maxDepth) compare operation.
 */
public class StringTree {
  int l;
  int r;
  int log;
  TrieNode root;
  TrieNode max;
  TrieNode min;

  public TrieNode getRoot() {
    return root;
  }

  public TrieNode getMax() {
    return max;
  }

  public TrieNode getMin() {
    return min;
  }

  public StringTree(int l, int r, int maxDepth) {
    if (maxDepth == Integer.MAX_VALUE) {
      maxDepth--;
    }

    this.l = l;
    this.r = r;
    log = Log2.ceilLog(maxDepth + 1);
    root = newTrieNode();
    max = newTrieNode(root, 0);
    min = newTrieNode(root, 0);
    max.depth = maxDepth + 1;
    min.depth = -1;
  }

  private TrieNode newTrieNode() {
    TrieNode ans = new TrieNode();
    ans.next = new TrieNode[r - l + 1];
    ans.jump = new TrieNode[log + 1];
    return ans;
  }

  public TrieNode newTrieNode(TrieNode p, int c) {
    TrieNode node = newTrieNode();
    node.depth = p.depth + 1;
    node.value = c;
    node.jump[0] = p;
    for (int i = 0; node.jump[i] != null; i++) {
      node.jump[i + 1] = node.jump[i].jump[i];
    }
    return node;
  }

  public TrieNode insert(TrieNode root, int digit) {
    if (root.next[digit - l] == null) {
      root.next[digit - l] = newTrieNode(root, digit);
    }
    return root.next[digit - l];
  }

  public static class TrieNode {
    int depth;
    int value;
    TrieNode[] next;
    TrieNode[] jump;

    public int getValue() {
      return value;
    }

    public TrieNode getParent() {
      return jump[0];
    }

    public static int compare(TrieNode a, TrieNode b) {
      if (a.depth != b.depth) {
        return a.depth - b.depth;
      }
      if (a == b) {
        return 0;
      }
      TrieNode lca = lca(a, b);
      a = gotoDepth(a, lca.depth + 1);
      b = gotoDepth(b, lca.depth + 1);
      return a.value - b.value;
    }

    public static TrieNode lca(TrieNode a, TrieNode b) {
      if (a.depth > b.depth) {
        TrieNode tmp = a;
        a = b;
        b = tmp;
      }
      b = gotoDepth(b, a.depth);
      if (a == b) {
        return a;
      }
      for (int i = 20; i >= 0; i--) {
        if (a.jump[i] == b.jump[i]) {
          continue;
        }
        a = a.jump[i];
        b = b.jump[i];
      }
      return a.jump[0];
    }

    public static TrieNode gotoDepth(TrieNode t, int d) {
      if (t.depth == d) {
        return t;
      }
      int log = Log2.floorLog(t.depth - d);
      return gotoDepth(t.jump[log], d);
    }

    @Override
    public String toString() {
      return jump[0] == null ? "" : jump[0].toString() + value;
    }
  }
}
