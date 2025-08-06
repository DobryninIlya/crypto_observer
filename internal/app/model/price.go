package model

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

type Decimal struct {
	IntPart  int64 `json:"int"`  // Целая часть
	FracPart int64 `json:"frac"` // Дробная часть (всегда 8 знаков)
}

func (d *Decimal) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" {
		return nil
	}

	if strings.ContainsAny(s, "eE") {
		f, _, err := big.ParseFloat(s, 10, 64, big.ToNearestEven)
		if err != nil {
			return fmt.Errorf("failed to parse scientific notation: %w", err)
		}

		s = f.Text('f', 8) // 8 знаков после запятой
	}

	parts := strings.SplitN(s, ".", 2)
	var err error

	d.IntPart, err = strconv.ParseInt(parts[0], 10, 64)
	if err != nil {
		return fmt.Errorf("invalid integer part: %w", err)
	}

	if len(parts) == 2 {
		fracStr := parts[1]

		if len(fracStr) < 8 {
			fracStr += strings.Repeat("0", 8-len(fracStr))
		} else if len(fracStr) > 8 {
			fracStr = fracStr[:8]
		}

		d.FracPart, err = strconv.ParseInt(fracStr, 10, 64)
		if err != nil {
			return fmt.Errorf("invalid fraction part: %w", err)
		}
	} else {
		d.FracPart = 0
	}

	return nil
}

func (d Decimal) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

func (d Decimal) String() string {
	fracStr := fmt.Sprintf("%08d", d.FracPart)
	return fmt.Sprintf("%d.%s", d.IntPart, fracStr)
}

func (d Decimal) IsZero() bool {
	return d.IntPart == 0 && d.FracPart == 0
}
