/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Solving problem of settings queens on board that do not terrify each other
           Problem is solved with simulated annealing algorithm
*/

package main

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/aosorgin/gotools/algorithm/ml"
)

type queensOnBoard struct {
	queens      []byte
	queensCount *big.Int
}

func initializeSolution(boardSize int) *queensOnBoard {
	s := &queensOnBoard{
		queens:      make([]byte, boardSize),
		queensCount: big.NewInt(int64(boardSize)),
	}

	for i := 0; i < boardSize; i++ {
		s.queens[i] = byte(i)
	}
	return s
}

func (s *queensOnBoard) GetChanged() (ml.Solution, error) {
	var queen1, queen2 int

	x, err := rand.Int(rand.Reader, s.queensCount)
	if err != nil {
		return nil, err
	}
	queen1 = int(x.Int64())

	for {
		y, err := rand.Int(rand.Reader, s.queensCount)
		if err != nil {
			return nil, err
		}
		if x != y {
			queen1 = int(y.Int64())
			break
		}
	}

	res := &queensOnBoard{
		queens:      make([]byte, len(s.queens)),
		queensCount: s.queensCount,
	}

	for i := 0; i < len(s.queens); i++ {
		res.queens[i] = s.queens[i]
	}

	res.queens[queen1] = s.queens[queen2]
	res.queens[queen2] = s.queens[queen1]
	return res, nil
}

var (
	xDelta = [4]int{-1, -1, 1, 1}
	yDelta = [4]int{-1, 1, -1, 1}
)

func (s *queensOnBoard) CalculateEnergy() (float64, error) {
	matches := 0
	for i := 0; i < len(s.queens); i++ {
		for k := 0; k < 4; k++ {
			tempx := i
			tempy := int(s.queens[i])
			xDiff := xDelta[k]
			yDiff := yDelta[k]

			for {
				tempx += xDiff
				tempy += yDiff
				if tempx < 0 || tempx >= len(s.queens) || tempy < 0 || tempy >= len(s.queens) {
					break
				}
				if s.queens[tempx] == byte(tempy) {
					matches++
				}
			}
		}
	}
	return float64(matches), nil
}

func (s *queensOnBoard) Print() {
	for i := 0; i < len(s.queens); i++ {
		for j := 0; j < len(s.queens); j++ {
			if int(s.queens[i]) == j {
				fmt.Print("Q")
			} else {
				fmt.Print(".")
			}
		}
		fmt.Print("\n")
	}
	fmt.Print("\n")
}

func main() {
	s, err := ml.SimulateAnnealing(initializeSolution(40), ml.GetGeometryDecreasingTemperature(float64(0.5), float64(30), float64(.98)), 100, float64(0.7))
	if err != nil {
		return
	}
	res := s.(*queensOnBoard)
	res.Print()
}
