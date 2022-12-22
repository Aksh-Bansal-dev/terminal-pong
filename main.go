package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Playarea struct {
	n       int
	m       int
	mat     [][]byte
	p1      int
	p2      int
	ball    [2]int
	ballDir [2]int
	batLen  int
	p1Dir   int
	p2Dir   int
}

const (
	frameTime = 50

	P1 = iota
	P2
	LEFT_MOVE
	RIGHT_MOVE
)

func main() {
	config := getConfig("./config.json")
	playarea := newPlayarea(config)

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
					debounceTimer = time.Now().Add(time.Millisecond * frameTime / 2)
				}
				prevChar = input
			}
		}
	}()

	iteration := 0
	for {
		reading := true
		for reading {
			select {
			case stdin, _ := <-ch:
				if stdin == "a" {
					playarea.move(P2, LEFT_MOVE, config.BatSpeed)
				} else if stdin == "f" {
					playarea.move(P2, RIGHT_MOVE, config.BatSpeed)
				} else if stdin == "h" {
					playarea.move(P1, LEFT_MOVE, config.BatSpeed)
				} else if stdin == "l" {
					playarea.move(P1, RIGHT_MOVE, config.BatSpeed)
				}
			default:
				reading = false
			}
		}
		end := 0
		if iteration == 0 {
			end = playarea.moveBall()
		}
		playarea.draw()
		if end > 0 {
			fmt.Printf("\nPlayer%d Won!\n", end)
			os.Exit(0)
		}
		frameSleep()
		iteration = (iteration + 1) % (5 - config.BallSpeed)
	}
}

func (p *Playarea) draw() {
	mat := p.mat
	cls()
	for i := 0; i < len(mat); i++ {
		fmt.Printf("%s\n", mat[i])
	}
}

func (p *Playarea) moveBall() int {
	ball := p.ball
	p.mat[ball[0]][ball[1]] = ' '
	ball[0] += p.ballDir[0]
	ball[1] += p.ballDir[1]
	if ball[0] == 1 && (p.p1 > ball[1] || p.p1+p.batLen < ball[1]) {
		p.mat[ball[0]][ball[1]] = 'o'
		return 2
	} else if ball[0] == p.n && (p.p2 > ball[1] || p.p2+p.batLen < ball[1]) {
		p.mat[ball[0]][ball[1]] = 'o'
		return 1
	} else if ball[0] == p.n && (p.p2 <= ball[1] && p.p2+p.batLen >= ball[1]) {
		p.ballDir = [2]int{-1, p.p2Dir}
	} else if ball[0] == 1 && (p.p1 <= ball[1] && p.p1+p.batLen >= ball[1]) {
		p.ballDir = [2]int{1, p.p1Dir}
	} else {
		p.mat[ball[0]][ball[1]] = 'o'
		p.ball = ball
		if ball[1] == p.m || ball[1] == 1 {
			p.ballDir[1] *= -1
		}
	}
	return 0
}

func (p *Playarea) move(player, dir int, batSpeed int) {
	row := 1
	if player == P2 {
		row = p.n
	}
	for i := 1; i <= p.m; i++ {
		p.mat[row][i] = ' '
	}

	if dir == LEFT_MOVE {
		if player == P1 {
			if p.p1 > 1 {
				p.p1Dir = -1
				p.p1 -= batSpeed
			}
		} else {
			if p.p2 > 1 {
				p.p2Dir = -1
				p.p2 -= batSpeed
			}
		}
	} else {
		if player == P1 {
			if p.p1 < p.m-p.batLen {
				p.p1Dir = 1
				p.p1 += batSpeed
			}
		} else {
			if p.p2 < p.m-p.batLen {
				p.p2Dir = 1
				p.p2 += batSpeed
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

func newPlayarea(config Config) Playarea {
	n := config.Rows
	m := config.Cols
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
	for i := m/2 - config.BatLength/2; i < m/2+config.BatLength/2; i++ {
		mat[n][i] = '='
		mat[1][i] = '='
	}
	mat[n/2][m/2] = 'o'
	return Playarea{
		n:       n,
		m:       m,
		mat:     mat,
		p1:      m/2 - config.BatLength/2,
		p2:      m/2 - config.BatLength/2,
		p1Dir:   0,
		p2Dir:   0,
		ball:    [2]int{n / 2, m / 2},
		ballDir: [2]int{1, 0},
		batLen:  config.BatLength}
}

func frameSleep() {
	time.Sleep(frameTime * time.Millisecond)
}

func cls() {
	cmd := exec.Command("clear") //Linux example, its tested
	cmd.Stdout = os.Stdout
	cmd.Run()
}
