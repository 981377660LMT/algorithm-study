/**
 * 四叉树
 * 00 01
 * 10 11
 */
public class QuadTree {
  private QuadTree s00, s01, s10, s11;
  private static QuadTree NIL = new QuadTree();

  private QuadTree() {
  }

  private void modify() {
      if (this == NIL) {
          return;
      }
  }

  private void pushUp() {
  }

  private void pushDown() {
      if (this == NIL) {
          return;
      }
  }

  public static QuadTree build(int x1, int x2, int y1, int y2) {
      if (x1 > x2 || y1 > y2) {
          return NIL;
      }
      QuadTree ans = new QuadTree();
      int xm = (x1 + x2) >> 1;
      int ym = (y1 + y2) >> 1;
      if (x1 < x2 || y1 < y2) {
          ans.s00 = build(x1, xm, y1, ym);
          ans.s01 = build(x1, xm, ym + 1, y2);
          ans.s10 = build(xm + 1, x2, y1, ym);
          ans.s11 = build(xm + 1, x2, ym + 1, y2);
          ans.pushUp();
      } else {
      }
      return ans;
  }

  void update(int tx1, int tx2, int ty1, int ty2, int x1, int x2, int y1,
              int y2) {
      if ((tx1 > x2 || tx2 < x1 || ty1 > y2 || ty2 < y1)) {
          return;
      }
      if ((tx1 <= x1 && x2 <= tx2 && ty1 <= y1 && y2 <= ty2)) {
          return;
      }
      int mx = (x1 + x2) >> 1;
      int my = (y1 + y2) >> 1;
      pushDown();
      s00.update(tx1, tx2, ty1, ty2, x1, mx, y1, my);
      s01.update(tx1, tx2, ty1, ty2, x1, mx, my + 1, y2);
      s10.update(tx1, tx2, ty1, ty2, mx + 1, x2, y1, my);
      s11.update(tx1, tx2, ty1, ty2, mx + 1, x2, my + 1, y2);
      pushUp();
  }

  void query(int tx1, int tx2, int ty1, int ty2, int x1, int x2, int y1,
             int y2) {
      if ((tx1 > x2 || tx2 < x1 || ty1 > y2 || ty2 < y1)) {
          return;
      }
      if ((tx1 <= x1 && x2 <= tx2 && ty1 <= y1 && y2 <= ty2)) {
          modify();
          return;
      }
      int mx = (x1 + x2) >> 1;
      int my = (y1 + y2) >> 1;
      pushDown();
      s00.query(tx1, tx2, ty1, ty2, x1, mx, y1, my);
      s01.query(tx1, tx2, ty1, ty2, x1, mx, my + 1, y2);
      s10.query(tx1, tx2, ty1, ty2, mx + 1, x2, y1, my);
      s11.query(tx1, tx2, ty1, ty2, mx + 1, x2, my + 1, y2);
  }
};
