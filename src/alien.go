package main

type Alien struct {
	i    int
	Loc  *City
	dead bool
}

func NewAlien(loc *City, i int) *Alien {
	alien := &Alien{Loc: loc, i: i}
	if loc != nil {
		loc.Aliens = append(loc.Aliens, alien)
	}
	return alien
}

func (a *Alien) move(dir string) bool {
	if a.dead {
		return false
	}
	var next *City
	var set bool

	switch dir {
	case east:
		if a.Loc != nil {
			next = a.Loc.Directions[east]
		}
		set = true
	case west:
		if a.Loc != nil {
			next = a.Loc.Directions[west]
		}
		set = true
	case north:
		if a.Loc != nil {
			next = a.Loc.Directions[north]
		}
		set = true
	case south:
		if a.Loc != nil {
			next = a.Loc.Directions[south]
		}
		set = true
	}

	if set {
		if a.Loc != nil {
			// this can be set as nil
			// because before moving
			// the user needs to check if the city is destroyed, ie has 2 aliens
			// this code will execute if there is only 1 alien
			a.Loc.Aliens = nil
		}
		if next != nil {
			if next.isDestroyed() {
				return false
			}
			a.Loc = next
		}
		if a.Loc != nil {
			a.Loc.Aliens = append(a.Loc.Aliens, a)
			return true
		}
	}
	return false
}
