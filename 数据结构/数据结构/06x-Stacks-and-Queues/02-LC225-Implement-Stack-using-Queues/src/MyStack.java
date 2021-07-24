import java.util.*;


class MyStack {

    private Queue<Integer> q;

    /** Initialize your data structure here. */
    public MyStack() {
        q = new LinkedList<>();
    }

    /** Push element x onto stack. */
    public void push(int x) {
        q.add(x);
    }

    /** Removes the element on top of the stack and returns that element. */
    public int pop() {

        // 创建另外一个队列 q2
        Queue<Integer> q2 = new LinkedList<>();

        // 除了最后一个元素，将 q 中的所有元素放入 q2
        while (q.size() > 1)
            q2.add(q.remove());

        // q 中剩下的最后一个元素就是“栈顶”元素
        int ret = q.remove();

        // 此时 q2 是整个数据结构存储的所有其他数据，赋值给 q
        q = q2;

        // 返回“栈顶元素”
        return ret;
    }

    /** Get the top element. */
    public int top() {
        int ret = pop();
        push(ret);
        return ret;
    }

    /** Returns whether the stack is empty. */
    public boolean empty() {
        return q.isEmpty();
    }
}

/**
 * Your MyStack object will be instantiated and called as such:
 * MyStack obj = new MyStack();
 * obj.push(x);
 * int param_2 = obj.pop();
 * int param_3 = obj.top();
 * boolean param_4 = obj.empty();
 */