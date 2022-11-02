package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type Event struct {
	time   float64
	kind   uint8
	weight int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, A int
	fmt.Fscan(in, &n, &A)
	fish := make([][3]int, n) // weight, pos, speed
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &fish[i][0], &fish[i][1], &fish[i][2])
	}

	res := 0
	for i := 0; i < n; i++ {
		pos1, speed1 := fish[i][1], fish[i][2]
		events := make([]Event, 0, n)
		for j := 0; j < n; j++ {
			weight2, pos2, speed2 := fish[j][0], fish[j][1], fish[j][2]
			posDiff, speedDiff := pos2-pos1, speed2-speed1
			if speedDiff == 0 {
				if 0 <= posDiff && posDiff <= A {
					events = append(events, Event{0, 0, weight2})
				}
			} else if speedDiff < 0 {
				left, right := float64(A-posDiff)/float64(speedDiff), float64(-posDiff)/float64(speedDiff)
				events = append(events, Event{left, 0, weight2})
				events = append(events, Event{right, 1, -weight2})
			} else {
				left, right := float64(-posDiff)/float64(speedDiff), float64(A-posDiff)/float64(speedDiff)
				events = append(events, Event{left, 0, weight2})
				events = append(events, Event{right, 1, -weight2})
			}
		}

		sort.Slice(events, func(i, j int) bool {
			if events[i].time == events[j].time {
				return events[i].kind < events[j].kind
			}
			return events[i].time < events[j].time
		})

		curSum := 0
		for _, e := range events {
			curSum += e.weight
			if e.time >= 0 {
				res = max(res, curSum)
			}
		}
	}

	fmt.Fprintln(out, res)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
