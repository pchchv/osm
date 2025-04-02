package mputil

func compact(ms MultiSegment) MultiSegment {
	var at int
	for _, s := range ms {
		if len(s.Line) > 1 {
			ms[at] = s
			at++
		}
	}

	return ms[:at]
}
