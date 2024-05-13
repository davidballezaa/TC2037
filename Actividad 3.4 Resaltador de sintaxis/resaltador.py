import re

# Define las expresiones regulares para diferentes tokens.
regex_patterns = [
    ('COMMENT', r'#.*'),
    ('STRING', r'(\'[^\']*\'|\"[^\"]*\")'),
    ('NUMBER', r'\b\d+(\.\d*)?([eE][+-]?\d+)?\b'),
    ('RESERVED_WORD', r'\b(if|else|elif|for|while|break|continue|return|and|or|not|is|in|import|def|class|try|except|finally|with|as|from|lambda|nonlocal|global|yield|assert|pass|raise|True|False|None)\b'),
    ('IDENTIFIER', r'\b[a-zA-Z_][a-zA-Z0-9_]*\b'),
    ('OPERATOR', r'[\+\-\*/%==!=<>&|^~]'),
    ('DELIMITER', r'[()\[\]{};:,]')
]

def tokenize(source_code):
    tokens = []
    last_end = 0
    while last_end < len(source_code):
        match_found = False
        for category, pattern in regex_patterns:
            regex = re.compile(pattern)
            match = regex.search(source_code, last_end)
            if match and match.start() == last_end:
                tokens.append((category, match.group(), match.start()))
                last_end = match.end()
                match_found = True
                break
        if not match_found:
            last_end += 1
    return sorted(tokens, key=lambda x: x[2])

def highlight_code(source_code, tokens):
    highlighted_html = '<html><head><style>'
    styles = {
        'RESERVED_WORD': 'color: red;',
        'IDENTIFIER': 'color: blue;',
        'OPERATOR': 'color: purple;',
        'DELIMITER': 'color: orange;',
        'NUMBER': 'color: green;',
        'STRING': 'color: brown;',
        'COMMENT': 'color: gray;'
    }
    for category, style in styles.items():
        highlighted_html += f'.{category} {{{style}}}'
    highlighted_html += '</style></head><body>'
    
    # Añade el encabezado con las descripciones de los colores
    highlighted_html += '<h2>Token Colors</h2><ul>'
    for category, style in styles.items():
        highlighted_html += f'<li><span class="{category}">{category}</span>: {style}</li>'
    highlighted_html += '</ul><hr><pre>'

    last_index = 0
    for category, value, start in tokens:
        highlighted_html += f'{source_code[last_index:start]}<span class="{category}">{value}</span>'
        last_index = start + len(value)

    highlighted_html += source_code[last_index:]  # Resto del código sin tokens
    highlighted_html += '</pre></body></html>'

    return highlighted_html

def main():
    with open('entrada.txt', 'r') as file:
        source_code = file.read()
    
    tokens = tokenize(source_code)
    html_output = highlight_code(source_code, tokens)
    
    with open('output.html', 'w') as html_file:
        html_file.write(html_output)

    print("HTML output has been generated and saved as 'output.html'.")

if __name__ == "__main__":
    main()
