package indicators

import (
	"strconv"

	stocks "./stocks"
)

//Go grab data and format it into array
func GetData(code string) []float64 {
	array := make([]float64, 0)
	//set api key
	stocks.SetStockKey("JNHHI33V09S504LS")
	d, _ := stocks.GetDailyStockData("av", code)
	var f int
	//make array
	for i := len(d) - 1; i >= 0; i-- {
		if i != 0 {
			g, _ := strconv.ParseFloat(d[i][5], 64)
			array = append(array, g)
		}
		f++
	}
	return array
}

//Calculate RSI
func RSI(array []float64) float64 {
	//local global values
	var perNum int
	var loss float64
	var gain float64
	var rs float64
	var rsi float64
	for number, value := range array {
		// first RSI calculations
		if number <= 15 {
			change := value - array[perNum]
			if change >= 0 {
				gain = (gain + change) / 2
			} else {
				loss = (loss - change) / 2
			}
			perNum = number
			rs = gain / loss
			rsi = 100 - (100 / (1 + rs))
		}
		//RSI after first calculations
		if number >= 15 {
			change := value - array[number-1]
			gain = gain * 13
			loss = loss * 13
			if change >= 0 {
				gain = (gain + change) / 14
				loss = loss / 14
			} else {
				loss = (loss + change) / 14
				gain = gain / 14
			}
			rs = gain / (-1 * loss)
			rsi = 100 - (100 / (1 + rs))
		}

	}
	return rsi
}

//EMA calculations
func EMA(smoothing float64, days float64, perEMA float64, cValue float64) float64 {
	return (cValue * (smoothing / (1 + days))) + perEMA*(1-(smoothing/(1+days)))
}

//EMA with array
func SEMA(array []float64, days float64) float64 {
	var ema float64
	for num, y := range array {
		if num > int(days-1) {
			ema = EMA(2, days, ema, y)
		} else if num == int(days-1) {
			ema = SMA(0, int(days-1), array)
		}
	}
	return ema
}

//SMA calculations
func SMA(startday int, endday int, array []float64) float64 {
	var f int
	var num float64
	for f <= endday {
		if f >= startday {
			num = num + array[f]
		}
		f++
	}

	return (num / (float64(endday-startday) + 1))
}

//Calculate MACD indicator
func MACD(array []float64) float64 {
	var per12EMA float64
	var per26EMA float64
	var macd float64
	for num, value := range array {
		if num >= 26 {
			if num == 26 {
				per12EMA = EMA(2, 12, SMA(14, 26, array), value)
				per26EMA = EMA(2, 26, SMA(0, 26, array), value)
				macd = per12EMA - per26EMA
			} else {
				per12EMA = EMA(2, 12, per12EMA, value)
				per26EMA = EMA(2, 26, per26EMA, value)
				macd = per12EMA - per26EMA
			}
		}
	}
	return macd
}

//Calculate Aroon
func AROON(array []float64) float64 {
	//return value
	var aroon float64
	//counter
	var f int
	//smallest day number
	var sDay int
	//Largest day number
	var lDay int
	//Smallest number
	var smallesNumber float64
	//Largest number memory
	var largestNumber float64
	//Range over array
	for num, _ := range array {
		//make sure it has 25 days pior
		if num > 24 {
			// loop through a range of numbers
			for num-25 <= f && f <= num {
				//if number is larger rember it and the day
				if array[f] >= largestNumber {
					largestNumber = array[f]
					lDay = (num - f)
				}
				//if number is smaller rember it and the day
				if array[f] <= smallesNumber {
					smallesNumber = array[f]
					sDay = (num - f)
				}
				f++
			}
			//aroon up
			aroonUp := 100 * ((25 - float64(lDay)) / 25)
			//aroon down
			aroonDown := 100 * ((25 - float64(sDay)) / 25)
			//Calculate Aroon
			aroon = aroonUp - aroonDown
			//reset
			smallesNumber = 10000000000000000000000000000000
			largestNumber = 0
			f = (num + 1) - 25
		}
	}
	return aroon
}
func PPO(array []float64) float64 {
	day12SMA := SMA(14, 25, array)
	day26SMA := SMA(0, 25, array)
	day12EMA := EMA(2, 12, day12SMA, array[26])
	day26EMA := EMA(2, 26, day26SMA, array[26])
	var app float64
	f := 27
	for f <= (len(array) - 1) {
		day12EMA = EMA(2, 12, day12EMA, array[f])
		day26EMA = EMA(2, 26, day26EMA, array[f])
		app = ((day12EMA - day26EMA) / day26EMA) * 100
		f++
	}
	return app
}
