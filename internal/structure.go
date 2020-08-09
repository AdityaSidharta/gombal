package internal

import "fmt"

type Tuple struct {
    Q, R string
}

func Zip(as []string, bs []string) ([]Tuple, error) {
	n := len(as)

    if len(as) != len(bs) {
        return nil, fmt.Errorf("zip: arguments must be of same length")
    }

    tuples := make([]Tuple, n, n)

    for i, _ := range as {
    	q := as[i]
    	r := bs[i]
        tuples[i] = Tuple{q, r}
    }

    return tuples, nil
}
