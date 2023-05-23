package format

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/drpaepper/go-openpomodoro"
)

const (
	DefaultFormat           = "%R⏱  %c%G🍅\n%d\n%t"
	DefaultExclamationPoint = "❗️"
)

type Formatter func(*openpomodoro.State) string

var FormatParts = map[string]Formatter{
	`%r`: timeRemaining(false),
	`%R`: timeRemaining(true),
	`%m`: minutesRemaining(false),
	`%M`: minutesRemaining(true),
	`%z`: glyphRemaining(false),
	`%Z`: glyphRemaining(true),
	`%p`: percentRemaining(false),
	`%P`: percentRemaining(true),
	`%l`: duration,
	`%L`: durationMinutes,

	`%d`: description,
	`%t`: tags,

	`%c`: goalComplete,
	`%g`: goalTotal(false),
	`%G`: goalTotal(true),
}

func Format(s *openpomodoro.State, f string) string {
	//if s.Pomodoro.IsInactive() {
	//	return ""
	//}

	result := f
	for part, replacement := range FormatParts {
		result = strings.Replace(result, part, replacement(s), -1)
	}
	result = strings.TrimSpace(result)
	return result
}

// DurationAsTime returns a duration string.
func DurationAsTime(d time.Duration) string {
	s := round(d.Seconds())
	return fmt.Sprintf("%d:%02d", s/60, s%60)
}

func timeRemaining(exclaim bool) Formatter {
	return func(s *openpomodoro.State) string {
		d := s.Pomodoro.Remaining()

		if s.Pomodoro.IsDone() {
			if exclaim {
				return DefaultExclamationPoint
			} else {
				return "0:00"
			}
		}

		return DurationAsTime(d)
	}
}

func minutesRemaining(exclaim bool) Formatter {
	return func(s *openpomodoro.State) string {
		if s.Pomodoro.IsDone() {
			if exclaim {
				return DefaultExclamationPoint
			} else {
				return "0"
			}
		}
		return defaultString(s.Pomodoro.RemainingMinutes())
	}
}

func percentRemaining(exclaim bool) Formatter {
	return func(s *openpomodoro.State) string {
		if s.Pomodoro.IsDone() {
			if exclaim {
				return DefaultExclamationPoint
			} else {
				return "0"
			}
		}
		return defaultString(s.Pomodoro.RemainingPercentage())
	}
}

func selectGlyph(isBreak bool, p int) string {

	if isBreak {

		if p > 83 {
			return "󰋙"
		} else if p > 67 {
			return "󰫃"
		} else if p > 50 {
			return "󰫄"
		} else if p > 33 {
			return "󰫅"
		} else if p > 17 {
			return "󰫆"
		} else if p > 0 {
			return "󰫇"
		} else {
			return "󰫈"
		}
	} else {
		if p > 88 {
			return "󰄰"
		} else if p > 75 {
			return "󰪞"
		} else if p > 63 {
			return "󰪟"
		} else if p > 50 {
			return "󰪠"
		} else if p > 38 {
			return "󰪡"
		} else if p > 25 {
			return "󰪢"
		} else if p > 13 {
			return "󰪣"
		} else if p > 0 {
			return "󰪤"
		} else {
			return "󰪥"
		}

	}

}

func glyphRemaining(exclaim bool) Formatter {
	return func(s *openpomodoro.State) string {
		if s.Pomodoro.IsDone() {
			if exclaim {
				return "󰚽"
			} else {
				if s.Pomodoro.Description == "BREAK" {
					return "󰫈"
				} else {
					return "󰪥"
				}
			}
		}
		if s.Pomodoro.Description == "BREAK" {
			return selectGlyph(true, s.Pomodoro.RemainingPercentage())
		} else {
			return selectGlyph(false, s.Pomodoro.RemainingPercentage())
		}
	}
}

func duration(s *openpomodoro.State) string {
	return DurationAsTime(s.Pomodoro.Duration)
}

func durationMinutes(s *openpomodoro.State) string {
	return defaultString(s.Pomodoro.DurationMinutes())
}

func description(s *openpomodoro.State) string {
	return s.Pomodoro.Description
}

func tags(s *openpomodoro.State) string {
	return strings.Join(s.Pomodoro.Tags, ", ")
}

func goalComplete(s *openpomodoro.State) string {
	if s.History == nil {
		return "0"
	}
	return fmt.Sprint(s.History.Date(time.Now()).Count())
}

func goalTotal(slash bool) Formatter {
	return func(s *openpomodoro.State) string {
		if s.Settings == nil || s.Settings.DailyGoal == 0 {
			return ""
		}

		result := fmt.Sprint(s.Settings.DailyGoal)
		if slash {
			result = fmt.Sprintf("/%s", result)
		}
		return result
	}
}

func defaultString(i interface{}) string {
	return fmt.Sprintf("%#v", i)
}

func round(f float64) int {
	return int(math.Floor(f + .5))
}
