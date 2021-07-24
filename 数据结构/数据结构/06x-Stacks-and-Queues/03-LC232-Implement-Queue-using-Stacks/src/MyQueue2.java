import java.util.Stack;


class MyQueue2 {

    private Stack<Integer> stack;
    int front;

    /** Initialize your data structure here. */
    public MyQueue2() {
        stack = new Stack<>();
    }

    /** Push element x to the back of queue. */
    public void push(int x) {
        if(empty()) front = x;
        stack.push(x);
    }

    /** Removes the element from in front of queue and returns that element. */
    public int pop() {

        Stack<Integer> stack2 = new Stack<>();

        while(stack.size() > 1) {
            front = stack.peek();
            stack2.push(stack.pop());
        }

        int ret = stack.pop();

        while(!stack2.isEmpty())
            stack.push(stack2.pop());

        return ret;
    }

    /** Get the front element. */
    public int peek() {
        return front;
    }

    /** Returns whether the queue is empty. */
    public boolean empty() {
        return stack.isEmpty();
    }
}