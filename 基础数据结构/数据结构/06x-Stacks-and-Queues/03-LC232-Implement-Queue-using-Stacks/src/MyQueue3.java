import java.util.Stack;


class MyQueue3 {

    private Stack<Integer> stack1;
    private Stack<Integer> stack2;
    int front;

    /** Initialize your data structure here. */
    public MyQueue3() {
        stack1 = new Stack<>();
        stack2 = new Stack<>();
    }

    /** Push element x to the back of queue. */
    public void push(int x) {

        if(stack1.isEmpty())
            front = x;
        stack1.push(x);
    }

    /** Removes the element from in front of queue and returns that element. */
    public int pop() {

        // 如果 stack2 不为空，直接返回 stack2 的栈首元素
        if(!stack2.isEmpty())
            return stack2.pop();

        while(stack1.size() > 1)
            stack2.push(stack1.pop());

        return stack1.pop();
    }

    /** Get the front element. */
    public int peek() {
        if(!stack2.isEmpty())
            return stack2.peek();
        return front;
    }

    /** Returns whether the queue is empty. */
    public boolean empty() {
        return stack1.isEmpty() && stack2.isEmpty();
    }
}