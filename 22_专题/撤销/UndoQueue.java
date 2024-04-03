package template.algo;

import java.util.ArrayDeque;
import java.util.Deque;


/**
 * (a, b) take the same effect as (b, a)
 */
public class FlaggedCommutativeUndoOperation implements UndoOperation {
  boolean flag;
  UndoOperation op;

  public FlaggedCommutativeUndoOperation(UndoOperation op) {
    this.op = op;
  }

  @Override
  public void apply() {
    op.apply();
  }

  @Override
  public void undo() {
    op.undo();
  }
}


/**
 * each operation will be invoked and undo O(log_2n) times
 */
public class UndoQueue {


  private Deque<FlaggedCommutativeUndoOperation> dq;
  private Deque<FlaggedCommutativeUndoOperation> bufA;
  private Deque<FlaggedCommutativeUndoOperation> bufB;

  public UndoQueue(int size) {
    dq = new ArrayDeque<>(size);
    bufA = new ArrayDeque<>(size);
    bufB = new ArrayDeque<>(size);
  }

  public void add(FlaggedCommutativeUndoOperation op) {
    op.flag = false;
    pushAndDo(op);
  }

  private void pushAndDo(FlaggedCommutativeUndoOperation op) {
    dq.addLast(op);
    op.apply();
  }

  

  private void popAndUndo() {
    FlaggedCommutativeUndoOperation ans = dq.removeLast();
    ans.undo();
    if (ans.flag) {
      bufA.addLast(ans);
    } else {
      bufB.addLast(ans);
    }
  }

  public FlaggedCommutativeUndoOperation remove() {
    if (!dq.peekLast().flag) {
      popAndUndo();
      while (!dq.isEmpty() && bufB.size() != bufA.size()) {
        popAndUndo();
      }
      if (dq.isEmpty()) {
        while (!bufB.isEmpty()) {
          FlaggedCommutativeUndoOperation ans = bufB.removeFirst();
          ans.flag = true;
          pushAndDo(ans);
        }
      } else {
        while (!bufB.isEmpty()) {
          FlaggedCommutativeUndoOperation ans = bufB.removeLast();
          pushAndDo(ans);
        }
      }
      while (!bufA.isEmpty()) {
        pushAndDo(bufA.removeLast());
      }
    }

    FlaggedCommutativeUndoOperation ans = dq.removeLast();
    ans.undo();
    return ans;
  }

  public int size() {
    return dq.size();
  }

  public boolean isEmpty() {
    return dq.isEmpty();
  }

  public void clear() {
    while (!isEmpty()) {
      popAndUndo();
    }
    bufA.clear();
    bufB.clear();
  }
}
