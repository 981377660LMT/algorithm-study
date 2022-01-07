/// Leetcode 20. Valid Parentheses
/// https://leetcode.com/problems/valid-parentheses/description/
/// 括号匹配问题
///
/// 使用LinkedListStack测试20号问题的代码在视频中没有提及
/// 该代码主要用于使用Leetcode上的问题测试我们的LinkedListStack：）
class Solution {

    private class LinkedList<E> {

        private class Node{
            public E e;
            public Node next;

            public Node(E e, Node next){
                this.e = e;
                this.next = next;
            }

            public Node(E e){
                this(e, null);
            }

            public Node(){
                this(null, null);
            }

            @Override
            public String toString(){
                return e.toString();
            }
        }

        private Node dummyHead;
        private int size;

        public LinkedList(){
            dummyHead = new Node();
            size = 0;
        }

        // 获取链表中的元素个数
        public int getSize(){
            return size;
        }

        // 返回链表是否为空
        public boolean isEmpty(){
            return size == 0;
        }

        // 在链表的index(0-based)位置添加新的元素e
        // 在链表中不是一个常用的操作，练习用：）
        public void add(int index, E e){

            if(index < 0 || index > size)
                throw new IllegalArgumentException("Add failed. Illegal index.");

            Node prev = dummyHead;
            for(int i = 0 ; i < index ; i ++)
                prev = prev.next;

            prev.next = new Node(e, prev.next);
            size ++;
        }

        // 在链表头添加新的元素e
        public void addFirst(E e){
            add(0, e);
        }

        // 在链表末尾添加新的元素e
        public void addLast(E e){
            add(size, e);
        }

        // 获得链表的第index(0-based)个位置的元素
        // 在链表中不是一个常用的操作，练习用：）
        public E get(int index){

            if(index < 0 || index >= size)
                throw new IllegalArgumentException("Get failed. Illegal index.");

            Node cur = dummyHead.next;
            for(int i = 0 ; i < index ; i ++)
                cur = cur.next;
            return cur.e;
        }

        // 获得链表的第一个元素
        public E getFirst(){
            return get(0);
        }

        // 获得链表的最后一个元素
        public E getLast(){
            return get(size - 1);
        }

        // 修改链表的第index(0-based)个位置的元素为e
        // 在链表中不是一个常用的操作，练习用：）
        public void set(int index, E e){
            if(index < 0 || index >= size)
                throw new IllegalArgumentException("Update failed. Illegal index.");

            Node cur = dummyHead.next;
            for(int i = 0 ; i < index ; i ++)
                cur = cur.next;
            cur.e = e;
        }

        // 查找链表中是否有元素e
        public boolean contains(E e){
            Node cur = dummyHead.next;
            while(cur != null){
                if(cur.e.equals(e))
                    return true;
                cur = cur.next;
            }
            return false;
        }

        // 从链表中删除index(0-based)位置的元素, 返回删除的元素
        // 在链表中不是一个常用的操作，练习用：）
        public E remove(int index){
            if(index < 0 || index >= size)
                throw new IllegalArgumentException("Remove failed. Index is illegal.");

            // E ret = findNode(index).e; // 两次遍历

            Node prev = dummyHead;
            for(int i = 0 ; i < index ; i ++)
                prev = prev.next;

            Node retNode = prev.next;
            prev.next = retNode.next;
            retNode.next = null;
            size --;

            return retNode.e;
        }

        // 从链表中删除第一个元素, 返回删除的元素
        public E removeFirst(){
            return remove(0);
        }

        // 从链表中删除最后一个元素, 返回删除的元素
        public E removeLast(){
            return remove(size - 1);
        }

        // 从链表中删除元素e
        public void removeElement(E e){

            Node prev = dummyHead;
            while(prev.next != null){
                if(prev.next.e.equals(e))
                    break;
                prev = prev.next;
            }

            if(prev.next != null){
                Node delNode = prev.next;
                prev.next = delNode.next;
                delNode.next = null;
                size --;
            }
        }

        @Override
        public String toString(){
            StringBuilder res = new StringBuilder();

            Node cur = dummyHead.next;
            while(cur != null){
                res.append(cur + "->");
                cur = cur.next;
            }
            res.append("NULL");

            return res.toString();
        }
    }

    private interface Stack<E> {

        int getSize();
        boolean isEmpty();
        void push(E e);
        E pop();
        E peek();
    }

    private class LinkedListStack<E> implements Stack<E> {

        private LinkedList<E> list;

        public LinkedListStack(){
            list = new LinkedList<>();
        }

        @Override
        public int getSize(){
            return list.getSize();
        }

        @Override
        public boolean isEmpty(){
            return list.isEmpty();
        }

        @Override
        public void push(E e){
            list.addFirst(e);
        }

        @Override
        public E pop(){
            return list.removeFirst();
        }

        @Override
        public E peek(){
            return list.getFirst();
        }

        @Override
        public String toString(){
            StringBuilder res = new StringBuilder();
            res.append("Stack: top ");
            res.append(list);
            return res.toString();
        }
    }

    public boolean isValid(String s) {

        LinkedListStack<Character> stack = new LinkedListStack<>();
        for(int i = 0 ; i < s.length() ; i ++){
            char c = s.charAt(i);
            if(c == '(' || c == '[' || c == '{')
                stack.push(c);
            else{
                if(stack.isEmpty())
                    return false;

                char topChar = stack.pop();
                if(c == ')' && topChar != '(')
                    return false;
                if(c == ']' && topChar != '[')
                    return false;
                if(c == '}' && topChar != '{')
                    return false;
            }
        }
        return stack.isEmpty();
    }

    public static void main(String[] args) {

        System.out.println((new Solution()).isValid("()[]{}"));
        System.out.println((new Solution()).isValid("([)]"));
    }
}
