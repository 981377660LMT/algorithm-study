class Solution {

    public ListNode reverseList(ListNode head) {

        if(head == null || head.next == null)
            return head;

        ListNode rev = reverseList(head.next);
        head.next.next = head;
        head.next = null;
        return rev;
    }
}