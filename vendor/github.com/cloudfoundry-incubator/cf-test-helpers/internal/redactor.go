package internal

import "strings"

type Redactor interface {
	Redact(toRedact string) string
}

type redactor struct {
	redactees []string
}

func NewRedactor(redactees ...string) Redactor {
	return &redactor{
		redactees: redactees,
	}
}

func (r *redactor) Redact(toRedact string) string {
	if len(r.redactees) == 0 {
		return toRedact
	}

	var out []string
	for _, candidate := range strings.Fields(toRedact) {
		if r.shouldBeRedacted(candidate) {
			out = append(out, "[REDACTED]")
		} else {
			out = append(out, candidate)
		}
	}

	return strings.Join(out, " ")
}

func (r *redactor) shouldBeRedacted(val string) bool {
	for _, v := range r.redactees {
		if v == val {
			return true
		}
	}

	return false
}
