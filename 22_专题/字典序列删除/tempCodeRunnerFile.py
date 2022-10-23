复次数保证
                while (
                    stack
                    and stack[-1] > char
                    and len(stack) + len(s) - i > k
                    and (stack[-1] != letter or need < remain)
                ):
                    top = stack.pop()
                    if top == letter:
                        need += 1