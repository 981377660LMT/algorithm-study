"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const isValid = (str) => {
    const stack = [];
    const map = new Map();
    map.set('(', ')').set('{', '}').set('[', ']');
    for (const s of str) {
        if (map.has(s)) {
            stack.push(s);
        }
        else {
            // 判断栈顶元素stack.slice(-1)[0]与当前右括号的关系，匹配则弹出左括号
            const last = stack.slice(-1)[0];
            if (map.get(last) === s) {
                stack.pop();
            }
            else {
                // 不匹配则返回false
                return false;
            }
        }
    }
    return stack.length === 0;
};
console.log(isValid('()'));
console.log(isValid('()('));
console.log(isValid('[]{})'));
console.log(isValid('[{]}'));
console.log(isValid(''));
