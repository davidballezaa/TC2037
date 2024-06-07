package banco

import (
	"sync"
	"testing"
)

func TestCuenta(t *testing.T) {
	cuenta := &Cuenta{Saldo: 100}

	// Prueba de Depósito
	cuenta.Depositar(50)
	if cuenta.Saldo != 150 {
		t.Errorf("Esperaba 150, pero obtuvo %d", cuenta.Saldo)
	}

	// Prueba de Retiro con saldo suficiente
	result := cuenta.Retirar(30)
	if !result || cuenta.Saldo != 120 {
		t.Errorf("Esperaba retirar 30 y obtener saldo 120, pero obtuvo %d", cuenta.Saldo)
	}

	// Prueba de Retiro con saldo insuficiente
	result = cuenta.Retirar(200)
	if result || cuenta.Saldo != 120 {
		t.Errorf("Esperaba no poder retirar 200 y mantener saldo 120, pero obtuvo %d", cuenta.Saldo)
	}
}

func TestTransferir(t *testing.T) {
	cuentaA := &Cuenta{Saldo: 200}
	cuentaB := &Cuenta{Saldo: 100}

	// Prueba de Transferencia con saldo suficiente
	result := Transferir(cuentaA, cuentaB, 50)
	if !result || cuentaA.Saldo != 150 || cuentaB.Saldo != 150 {
		t.Errorf("Esperaba transferencia exitosa y saldos 150 en ambas cuentas, pero obtuvo %d y %d", cuentaA.Saldo, cuentaB.Saldo)
	}

	// Prueba de Transferencia con saldo insuficiente
	result = Transferir(cuentaA, cuentaB, 300)
	if result || cuentaA.Saldo != 150 || cuentaB.Saldo != 150 {
		t.Errorf("Esperaba transferencia fallida y saldos sin cambios, pero obtuvo %d y %d", cuentaA.Saldo, cuentaB.Saldo)
	}
}

func TestConcurrencia(t *testing.T) {
	cuenta := &Cuenta{Saldo: 1000}
	var wg sync.WaitGroup
	numDepositos := 100
	numRetiros := 100

	// Ejecutar múltiples depósitos concurrentes
	for i := 0; i < numDepositos; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cuenta.Depositar(10)
		}()
	}

	// Ejecutar múltiples retiros concurrentes
	for i := 0; i < numRetiros; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cuenta.Retirar(10)
		}()
	}

	wg.Wait()

	// Verificar el saldo final
	expectedSaldo := 1000 + (numDepositos * 10) - (numRetiros * 10)
	if cuenta.Saldo != expectedSaldo {
		t.Errorf("Esperaba saldo final de %d, pero obtuvo %d", expectedSaldo, cuenta.Saldo)
	}
}
