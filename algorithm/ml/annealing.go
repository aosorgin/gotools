/*
Author:    Alexey Osorgin (alexey.osorgin@gmail.com)
Copyright: Alexey Osorgin, 2017

Brief:     Simulated annealing algorithm primitives
*/

package ml

import (
	"math"
)

type Solution interface {
	GetChanged() (Solution, error)
	CalculateEnergy() (float64, error)
}

type Temperature interface {
	Current() float64
	// if Next() returns false the finging of solutions is completing
	Next() bool
}

func SimulateAnnealing(s Solution, t Temperature, stepIterations int, probabilityToChange float64) (Solution, error) {
	best := s
	bestEnergy, err := s.CalculateEnergy()
	if err != nil {
		return nil, err
	}
	current := s
	currentEnergy := bestEnergy

	for t.Next() {
		for i := 0; i < stepIterations; i++ {
			working, err := current.GetChanged()
			if err != nil {
				return nil, err
			}

			energy, err := working.CalculateEnergy()
			useNew := (energy <= currentEnergy) || (math.Exp((currentEnergy-energy)/t.Current()) > probabilityToChange)

			if useNew == true {
				if energy < bestEnergy {
					best = working
					bestEnergy = energy
				}
				current = working
				currentEnergy = energy
			}
		}
	}
	return best, nil
}

/* Temperature interfaces */

type geometryDecreasingTemperature struct {
	min        float64
	current    float64
	multiplier float64
}

func (t *geometryDecreasingTemperature) Current() float64 {
	return t.current
}

func (t *geometryDecreasingTemperature) Next() bool {
	t.current *= t.multiplier
	return t.current >= t.min
}

func GetGeometryDecreasingTemperature(min float64, max float64, multiplier float64) Temperature {
	return &geometryDecreasingTemperature{
		min:        min,
		current:    max / multiplier,
		multiplier: multiplier,
	}
}
