package utils

import (
	"cryptoObserver/internal/app/model"
	"strconv"
	"strings"
)

func ParseDecimal(s string) (model.Decimal, error) {
	parts := strings.SplitN(s, ".", 2)
	var d model.Decimal
	var fracStr string
	d.IntPart, _ = strconv.ParseInt(parts[0], 10, 64)
	if len(parts) == 2 {
		fracStr = parts[1]
		// Дополнить нулями до 8 знаков
		if len(fracStr) < 8 {
			fracStr += strings.Repeat("0", 8-len(fracStr))
		} else if len(fracStr) > 8 {
			fracStr = fracStr[:8]
		}
		d.FracPart, _ = strconv.ParseInt(fracStr, 10, 64)
	}
	return d, nil
}

func DecimalToString(d model.Decimal) string {
	var fracStr string
	if d.FracPart == 0 {
		fracStr = "00000000"
	} else {
		fracStr = strconv.FormatInt(d.FracPart, 10)
		if len(fracStr) < 8 {
			fracStr = strings.Repeat("0", 8-len(fracStr)) + fracStr
		} else if len(fracStr) > 8 {
			fracStr = fracStr[:8]
		}
	}
	return strconv.FormatInt(d.IntPart, 10) + "." + fracStr
}
