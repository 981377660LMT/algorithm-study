package template.algo;

import java.util.*;


public class PriorityCommutativeUndoOperation implements UndoOperation {
  public long priority;
  int offsetToBottom;
  UndoOperation op;

  public PriorityCommutativeUndoOperation(long priority, UndoOperation op) {
    this.priority = priority;
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
 * long type unique priority, each operation will be done and undone at most ln n times
 *
 * credited to https://codeforces.com/blog/entry/111117
 */
public class UndoPriorityQueue {
  TreeSet<PriorityCommutativeUndoOperation> set =
      new TreeSet<>(Comparator.comparingLong(x -> x.priority));
  UndoStack stack;
  List<PriorityCommutativeUndoOperation> bufferForLowPriority;
  List<PriorityCommutativeUndoOperation> bufferForHighPriority;

  public UndoPriorityQueue(int cap) {
    stack = new UndoStack(cap);
    bufferForLowPriority = new ArrayList<>(cap);
    bufferForHighPriority = new ArrayList<>(cap);
  }

  public void push(PriorityCommutativeUndoOperation op) {
    if (!set.add(op)) {
      throw new IllegalArgumentException("Duplicate priority");
    }
    pushStack(op);
  }

  private void pushStack(PriorityCommutativeUndoOperation op) {
    op.offsetToBottom = stack.size();
    stack.push(op);
  }

  public PriorityCommutativeUndoOperation pop() {
    int k = 0;
    int size = size();
    bufferForLowPriority.clear();
    for (PriorityCommutativeUndoOperation op : set) {
      bufferForLowPriority.add(op);
      k = Math.max(k, size - op.offsetToBottom);
      op.offsetToBottom = -1;
      if (bufferForLowPriority.size() * 2 >= k) {
        break;
      }
    }
    if (k > 1) {
      bufferForHighPriority.clear();
      for (int i = 0; i < k; i++) {
        PriorityCommutativeUndoOperation op = (PriorityCommutativeUndoOperation) stack.pop();
        if (op.offsetToBottom != -1) {
          bufferForHighPriority.add(op);
        }
      }
      for (PriorityCommutativeUndoOperation op : bufferForHighPriority) {
        pushStack(op);
      }
      Collections.reverse(bufferForLowPriority);
      for (PriorityCommutativeUndoOperation op : bufferForLowPriority) {
        pushStack(op);
      }
    }
    PriorityCommutativeUndoOperation ans = set.pollFirst();
    stack.pop();
    return ans;
  }

  public int size() {
    return stack.size();
  }

  public boolean isEmpty() {
    return stack.isEmpty();
  }
}
