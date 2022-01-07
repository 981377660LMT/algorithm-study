import java.util.LinkedList;
import java.util.Queue;


// 小优化，使用一个变量记录栈顶元素
// 注意，提交给 Leetcode 的时候，需要将 MyStack2 改成是 MyStack
public class MyStack2 {

    private Queue<Integer> q;
    private int top;

    /** Initialize your data structure here. */
    public MyStack2() {
        q = new LinkedList<>();
    }

    /** Push element x onto stack. */
    public void push(int x) {
        q.add(x);
        top = x;
    }

    /** Removes the element on top of the stack and returns that element. */
    public int pop() {

        Queue<Integer> q2 = new LinkedList<>();
        while (q.size() > 1){
            // 每从 q 中取出一个元素，都给 top 赋值
            // top 最后存储的就是 q 中除了队尾元素以外的最后一个元素
            // 即新的栈顶元素
            top = q.peek();
            q2.add(q.remove());
        }

        int ret = q.remove();
        q = q2;
        return ret;
    }

    /** Get the top element. */
    public int top() {
        return top;
    }

    /** Returns whether the stack is empty. */
    public boolean empty() {
        return q.isEmpty();
    }
}
