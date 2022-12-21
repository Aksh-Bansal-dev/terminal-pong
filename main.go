package main

import (
	"os"
	"os/exec"
	"time"
)

const (
	P1 = iota
	P2
	LEFT_MOVE
	RIGHT_MOVE
)

func main() {
	var (
		n = 25
		m = 100
	)
	playarea := newPlayarea(n, m)

	// taking user input
	ch := make(chan string)
	go func() {
		// disable input buffering
		exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
		// do not display entered characters on the screen
		exec.Command("stty", "-F", "/dev/tty", "-echo").Run()
		var b []byte = make([]byte, 1)
		debounceTimer := time.Now()
		for {
			os.Stdin.Read(b)
			if debounceTimer.Sub(time.Now()).Milliseconds() <= 0 {
				ch <- string(b)
				debounceTimer = time.Now().Add(time.Millisecond * 50)
			}
		}
	}()

	for {
		reading := true
		for reading {
			select {
			case stdin, _ := <-ch:
				if stdin == "a" {
					playarea.move(P2, LEFT_MOVE)
				} else if stdin == "d" {
					playarea.move(P2, RIGHT_MOVE)
				}
			default:
				reading = false
			}
		}
		playarea.draw()
		frameSleep()
	}
}

func frameSleep() {
	time.Sleep(100 * time.Millisecond)
}

func cls() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
