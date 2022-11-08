package arknights

type storage uint64

func (s *storage) is6starsmode() bool {
	return *s&1 == 1
}

func (s *storage) setmode(is6stars bool) bool {
	if is6stars {
		*s |= 1
	} else {
		*s &= 0xffffffff_fffffffe
	}
	return is6stars
}
