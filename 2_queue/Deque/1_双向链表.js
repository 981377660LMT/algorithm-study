function ListNode(val) {
  this.val = val
  this.next = null
  this.prev = null
}

/**
 * Initialize your data structure here. Set the size of the deque to be k.
 * @param {number} k
 */
var MyCircularDeque = function (k) {
  this.head = new ListNode('head')
  this.last = new ListNode('last')
  this.head.next = this.last
  this.last.prev = this.head
  this.length = 0
  this.max = k
}

/**
 * Adds an item at the front of Deque. Return true if the operation is successful.
 * @param {number} value
 * @return {boolean}
 */
MyCircularDeque.prototype.insertFront = function (value) {
  if (this.isFull()) {
    return false
  }
  let now = new ListNode(value)
  let next = this.head.next
  now.next = next
  now.prev = this.head
  next.prev = now
  this.head.next = now
  this.length++

  return true
}

/**
 * Adds an item at the rear of Deque. Return true if the operation is successful.
 * @param {number} value
 * @return {boolean}
 */
MyCircularDeque.prototype.insertLast = function (value) {
  if (this.isFull()) {
    return false
  }
  let now = new ListNode(value)
  let prev = this.last.prev
  now.prev = prev
  now.next = this.last
  prev.next = now
  this.last.prev = now
  this.length++

  return true
}

/**
 * Deletes an item from the front of Deque. Return true if the operation is successful.
 * @return {boolean}
 */
MyCircularDeque.prototype.deleteFront = function () {
  if (this.isEmpty()) {
    return false
  }
  let now = this.head.next
  let next = now.next
  next.prev = this.head
  this.head.next = next
  this.length--

  return true
}

/**
 * Deletes an item from the rear of Deque. Return true if the operation is successful.
 * @return {boolean}
 */
MyCircularDeque.prototype.deleteLast = function () {
  if (this.isEmpty()) {
    return false
  }
  let now = this.last.prev
  let prev = now.prev
  prev.next = this.last
  this.last.prev = prev
  this.length--

  return true
}

/**
 * Get the front item from the deque.
 * @return {number}
 */
MyCircularDeque.prototype.getFront = function () {
  if (this.isEmpty()) {
    return -1
  }
  return this.head.next.val
}

/**
 * Get the last item from the deque.
 * @return {number}
 */
MyCircularDeque.prototype.getRear = function () {
  if (this.isEmpty()) {
    return -1
  }
  return this.last.prev.val
}

/**
 * Checks whether the circular deque is empty or not.
 * @return {boolean}
 */
MyCircularDeque.prototype.isEmpty = function () {
  return !this.length
}

/**
 * Checks whether the circular deque is full or not.
 * @return {boolean}
 */
MyCircularDeque.prototype.isFull = function () {
  return this.length >= this.max
}

/**
 * Your MyCircularDeque object will be instantiated and called as such:
 * var obj = new MyCircularDeque(k)
 * var param_1 = obj.insertFront(value)
 * var param_2 = obj.insertLast(value)
 * var param_3 = obj.deleteFront()
 * var param_4 = obj.deleteLast()
 * var param_5 = obj.getFront()
 * var param_6 = obj.getRear()
 * var param_7 = obj.isEmpty()
 * var param_8 = obj.isFull()
 */
