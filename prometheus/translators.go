package prometheus

import "io"

func writeType(name string, t string, w io.Writer) (int, error) {
	var written int
	n, err := w.Write([]byte("# TYPE "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(name))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(" "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(t))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte("\n"))
	written += n
	if err != nil {
		return written, err
	}

	return written, nil
}

//metric_without_timestamp_and_labels 12.47
func writeCounter(name string, c int64, w io.Writer) (int, error) {
	var written int
	n, err := w.Write([]byte(name))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(" "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = writeInt(w, c)
	written += n
	if err != nil {
		return written, err
	}
	/* todo: not sure what this is used for yet, go-metrics does not have it natively so might skip implementation
	 * todo: I assume it does not make much sense to add a timestamp now
	if metric.TimestampMs != nil {
		err = w.WriteByte(' ')
		written++
		if err != nil {
			return written, err
		}
		n, err = writeInt(w, *metric.TimestampMs)
		written += n
		if err != nil {
			return written, err
		}
	}
	*/
	n, err = w.Write([]byte("\n"))
	written += n
	if err != nil {
		return written, err
	}

	return written, nil
}

func writeGauge(name string, g float64, w io.Writer) (int, error) {
	var written int
	n, err := w.Write([]byte(name))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(" "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = writeFloat(w, g)
	written += n
	if err != nil {
		return written, err
	}

	n, err = w.Write([]byte("\n"))
	written += n
	if err != nil {
		return written, err
	}

	return written, nil
}

func writeGaugeQuantile(name string, g float64, q float64, w io.Writer) (int, error) {
	var written int
	n, err := w.Write([]byte(name))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte("{quantile=\""))
	written += n
	if err != nil {
		return written, err
	}
	n, err = writeFloat(w, q)
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte("\"}"))
	written += n
	if err != nil {
		return written, err
	}

	n, err = w.Write([]byte(" "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = writeFloat(w, g)
	written += n
	if err != nil {
		return written, err
	}

	n, err = w.Write([]byte("\n"))
	written += n
	if err != nil {
		return written, err
	}

	return written, nil
}

/*
func writeGaugeLabel(name string, g float64, label string, w io.Writer) (int, error) {
	var written int
	n, err := w.Write([]byte(name))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte("{le=\""))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte(label))
	written += n
	if err != nil {
		return written, err
	}
	n, err = w.Write([]byte("\"}"))
	written += n
	if err != nil {
		return written, err
	}

	n, err = w.Write([]byte(" "))
	written += n
	if err != nil {
		return written, err
	}
	n, err = writeFloat(w, g)
	written += n
	if err != nil {
		return written, err
	}

	n, err = w.Write([]byte("\n"))
	written += n
	if err != nil {
		return written, err
	}

	return written, nil
}
*/
