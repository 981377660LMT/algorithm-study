                            nr, nc = r + dr, c + dc
                            if (
                                0 <= nr < ROW
                                and 0 <= nc < COL
                                and grid[nr][nc] == "W"
                                and (nr, nc) not in visited
                            ):
                                visited.add((nr, nc))
                                queue.append((nr, nc))
