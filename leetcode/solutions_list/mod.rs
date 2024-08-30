// as_ref 提供了一种安全且方便的方式来访问 Option 中的值，而不必担心所有权问题.

struct Solution;

#[derive(PartialEq, Eq, Clone, Debug)]
// Definition for singly-linked list.
pub struct ListNode {
    pub val: i32,
    pub next: Option<Box<ListNode>>,
}

impl ListNode {
    #[inline]
    fn new(val: i32) -> Self {
        ListNode { next: None, val }
    }
}

impl Solution {
    pub fn merge_two_lists(
        l1: Option<Box<ListNode>>,
        l2: Option<Box<ListNode>>,
    ) -> Option<Box<ListNode>> {
        match (l1, l2) {
            (None, None) => None,
            (None, r) => r,
            (l, None) => l,
            (Some(mut l), Some(mut r)) => {
                if l.val <= r.val {
                    l.next = Self::merge_two_lists(l.next, Some(r));
                    Some(l)
                } else {
                    r.next = Self::merge_two_lists(Some(l), r.next);
                    Some(r)
                }
            }
        }
    }

    pub fn merge_two_lists_iteratively(
        mut l1: Option<Box<ListNode>>,
        mut l2: Option<Box<ListNode>>,
    ) -> Option<Box<ListNode>> {
        let mut dummy = ListNode::new(-1);
        let mut cur = &mut dummy;
        // as_ref 提供了一种安全且方便的方式来访问 Option 中的值，而不必担心所有权问题
        while let (Some(n1), Some(n2)) = (l1.as_ref(), l2.as_ref()) {
            if n1.val < n2.val {
                // 将较小链表连接到新链表尾节点，所有权移动
                cur.next = l1;
                // 将cur尾节点指向它的后继节点
                cur = cur.next.as_mut().unwrap();
                // 将链表从尾节点取下来，将所有权返给较小的链表
                l1 = cur.next.take();
            } else {
                cur.next = l2;
                cur = cur.next.as_mut().unwrap();
                l2 = cur.next.take();
            }
        }

        cur.next = if l1.is_some() { l1 } else { l2 };
        dummy.next
    }

    pub fn reverse_list(head: Option<Box<ListNode>>) -> Option<Box<ListNode>> {
        let mut old_list = head;
        let mut new_list: Option<Box<ListNode>> = None;
        while old_list.is_some() {
            let mut head = old_list.take().unwrap();
            let rem = head.next.take();
            head.next = new_list;
            new_list = Some(head);
            old_list = rem;
        }
        new_list
    }

    pub fn add_two_numbers(
        mut l1: Option<Box<ListNode>>,
        mut l2: Option<Box<ListNode>>,
    ) -> Option<Box<ListNode>> {
        let mut dummy = ListNode::new(0);
        let mut cur = &mut dummy;
        let mut carry = 0;
        while l1.is_some() || l2.is_some() {
            let v1 = l1.as_ref().map_or(0, |n| n.val);
            let v2 = l2.as_ref().map_or(0, |n| n.val);
            let sum = v1 + v2 + carry;
            cur.next = Some(Box::new(ListNode::new(sum % 10)));
            carry = sum / 10;
            l1 = l1.and_then(|n| n.next);
            l2 = l2.and_then(|n| n.next);
            cur = cur.next.as_mut().unwrap();
        }
        if carry > 0 {
            cur.next = Some(Box::new(ListNode::new(carry)));
        }
        dummy.next
    }

    pub fn add_two_numbers_2(
        l1: Option<Box<ListNode>>,
        l2: Option<Box<ListNode>>,
    ) -> Option<Box<ListNode>> {
        fn to_vec(mut l: Option<Box<ListNode>>) -> Vec<i32> {
            let mut res = vec![];
            while let Some(node) = l {
                res.push(node.val);
                l = node.next;
            }
            res
        }
        let (mut stack1, mut stack2, mut sum, mut res) = (to_vec(l1), to_vec(l2), 0, None);
        while !stack1.is_empty() || !stack2.is_empty() || sum != 0 {
            sum += stack1.pop().unwrap_or(0) + stack2.pop().unwrap_or(0);
            let mut node = ListNode::new(sum % 10);
            node.next = res.take();
            res = Some(Box::new(node.clone()));
            sum /= 10;
        }
        res
    }
}
