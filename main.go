package main

import (
	"fmt"
	"math"
	"os"
	"slices"
	"sync"
	"time"
)

type MechanismConfig struct {
	U     float64 // заданное передаточное отношение
	K     float64 // число саттелитов
	E     float64 // максимальное отклонение
	Ha    float64 // коэффициент высоты головки зубьев, ha*
	FlagK bool    // для того, чтобы установить проверку на кратность K
}

// структура данных для вывода результатов
type Result struct {
	U1h float64 // передаточное отношение
	Z1  float64 // число зубьев Z1, Z2, Z3, Z4
	Z2  float64
	Z3  float64
	Z4  float64
	P   int     // произвольное дополнительное число оборотов водила при сборке
	G   float64 // габарит
}

func calculatingResults(config MechanismConfig) []Result {
	var results []Result
	var mu sync.Mutex
	wg := sync.WaitGroup{}

	for Z1 := 18.0; Z1 <= 100; Z1++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for Z2 := 18.0; Z2 <= 100; Z2++ {
				for Z3 := 21.0; Z3 <= 100; Z3++ {
					// вычисляем число зубьев Z4 из условия соосности
					Z4 := Z1 + Z2 + Z3
					if config.FlagK {
						// Z1 и Z4 должны быть некратны K
						if (int(Z1)%int(config.K) == 0) && (int(Z4)%int(config.K) == 0) {
							continue
						}
					}

					// проверка на допустимость чисел зубьев
					if Z4-Z3 < 8 && Z4 <= 85 {
						continue // условие не выполнено, переходим к следующей итерации
					}

					// вычисление передаточного отношения редуктора
					U1h := 1 + (Z2*Z4)/(Z1*Z3)
					// проверка на точность совпадения вычисленного передаточного отнешния
					// чисел зубьев с заданным
					if math.Abs((U1h-config.U)/(config.U)) > config.E {
						continue // условие не выполнено, переходим к следующей итерации
					}

					// условие соседства
					if math.Sin((math.Pi)/(config.K)) <= (math.Max(Z2, Z3)+2*config.Ha)/(Z1+Z2) {
						continue // условие не выполнено, переходим к следующей итерации
					}

					// условие сборки
					validP := -1
					for p := 0; p <= 3; p++ {
						temp := ((U1h * Z1) / (config.K)) * (1 + config.K*float64(p))
						// проверяем, является ли temp целым числом
						if math.Mod(temp, 1) == 0 {
							validP = p
							break // найдено целое P
						}
					}

					if validP != -1 {
						G := Z1 + Z2 + Z3
						result := Result{
							U1h: float64(U1h), Z1: Z1, Z2: Z2, Z3: Z3, Z4: Z4, P: validP, G: G,
						}

						mu.Lock()
						results = append(results, result)
						mu.Unlock()
					}
				}
			}
		}()
	}

	wg.Wait()
	return results
}

func Sort(results []Result) {
	slices.SortFunc(results, func(a, b Result) int {
		if a.G < b.G {
			return -1
		} else if a.G > b.G {
			return 1
		} else {
			return 0
		}
	})
}

func PrintResults(results []Result) {
	file, err := os.Create("static/results.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка при создании файла: %v\n", err)
		return
	}
	defer file.Close()

	for i, r := range results {
		fmt.Fprintf(file, "%d| U1h: %.5f, Z1: %.0f, Z2: %.0f, Z3: %.0f, Z4: %.0f, p: %d, G: %.0f  \n",
			i, r.U1h, r.Z1, r.Z2, r.Z3, r.Z4, r.P, r.G)
	}
}

func main() {
	start := time.Now()
	config := MechanismConfig{
		U:     16.46103896,
		K:     3,
		E:     0.05,
		Ha:    1.0, // ha*
		FlagK: true,
	}

	results := calculatingResults(config)

	Sort(results)

	PrintResults(results)

	elapsed := time.Since(start)
	fmt.Fprintf(os.Stderr, "Время выполнения программы: %v\n", elapsed)
}
