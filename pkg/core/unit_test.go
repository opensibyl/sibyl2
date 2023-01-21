package core

import "testing"

func TestSpan_HasInteraction(t *testing.T) {
	t.Parallel()
	a := &Span{
		Start: Point{
			1,
			2,
		},
		End: Point{
			1,
			3,
		},
	}

	b := &Span{
		Start: Point{
			2,
			2,
		},
		End: Point{
			2,
			3,
		},
	}

	if a.HasInteraction(b) {
		panic(nil)
	}

	c := &Span{
		Start: Point{
			1,
			2,
		},
		End: Point{
			5,
			10,
		},
	}

	d := &Span{
		Start: Point{
			5,
			11,
		},
		End: Point{
			10,
			13,
		},
	}

	if c.HasInteraction(d) {
		panic(nil)
	}

	e := &Span{
		Start: Point{
			1,
			2,
		},
		End: Point{
			5,
			10,
		},
	}

	f := &Span{
		Start: Point{
			5,
			8,
		},
		End: Point{
			10,
			13,
		},
	}

	if !e.HasInteraction(f) {
		panic(nil)
	}
}
