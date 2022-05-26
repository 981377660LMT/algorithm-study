import { AvlTree } from 'datastructures-js'

const avl = new AvlTree()
avl.insert(2, undefined)
console.log(avl.ceil(1)?.getValue())
console.log(avl.ceil(1)?.getKey())
