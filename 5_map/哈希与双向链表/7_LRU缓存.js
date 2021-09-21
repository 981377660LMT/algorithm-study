function ListNode(key, val) {
  this.key = key
  this.val = val
  this.pre = this.next = null
}

var LRUCache = function (capacity) {
  this.capacity = capacity
  this.size = 0
  this.data = {}
  this.head = new ListNode()
  this.tail = new ListNode()
  this.head.next = this.tail
  this.tail.pre = this.head
}

function get(key) {
  if (this.data[key] !== undefined) {
    let node = this.data[key]
    this.removeNode(node)
    this.appendHead(node)
    return node.val
  } else {
    return -1
  }
}

function put(key, value) {
  let node
  if (this.data[key] !== undefined) {
    node = this.data[key]
    this.removeNode(node)
    node.val = value
  } else {
    node = new ListNode(key, value)
    this.data[key] = node
    if (this.size < this.capacity) {
      this.size++
    } else {
      key = this.removeTail()
      delete this.data[key]
    }
  }
  this.appendHead(node)
}

function removeNode(node) {
  let preNode = node.pre,
    nextNode = node.next
  preNode.next = nextNode
  nextNode.pre = preNode
}

function appendHead(node) {
  let firstNode = this.head.next
  this.head.next = node
  node.pre = this.head
  node.next = firstNode
  firstNode.pre = node
}

function removeTail() {
  let key = this.tail.pre.key
  this.removeNode(this.tail.pre)
  return key
}

const lru = LRUCache(Infinity)
put(1, 1)
console.log(get(1))
