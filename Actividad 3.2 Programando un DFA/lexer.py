import sys
import os

NUM_ESPACIOS = 40
SIMBOLOS = {
    "=": "Asignación",
    "+": "Suma",
    "-": "Resta",
    "*": "Multiplicación",
    "/": "División",
    "^": "Potencia",
    "(": "Paréntesis que abre",
    ")": "Paréntesis que cierra"
}

def salida(palabra: str, resultado: str, archivo_salida):
    palabra = palabra.strip()  # Limpiar espacios en blanco alrededor del token
    if palabra not in [" ", "\n", "\t", ""] and len(palabra) > 0:
        archivo_salida.write(f"{palabra.ljust(NUM_ESPACIOS)}{resultado}\n")

def procesarLinea(linea, archivo_salida):
    estado = 0  # Estado inicial
    buffer = ""  # Acumula caracteres para el token actual
    indice = 0  # Posición actual en la línea
    longitud = len(linea)  # Longitud total de la línea

    while indice < longitud:
        caracter = linea[indice]  # Carácter actual

        if caracter == "/" and indice + 1 < longitud and linea[indice+1] == "/":  # Inicio de un comentario
            buffer += linea[indice:].strip()  # Tomar todo hasta el final como un comentario
            salida(buffer, "Comentario", archivo_salida)
            break  # No es necesario revisar el resto de la línea
        elif estado == 0:  # Estado inicial, buscando el inicio de un token
            if caracter.isalpha():  # Inicio de una variable
                buffer = caracter
                estado = 1  # Estado de variable
            elif caracter.isdigit() or (caracter == "-" and indice + 1 < longitud and (linea[indice+1].isdigit() or linea[indice+1] == '.')):  # Inicio de un número, posiblemente negativo
                buffer = caracter
                estado = 2  # Estado de número
            elif caracter in SIMBOLOS:  # Símbolo
                if buffer:  # Si hay algo en el buffer, sacarlo
                    salida(buffer, "Error", archivo_salida)
                    buffer = ""
                salida(caracter, SIMBOLOS[caracter], archivo_salida)
            elif not caracter.isspace():  # Cualquier otro carácter que no sea espacio en blanco
                buffer = caracter
                estado = 6  # Estado de error
        else:  # Estados de procesamiento de tokens
            if caracter.isspace() or caracter in SIMBOLOS or indice == longitud - 1:  # Delimitadores de tokens
                if buffer.endswith('E') or buffer.endswith('e'):  # Notación exponencial
                    if indice + 1 < longitud and (linea[indice+1] == '-' or linea[indice+1].isdigit()):
                        buffer += caracter  # Incluir 'E' o 'e' en el buffer
                        indice += 1
                        continue
                if estado == 2 and ('E' in buffer or 'e' in buffer or '.' in buffer or buffer.startswith('-')):  # Verificar si es número real
                    tipo_token = "Real"
                elif estado == 2:  # Entero
                    tipo_token = "Entero"
                elif estado == 1:  # Variable
                    tipo_token = "Variable"
                else:
                    tipo_token = "Error"
                
                if indice == longitud - 1 and not caracter.isspace() and not caracter in SIMBOLOS:  # Agregar último carácter si es parte del token
                    buffer += caracter
                salida(buffer, tipo_token, archivo_salida)
                buffer = ""
                estado = 0  # Reiniciar al estado inicial
                
                if caracter in SIMBOLOS and caracter != '/':  # Si el carácter actual es un símbolo, sacarlo inmediatamente
                    salida(caracter, SIMBOLOS[caracter], archivo_salida)
            else:
                buffer += caracter  # Continuar construyendo el token

        indice += 1  # Mover al siguiente carácter

def lexerAritmetico(archivo):
    try:
        with open(archivo, 'r') as archivo_entrada, open("salida.txt", 'w', encoding='utf-8') as archivo_salida:
            # Escribir el encabezado
            archivo_salida.write("Token".ljust(NUM_ESPACIOS) + "Tipo\n")
            for linea in archivo_entrada:
                procesarLinea(linea, archivo_salida)
    except FileNotFoundError:
        print(f'ERROR: El archivo "{archivo}" no se ha encontrado.')

if __name__ == '__main__':
    archivo = sys.argv[1]
    lexerAritmetico(archivo)
