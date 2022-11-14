class Solution {

  static class FHQTreap {

    static int INF = 0x7FFFFFFF;
    static int count = 0;

    static class Node {

      Node[] ch;
      int val, sz, rk, dp, maxDp;

      public Node(int v, int r) {
        val = v;
        sz = 1;
        rk = (int) (Math.random() * 10000);
        ch = new Node[2];
        dp = maxDp = r;
      }

      public void up() {
        count++;
        sz = 1;
        if (ch[0] != null)
          sz += ch[0].sz;
        if (ch[1] != null)
          sz += ch[1].sz;
        // 注意不是maxDp，应该使用dp
        maxDp = Math.max(dp,
            Math.max(ch[0] == null ? 0 : ch[0].maxDp, ch[1] == null ? 0 : ch[1].maxDp));
      }

      public void down() {}// 需要时补充即可
    }

    Node ROOT;

    public Node[] split(Node root, int val) {
      if (root == null)
        return new Node[2];
      root.down();
      if (root.val <= val) {
        Node[] p = split(root.ch[1], val);
        root.ch[1] = p[0];
        root.up();
        return new Node[] {root, p[1]};
      } else {
        Node[] p = split(root.ch[0], val);
        root.ch[0] = p[1];
        root.up();
        return new Node[] {p[0], root};
      }
    }

    public Node merge(Node x, Node y) {
      if (x == null && y == null)
        return null;
      if (x == null || y == null)
        return x == null ? y : x;
      x.down();
      y.down();
      if (x.rk < y.rk) {
        x.ch[1] = merge(x.ch[1], y);
        x.up();
        return x;
      } else {
        y.ch[0] = merge(x, y.ch[0]);
        y.up();
        return y;
      }
    }

    public int insert(int val, int k) {
      Node[] t = split(ROOT, val - 1);
      Node[] p = split(t[0], val - k - 1);// [-inf,val - k - 1],[val-k,val - 1],[val,inf]
      int h = p[1] == null ? 0 : p[1].maxDp;
      t[0] = merge(p[0], p[1]);// 重新合并回来，只是为了求maxDp
      ROOT = merge(merge(t[0], new Node(val, h + 1)), t[1]);// 整颗树的根结点，插入的同时合并
      return h;
    }
  }

  public int lengthOfLIS(int[] nums, int k) {
    FHQTreap tr = new FHQTreap();
    int res = 0;
    for (int num : nums)
      res = Math.max(res, tr.insert(num, k) + 1);// h+1
    System.out.println("count = " + FHQTreap.count);
    return res;
  }


  public static void main(String[] args) {
    // nums := make([]int, 1e5)
    // for i := range nums {
    // if i <= 5e4 {
    // nums[i] = i + 5e4
    // } else {
    // nums[i] = 5e4
    // }
    // }

    int[] nums = new int[100000];
    for (int i = 0; i < nums.length; i++) {
      if (i <= 50000) {
        nums[i] = i + 50000;
      } else {
        nums[i] = 50000;
      }
    }
    int k = 50000;
    System.out.println(new Solution().lengthOfLIS(nums, k));
  }
}

