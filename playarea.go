package main

import "fmt"

type Playarea struct {
	n      int
	m      int
	mat    [][]byte
	p1     int
	p2     int
	ball   [2]int
	batLen int
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
