package utils

import (
	"fmt"
	"time"
)

var Loc, _ = time.LoadLocation("Europe/Belgrade")

func TimeAgo(t time.Time) string {
	now := time.Now().In(Loc)
	diff := now.Sub(t)

	switch {
	case diff < time.Second:
		return "sada"
	case diff < time.Minute:
		sec := int(diff.Seconds())
		if sec < 5 {
			return fmt.Sprintf("pre %d sekunde", sec)
		}
		return fmt.Sprintf("pre %d sekundi", sec)
	case diff < time.Hour:
		min := int(diff.Minutes())
		if min == 1 {
			return "pre 1 minut"
		}
		return fmt.Sprintf("pre %d minuta", min)
	case diff < 24*time.Hour:
		h := int(diff.Hours())
		if h == 1 {
			return "pre 1 sat"
		} else if h < 5 {
			return fmt.Sprintf("pre %d sata", h)
		}
		return fmt.Sprintf("pre %d sati", h)
	case diff < 30*24*time.Hour:
		d := int(diff.Hours() / 24)
		if d == 1 {
			return "pre 1 dan"
		}
		return fmt.Sprintf("pre %d dana", d)
	case diff < 365*24*time.Hour:
		m := int(diff.Hours() / 24 / 30)
		if m == 1 {
			return "pre 1 mesec"
		}
		return fmt.Sprintf("pre %d meseci", m)
	default:
		y := int(diff.Hours() / 24 / 365)
		if y == 1 {
			return "pre godinu dana"
		}
		return fmt.Sprintf("pre %d godina", y)
	}
}
