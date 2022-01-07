import java.util.LinkedList;
import java.util.Queue;


// push 的过程使用两个 queue
// 注意，提交给 Leetcode 的时候，需要将 MyStack3 改成是 MyStack
public class MyStack3 {

    private Queue<Integer> q;

    /** Initialize your data structure here. */
    public MyStack3() {
        q = new LinkedList<>();
    }

    /** Push element x onto stack. */
    public void push(int x) {

        Queue<Integer> q2 = new LinkedList<>();

//        while(!q.isEmpty())
//            q2.add(q.remove());
//
//        q.add(x);
//
//        while (!q2.isEmpty())
//            q.add(q2.remove());

        q2.add(x);
        while(!q.isEmpty())
            q2.add(q.remove());

        q = q2;
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
