import java.util.LinkedList;
import java.util.Queue;


// push 的过程只使用一个 queue
// 注意，提交给 Leetcode 的时候，需要将 MyStack4 改成是 MyStack
public class MyStack4 {

    private Queue<Integer> q;

    /** Initialize your data structure here. */
    public MyStack4() {
        q = new LinkedList<>();
    }

    /** Push element x onto stack. */
    public void push(int x) {

        // 首先，将 x 入队
        q.add(x);

        // 执行 n - 1 次出队再入队的操作
        for(int i = 1; i < q.size(); i ++)
            q.add(q.remove());
    }

    /** Removes the element on top of the stack and returns that element. */
    public int pop() {
        return q.remove();
    }

    /** Get the top element. */
    public int top() {
        return q.peek();
    }

    /** Returns whether the stack is empty. */
    public boolean empty() {
        return q.isEmpty();
    }
}
