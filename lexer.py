import sys
import os

NUM_ESPACIOS = 40
SYMBOLS = {
    "=": "Asignación",
    "+": "Suma",
    "-": "Resta",
    "*": "Multiplicación",
    "/": "División",
    "^": "Potencia",
    "(": "Paréntesis que abre",
    ")": "Paréntesis que cierra"
}

def output(word: str, result: str, outputFile):
    """Writes the analyzed token and its classification to the output file."""
    word = word.strip()  # Clean whitespace around the token
    if word not in [" ", "\n", "\t", ""] and len(word) > 0:
        outputFile.write(f"{word.ljust(NUM_ESPACIOS)}{result}\n")

def processLine(line, outputFile):
    """Analyzes each line and classifies each token according to the defined lexical rules."""
    state = 0  # Initial state
    buffer = ""  # Accumulates characters for the current token
    index = 0  # Current position in the line
    length = len(line)  # Total length of the line

    while index < length:
        char = line[index]  # Current character

        if char == "/" and index + 1 < length and line[index+1] == "/":  # Start of a comment
            buffer += line[index:].strip()  # Take everything to the end as a comment
            output(buffer, "Comentario", outputFile)
            break  # No need to check the rest of the line
        elif state == 0:  # Initial state, looking for the start of a token
            if char.isalpha():  # Start of a variable
                buffer = char
                state = 1  # Variable state
            elif char.isdigit() or (char == "-" and index + 1 < length and (line[index+1].isdigit() or line[index+1] == '.')):  # Start of a number, possibly negative
                buffer = char
                state = 2  # Number state
            elif char in SYMBOLS:  # Symbol
                if buffer:  # If there's something in the buffer, output it
                    output(buffer, "Error", outputFile)
                    buffer = ""
                output(char, SYMBOLS[char], outputFile)
            elif not char.isspace():  # Any other character not whitespace
                buffer = char
                state = 6  # Error state
        else:  # Token processing states
            if char.isspace() or char in SYMBOLS or index == length - 1:  # Token delimiters
                if buffer.endswith('E') or buffer.endswith('e'):  # Exponential notation
                    if index + 1 < length and (line[index+1] == '-' or line[index+1].isdigit()):
                        buffer += char  # Include 'E' or 'e' in the buffer
                        index += 1
                        continue
                if state == 2 and ('E' in buffer or 'e' in buffer or '.' in buffer or buffer.startswith('-')):  # Check for real number
                    token_type = "Real"
                elif state == 2:  # Integer
                    token_type = "Entero"
                elif state == 1:  # Variable
                    token_type = "Variable"
                else:
                    token_type = "Error"
                
                if index == length - 1 and not char.isspace() and not char in SYMBOLS:  # Add last character if it's part of the token
                    buffer += char
                output(buffer, token_type, outputFile)
                buffer = ""
                state = 0  # Reset to initial state
                
                if char in SYMBOLS and char != '/':  # If the current char is a symbol, output it immediately
                    output(char, SYMBOLS[char], outputFile)
            else:
                buffer += char  # Continue building the token

        index += 1  # Move to the next character

def lexerAritmetico(nombre_archivo):
    """Main function to read the file and analyze each line."""
    try:
        with open(nombre_archivo, 'r') as inputFile, open("output.txt", 'w', encoding='utf-8') as outputFile:
            # Write the header
            outputFile.write("Token".ljust(NUM_ESPACIOS) + "Tipo\n")
            for line in inputFile:
                processLine(line, outputFile)
    except FileNotFoundError:
        print(f'ERROR: El archivo "{nombre_archivo}" no se ha encontrado.')

if __name__ == '__main__':
    if len(sys.argv) != 2:
        print("USO: python script.py [ARCHIVO_CON_EXPRESIONES.txt]")
        sys.exit()

    archivo = sys.argv[1]

    if not archivo.endswith(".txt"):
        print("Debes proveer un archivo .txt")
        sys.exit()

    if not os.path.isfile(archivo):
        print("Este archivo no se encuentra en el directorio actual")
        sys.exit()

    lexerAritmetico(archivo)
