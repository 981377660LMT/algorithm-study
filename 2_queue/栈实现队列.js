// 使用两个数组的栈方法（push, pop） 实现队列
/**
 * Initialize your data structure here.
 */
var MyQueue = function () {
  this.stack1 = []
  this.stack2 = []
}

/**
 * Push element x to the back of queue.
 * @param {number} x
 * @return {void}
 */
MyQueue.prototype.push = function (x) {
  this.stack1.push(x)
}

/**
 * Removes the element from in front of queue and returns that element.
 * @return {number}
 */
MyQueue.prototype.pop = function () {
  const size = this.stack2.length
  if (size) {
    return this.stack2.pop()
  }
  while (this.stack1.length) {
    this.stack2.push(this.stack1.pop())
  }
  return this.stack2.pop()
}

/**
 * Get the front element.
 * @return {number}
 */
MyQueue.prototype.peek = function () {
  const x = this.pop()
  this.stack2.push(x)
  return x
}

/**
 * Returns whether the queue is empty.
 * @return {boolean}
 */
MyQueue.prototype.empty = function () {
  return !this.stack1.length && !this.stack2.length
}
