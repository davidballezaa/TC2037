package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// Definición de las matrices
type Matriz struct {
	data [][]int
	mu   sync.Mutex
}

// Función para inicializar una matriz de tamaño p x q
func NuevaMatriz(p, q int) *Matriz {
	data := make([][]int, p)
	for i := range data {
		data[i] = make([]int, q)
	}
	return &Matriz{data: data}
}

// Función para imprimir una matriz
func (m *Matriz) Imprimir() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for _, fila := range m.data {
		for _, val := range fila {
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
}

// Función para multiplicar dos matrices concurrentemente
func MultiplicarElemento(A, B, C *Matriz, i, j, q int, c chan bool) {
	C.data[i][j] = 0
	for k := 0; k < q; k++ {
		C.mu.Lock()
		C.data[i][j] += A.data[i][k] * B.data[k][j]
		C.mu.Unlock()
	}
	c <- true
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Semilla para generar números aleatorios

	// Dimensiones de las matrices
	p, q, r := 3, 3, 3

	// Inicialización de matrices
	A := NuevaMatriz(p, q)
	B := NuevaMatriz(q, r)
	C := NuevaMatriz(p, r)

	// Rellenar matrices A y B con valores de ejemplo
	A.data = [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7, 8, 9},
	}

	B.data = [][]int{
		{9, 8, 7},
		{6, 5, 4},
		{3, 2, 1},
	}

	// Mostrar matrices A y B
	fmt.Println("Matriz A:")
	A.Imprimir()

	fmt.Println("Matriz B:")
	B.Imprimir()

	// Canal para sincronizar la finalización de los cálculos
	c := make(chan bool)

	// Lanzar goroutines para cada elemento de la matriz resultante
	for i := 0; i < p; i++ {
		for j := 0; j < r; j++ {
			go MultiplicarElemento(A, B, C, i, j, q, c)
		}
	}

	// Esperar a que todas las goroutines terminen
	for i := 0; i < p*r; i++ {
		<-c
		fmt.Println("¡Un elemento ha sido calculado!")
	}

	// Mostrar la matriz resultante C
	fmt.Println("Matriz C (resultado de A x B):")
	C.Imprimir()
}
