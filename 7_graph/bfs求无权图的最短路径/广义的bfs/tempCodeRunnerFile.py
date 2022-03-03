 clean(b[:i] + h[j] + b[i:])
                        nextH = h[:j] + h[j + 1 :]
                        if (nextB, nextH) not in visited:
                            visited.add((nextB, nextH))
                            queue.app