package banco

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	"unsafe"
)

// Definición de la estructura Cuenta con un saldo y un mutex para sincronización
type Cuenta struct {
	Saldo int
	Mu    sync.Mutex
}

// Función para depositar dinero en una cuenta
func (c *Cuenta) Depositar(monto int) {
	c.Mu.Lock()         // Inicio de la región crítica
	defer c.Mu.Unlock() // Fin de la región crítica
	c.Saldo += monto
}

// Función para retirar dinero de una cuenta
func (c *Cuenta) Retirar(monto int) bool {
	c.Mu.Lock()         // Inicio de la región crítica
	defer c.Mu.Unlock() // Asegurar el desbloqueo al finalizar la función

	if c.Saldo >= monto {
		c.Saldo -= monto
		return true
	}
	return false
}

// Función para transferir dinero entre dos cuentas
func Transferir(cuentaOrigen, cuentaDestino *Cuenta, monto int) bool {
	if cuentaOrigen == cuentaDestino {
		return false
	}

	// Ordenar los bloqueos para evitar deadlock usando las direcciones de memoria
	if uintptr(unsafe.Pointer(cuentaOrigen)) < uintptr(unsafe.Pointer(cuentaDestino)) {
		cuentaOrigen.Mu.Lock()
		defer cuentaOrigen.Mu.Unlock()
		cuentaDestino.Mu.Lock()
		defer cuentaDestino.Mu.Unlock()
	} else {
		cuentaDestino.Mu.Lock()
		defer cuentaDestino.Mu.Unlock()
		cuentaOrigen.Mu.Lock()
		defer cuentaOrigen.Mu.Unlock()
	}

	if cuentaOrigen.Saldo >= monto {
		cuentaOrigen.Saldo -= monto
		cuentaDestino.Saldo += monto
		return true
	}
	return false
}

// Función que simula las transacciones de un cliente
func Cliente(id int, cuentas []*Cuenta, c chan bool) {
	for i := 0; i < 10; i++ { // Realizar 10 transacciones por cliente
		op := rand.Intn(3) // Elegir una operación aleatoria: 0=Depositar, 1=Retirar, 2=Transferir
		cuenta := cuentas[rand.Intn(len(cuentas))]

		switch op {
		case 0: // Depositar
			monto := rand.Intn(100)
			cuenta.Depositar(monto)
			fmt.Printf("Cliente %d depositó %d en cuenta %p. Nuevo saldo: %d\n", id, monto, cuenta, cuenta.Saldo)
		case 1: // Retirar
			monto := rand.Intn(100)
			if cuenta.Retirar(monto) {
				fmt.Printf("Cliente %d retiró %d de cuenta %p. Nuevo saldo: %d\n", id, monto, cuenta, cuenta.Saldo)
			} else {
				fmt.Printf("Cliente %d intentó retirar %d de cuenta %p, pero no tenía suficientes fondos.\n", id, monto, cuenta)
			}
		case 2: // Transferir
			cuentaDestino := cuentas[rand.Intn(len(cuentas))]
			monto := rand.Intn(100)
			if Transferir(cuenta, cuentaDestino, monto) {
				fmt.Printf("Cliente %d transfirió %d de cuenta %p a cuenta %p. Nuevo saldo cuenta origen: %d, nuevo saldo cuenta destino: %d\n", id, monto, cuenta, cuentaDestino, cuenta.Saldo, cuentaDestino.Saldo)
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
		cuentas[i] = &Cuenta{Saldo: 1000}
	}

	// Canal para sincronizar la finalización de los clientes
	c := make(chan bool)

	// Lanzar clientes
	for i := 0; i < 10; i++ {
		go Cliente(i, cuentas, c)
	}

	// Esperar a que todos los clientes terminen
	for i := 0; i < 10; i++ {
		<-c
		fmt.Println("¡Un cliente ha terminado su trabajo!")
	}

	// Mostrar los saldos finales de las cuentas
	fmt.Println("Transacciones completadas. Saldos finales:")
	for i, cuenta := range cuentas {
		fmt.Printf("Cuenta %d: %d\n", i, cuenta.Saldo)
	}
}
