package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Playarea struct {
	n      int
	m      int
	mat    [][]byte
	p1     int
	p2     int
	ball   [2]int
	batLen int
}

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
		prevChar := "$"
		for {
			os.Stdin.Read(b)
			if debounceTimer.Sub(time.Now()).Milliseconds() <= 0 {
				input := string(b)
				ch <- input
				if input == prevChar {
					debounceTimer = time.Now().Add(time.Millisecond * 50)
				}
				prevChar = input
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
				} else if stdin == "h" {
					playarea.move(P1, LEFT_MOVE)
				} else if stdin == "l" {
					playarea.move(P1, RIGHT_MOVE)
				}
			default:
				reading = false
			}
		}
		playarea.draw()
		frameSleep()
	}
}

func (p *Playarea) draw() {
	mat := p.mat
	cls()
	for i := 0; i < len(mat); i++ {
		fmt.Printf("%s\n", mat[i])
	}
}

func (p *Playarea) move(player, dir int) {
	row := 1
	if player == P2 {
		row = p.n
	}
	for i := 1; i < p.m; i++ {
		p.mat[row][i] = ' '
	}

	if dir == LEFT_MOVE {
		if player == P1 {
			if p.p1 > 1 {
				p.p1--
			}
		} else {
			if p.p2 > 1 {
				p.p2--
			}
		}
	} else {
		if player == P1 {
			if p.p1 < p.m-p.batLen {
				p.p1++
			}
		} else {
			if p.p2 < p.m-p.batLen {
				p.p2++
			}
		}
	}

	for i := 0; i < p.batLen; i++ {
		if player == P1 {
			p.mat[row][p.p1+i] = '='
		} else {
			p.mat[row][p.p2+i] = '='
		}
	}
}

func newPlayarea(n, m int) Playarea {
	mat := [][]byte{}
	for i := -1; i <= n; i++ {
		temp := []byte{}
		for j := -1; j <= m; j++ {
			if i == -1 || i == n {
				temp = append(temp, '-')
			} else if j == -1 || j == m {
				temp = append(temp, '|')
			} else {
				temp = append(temp, ' ')
			}
		}
		mat = append(mat, temp)
	}
	for i := m/2 - 4; i < m/2+4; i++ {
		mat[n][i] = '='
		mat[1][i] = '='
	}
	mat[n/2][m/2] = 'o'
	return Playarea{n: n, m: m, mat: mat, p1: m/2 - 4, p2: m/2 - 4, ball: [2]int{n / 2, m / 2}, batLen: 8}
}

func frameSleep() {
	time.Sleep(100 * time.Millisecond)
}

func cls() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
