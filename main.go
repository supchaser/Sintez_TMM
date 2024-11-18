package main

import (
	"fmt"
	"math"
	"sort"
)

// структура данных для вывода результатов
type Result struct {
	U1h float64
	Z1  float64
	Z2  float64
	Z3  float64
	Z4  float64
	P   int
	G   float64
}

// Для двухрядного планетарного механизма со смешанным зацеплением
func main() {
	U14 := 16.46103896 // заданное передаточное отношениеы
	var K float64
	K = 3     // число саттелитов
	e := 0.05 // максимальное отклонение
	pi := math.Pi
	ha := 1.0

	results := []Result{}

	for Z1 := 18.0; Z1 <= 100; Z1++ {
		for Z2 := 18.0; Z2 <= 100; Z2++ {
			for Z3 := 21.0; Z3 <= 100; Z3++ {
				// вычисляем число зубьев Z4 из условия соосности
				Z4 := Z1 + Z2 + Z3
				// Z1 и Z4 должны быть некратны K
				if (int(Z1)%int(K) == 0) && (int(Z4)%int(K) == 0) {
					continue
				}
				// проверка на допустимость чисел зубьев
				if Z4-Z3 < 8 && Z4 <= 85 {
					continue // условие не выполнено, переходим к следующей итерации
				}

				// вычисление передаточного отношения редуктора
				U1h := 1 + (Z2*Z4)/(Z1*Z3)
				// проверка на точность совпадения вычисленного передаточного отнешния
				// чисел зубьев с заданным
				if math.Abs((U1h-U14)/(U14)) > e {
					continue // условие не выполнено, переходим к следующей итерации
				}

				// условие соседства
				if math.Sin((pi)/(K)) <= (math.Max(Z2, Z3)+2*ha)/(Z1+Z2) {
					continue // условие не выполнено, переходим к следующей итерации
				}

				// условие сборки
				validP := -1
				for p := 0; p <= 3; p++ {
					temp := ((U1h * Z1) / (K)) * (1 + K*float64(p))
					// Проверяем, является ли temp целым числом
					if math.Mod(temp, 1) == 0 {
						validP = p
						break // найдено целое P
					}
				}

				if validP != -1 {
					// габарит
					G := Z1 + Z2 + Z3
					results = append(results, Result{
						U1h: float64(U1h), Z1: Z1, Z2: Z2, Z3: Z3, Z4: Z4, P: validP, G: G,
					}) // кладем полученный результат в массив массивов Result
				}
			}
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].G < results[j].G
	})

	for i, r := range results {
		fmt.Printf("%d| U1h: %.5f, Z1: %.0f, Z2: %.0f, Z3: %.0f, Z4: %.0f, p: %d, G: %.0f  \n",
			i, r.U1h, r.Z1, r.Z2, r.Z3, r.Z4, r.P, r.G)
	}
}
