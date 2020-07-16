// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 156.

// Package geometry defines simple types for plane geometry.
//!+point
package main

import (
	"fmt"
	"math"
)

type Point struct {
	X, Y float64
}

// traditional function
func Distance(p, q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// same thing, but as a method of the Point type
func (p Point) Distance(q Point) float64 {
	return math.Hypot(q.X-p.X, q.Y-p.Y)
}

// if the receiver need to be update, pass a pointer
func (p *Point) ScaleBy(factor float64) {
	p.X *= factor
	p.Y *= factor
}

func (p Point) foo(factor float64) {
	p.X *= factor
	p.Y *= factor
}

//!-point

//!+path

// A Path is a journey connecting the points with straight lines.
type Path []Point

// Distance returns the distance traveled along the path.
func (path Path) Distance() float64 {
	sum := 0.0
	for i := range path {
		if i > 0 {
			sum += path[i-1].Distance(path[i])
		}
	}
	return sum
}

func main() {
	p := Point{
		X: 1,
		Y: 2,
	}
	fmt.Printf("%+v", p)

	p.ScaleBy(3)
	fmt.Printf("%+v", p)

	(&p).ScaleBy(4)
	fmt.Printf("%+v", p)

	(&p).foo(4)
	fmt.Printf("%+v", p)

}

//!-path
