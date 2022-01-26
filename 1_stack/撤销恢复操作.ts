type InputWord = 'undo' | 'redo' | (string & {})

const input = 'hello undo redo world.'
const words: InputWord[] = input.split(' ')
const stack: InputWord[] = []
let undo: Exclude<InputWord, 'undo' | 'redo'>[] = []

for (const word of words) {
  if (word === 'undo') {
    undo.push(stack.pop()!)
  } else if (word === 'redo') {
    stack.push(undo.pop()!)
  } else {
    stack.push(word)
    undo = []
  }
}

console.log(stack.join(' '))
export {}
