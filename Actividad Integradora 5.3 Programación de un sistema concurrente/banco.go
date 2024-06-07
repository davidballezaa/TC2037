package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

// Definición de la estructura Cuenta con un saldo y un mutex para sincronización
type Cuenta struct {
	saldo int
	mu    sync.Mutex
}

// Función para depositar dinero en una cuenta
func (c *Cuenta) Depositar(monto int) {
	c.mu.Lock()         // Inicio de la región crítica
	c.saldo += monto
	c.mu.Unlock()       // Fin de la región crítica
}

// Función para retirar dinero de una cuenta
func (c *Cuenta) Retirar(monto int) bool {
	c.mu.Lock()         // Inicio de la región crítica
	defer c.mu.Unlock() // Asegurar el desbloqueo al finalizar la función

	if c.saldo >= monto {
		c.saldo -= monto
		return true
	}
	return false
}

// Función para transferir dinero entre dos cuentas
func Transferir(cuentaOrigen, cuentaDestino *Cuenta, monto int) bool {
	if cuentaOrigen == cuentaDestino {
		return false
	}

	/*
	Supongamos que tenemos dos cuentas, cuentaA y cuentaB, y dos hilos intentan transferir dinero simultáneamente entre estas cuentas:

	Hilo 1 quiere transferir de cuentaA a cuentaB.
	Hilo 2 quiere transferir de cuentaB a cuentaA.

	Sin seguir un orden de bloqueo, pasaría esto:
	Hilo 1 bloquea cuentaA y espera bloquear cuentaB.
	Hilo 2 bloquea cuentaB y espera bloquear cuentaA.
	deadlock...

	Decidimos usar las direcciones de memoria como parámetro para seguir un mismo orden de bloqueo (Línea 56).
	*/
	
	if uintptr(unsafe.Pointer(cuentaOrigen)) < uintptr(unsafe.Pointer(cuentaDestino)) {
		cuentaOrigen.mu.Lock()
		defer cuentaOrigen.mu.Unlock()
		cuentaDestino.mu.Lock()
		defer cuentaDestino.mu.Unlock()
	} else {
		cuentaDestino.mu.Lock()
		defer cuentaDestino.mu.Unlock()
		cuentaOrigen.mu.Lock()
		defer cuentaOrigen.mu.Unlock()
	}

	if cuentaOrigen.saldo >= monto {
		cuentaOrigen.saldo -= monto
		cuentaDestino.saldo += monto
		return true
	}
	return false
}

// Función que simula las transacciones de un cliente
func cliente(id int, cuentas []*Cuenta, c chan bool) {
	for i := 0; i < 10; i++ { // Realizar 10 transacciones por cliente
		op := rand.Intn(3) // Elegir una operación aleatoria: 0=Depositar, 1=Retirar, 2=Transferir
		cuenta := cuentas[rand.Intn(len(cuentas))]

		switch op {
		case 0: // Depositar
			monto := rand.Intn(100)
			cuenta.Depositar(monto)
			fmt.Printf("Cliente %d depositó %d en cuenta %p. Nuevo saldo: %d\n", id, monto, cuenta, cuenta.saldo)
		case 1: // Retirar
			monto := rand.Intn(100)
			if cuenta.Retirar(monto) {
				fmt.Printf("Cliente %d retiró %d de cuenta %p. Nuevo saldo: %d\n", id, monto, cuenta, cuenta.saldo)
			} else {
				fmt.Printf("Cliente %d intentó retirar %d de cuenta %p, pero no tenía suficientes fondos.\n", id, monto, cuenta)
			}
		case 2: // Transferir
			cuentaDestino := cuentas[rand.Intn(len(cuentas))]
			monto := rand.Intn(100)
			if Transferir(cuenta, cuentaDestino, monto) {
				fmt.Printf("Cliente %d transfirió %d de cuenta %p a cuenta %p. Nuevo saldo cuenta origen: %d, nuevo saldo cuenta destino: %d\n", id, monto, cuenta, cuentaDestino, cuenta.saldo, cuentaDestino.saldo)
			} else {
				fmt.Printf("Cliente %d intentó transferir %d de cuenta %p a cuenta %p, pero no tenía suficientes fondos.\n", id, monto, cuenta, cuentaDestino)
			}
		}
		time.Sleep(time.Millisecond * 100) // Esperar un poco entre transacciones para simular actividad real
	}
	c <- true
}

func main() {
	rand.Seed(time.Now().UnixNano()) // Semilla para generar números aleatorios

	// Crear cuentas bancarias
	cuentas := make([]*Cuenta, 5)
	for i := 0; i < 5; i++ {
		cuentas[i] = &Cuenta{saldo: 1000}
	}

	// Canal para sincronizar la finalización de los clientes
	c := make(chan bool)

	// Lanzar clientes
	for i := 0; i < 10; i++ {
		go cliente(i, cuentas, c)
	}

	// Esperar a que todos los clientes terminen
	for i := 0; i < 10; i++ {
		<-c
		fmt.Println("¡Un cliente ha terminado su trabajo!")
	}

	// Mostrar los saldos finales de las cuentas
	fmt.Println("Transacciones completadas. Saldos finales:")
	for i, cuenta := range cuentas {
		fmt.Printf("Cuenta %d: %d\n", i, cuenta.saldo)
	}
}
