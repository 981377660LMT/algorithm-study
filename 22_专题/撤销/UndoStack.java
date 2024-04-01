package template.algo;

import java.util.ArrayDeque;
import java.util.Deque;

public class UndoStack {
  private Deque<UndoOperation> dq;

  public UndoStack(int size) {
    dq = new ArrayDeque<>(size);
  }

  public void push(UndoOperation op) {
    dq.addLast(op);
    op.apply();
  }

  public UndoOperation pop() {
    UndoOperation ans = dq.removeLast();
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
      pop();
    }
  }
}
