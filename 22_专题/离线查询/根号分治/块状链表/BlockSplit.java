package template.algo;

import template.math.DigitUtils;
import template.utils.Int2ToObjectFunction;
import template.utils.SegmentUtils;


public class SegmentUtils {
  public static boolean enter(int L, int R, int l, int r) {
    return L <= l && R >= r;
  }

  public static boolean leave(int L, int R, int l, int r) {
    return L > r || R < l;
  }
}


public class BlockSplit<S, U> {
  int B;
  int n;
  Block<S, U>[] blocks;



  public BlockSplit(int n, int blockSize, Int2ToObjectFunction<Block<S, U>> supplier) {
    B = blockSize;
    this.n = n;
    blocks = new Block[DigitUtils.ceilDiv(n, B)];
    for (int i = 0; i < blocks.length; i++) {
      int l = i * B;
      int r = l + B - 1;
      r = Math.min(r, n - 1);
      blocks[i] = supplier.apply(l, r);
    }
  }


  public void update(int L, int R, U upd) {
    for (int i = 0; i < blocks.length; i++) {
      int l = i * B;
      int r = l + B - 1;
      r = Math.min(r, n - 1);
      if (SegmentUtils.leave(L, R, l, r)) {
        continue;
      } else if (SegmentUtils.enter(L, R, l, r)) {
        blocks[i].fullyUpdate(upd);
      } else {
        blocks[i].beforePartialUpdate();
        for (int j = Math.max(l, L), to = Math.min(r, R); j <= to; j++) {
          blocks[i].partialUpdate(j - l, upd);
        }
        blocks[i].afterPartialUpdate();
      }
    }
  }

  public void update(int index, U upd) {
    int blockId = index / B;
    blocks[blockId].partialUpdate(index - blockId * B, upd);
    blocks[blockId].afterPartialUpdate();
  }

  public void query(int L, int R, S sum) {
    for (int i = 0; i < blocks.length; i++) {
      int l = i * B;
      int r = l + B - 1;
      r = Math.min(r, n - 1);
      if (SegmentUtils.leave(L, R, l, r)) {
        continue;
      } else if (SegmentUtils.enter(L, R, l, r)) {
        blocks[i].fullyQuery(sum);
      } else {
        blocks[i].beforePartialQuery();
        for (int j = Math.max(l, L), to = Math.min(r, R); j <= to; j++) {
          blocks[i].partialQuery(j - l, sum);
        }
        blocks[i].afterPartialQuery();
      }
    }
  }

  public static interface Block<S, U> {
    void fullyUpdate(U upd);

    void partialUpdate(int index, U upd);

    void fullyQuery(S s);

    void partialQuery(int index, S s);

    default void beforePartialUpdate() {}

    default void afterPartialUpdate() {}

    default void beforePartialQuery() {}

    default void afterPartialQuery() {}
  }
}
