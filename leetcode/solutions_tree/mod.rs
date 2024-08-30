// !1.Rc允许节点有多个所有者，RefCell允许安全地修改这些共享的节点，即使它们在外部看起来是不可变的;
// !2,获取值用 root.borrow().val 即可.
// !3.为什么获取节点需要 let left = root.borrow_mut().left.take();
//    borrow_mut()方法，允许持有RefCell的不可变引用时，仍然能够获取其内部值的可变引用
//    take() 会将Option中的值取出，并留下None，临时取得所有权。和 & 一样。

use std::cell::RefCell;
use std::rc::Rc;

#[derive(Debug, PartialEq, Eq)]
pub struct TreeNode {
    pub val: i32,
    pub left: Option<Rc<RefCell<TreeNode>>>,
    pub right: Option<Rc<RefCell<TreeNode>>>,
}

impl TreeNode {
    #[inline]
    pub fn new(val: i32) -> Self {
        TreeNode {
            val,
            left: None,
            right: None,
        }
    }
}

struct Solution;

impl Solution {
    pub fn max_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        match root {
            Some(root) => {
                let left = root.borrow_mut().left.take();
                let right = root.borrow_mut().right.take();
                1 + Self::max_depth(left).max(Self::max_depth(right))
            }
            None => 0,
        }
    }

    pub fn min_depth(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        match root {
            Some(root) => {
                let left = root.borrow_mut().left.take();
                let right = root.borrow_mut().right.take();
                if left.is_none() {
                    Self::min_depth(right) + 1
                } else if right.is_none() {
                    Self::min_depth(left) + 1
                } else {
                    Self::min_depth(left).min(Self::min_depth(right)) + 1
                }
            }
            None => 0,
        }
    }

    // bfs.
    pub fn min_depth_2(root: Option<Rc<RefCell<TreeNode>>>) -> i32 {
        let mut queue = std::collections::VecDeque::new();
        queue.push_back((1, root));
        while !queue.is_empty() {
            if let Some((depth, Some(node))) = queue.pop_front() {
                match (&node.borrow().left, &node.borrow().right) {
                    (None, None) => {
                        return depth;
                    }
                    (l, r) => {
                        queue.push_back((depth + 1, l.clone()));
                        queue.push_back((depth + 1, r.clone()));
                    }
                }
            }
        }
        0
    }

    pub fn inorder_traversal(root: Option<Rc<RefCell<TreeNode>>>) -> Vec<i32> {
        fn dfs(node: Option<Rc<RefCell<TreeNode>>>, res: &mut Vec<i32>) {
            if let Some(node) = node {
                dfs(node.borrow_mut().left.take(), res);
                res.push(node.borrow().val);
                dfs(node.borrow_mut().right.take(), res);
            }
        }
        let mut res = vec![];
        dfs(root, &mut res);
        res
    }

    pub fn invert_tree(root: Option<Rc<RefCell<TreeNode>>>) -> Option<Rc<RefCell<TreeNode>>> {
        match root {
            Some(root) => {
                let left = Self::invert_tree(root.borrow_mut().left.take());
                let right = Self::invert_tree(root.borrow_mut().right.take());
                root.borrow_mut().left = right;
                root.borrow_mut().right = left;
                Some(root)
            }
            None => None,
        }
    }
}
