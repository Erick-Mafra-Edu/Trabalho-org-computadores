package cache

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// carrega endereços decimais ou hexadecimais e ignora comentários/linhas vazias.
func ReadAddresses(path string, addrBits uint) ([]addressInput, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("não foi possível abrir o arquivo: %w", err)
	}
	defer file.Close()

	var addresses []addressInput
	scanner := bufio.NewScanner(file)
	lineNumber := 0
	maxValue := maxAddressValue(addrBits)

	for scanner.Scan() {
		lineNumber++
		text := scanner.Text()
		if commentIndex := strings.Index(text, "#"); commentIndex >= 0 {
			text = text[:commentIndex]
		}
		text = strings.TrimSpace(text)
		if text == "" {
			continue
		}

		value, err := strconv.ParseUint(text, 0, 64)
		if err != nil {
			return nil, fmt.Errorf("endereço inválido na linha %d: %q", lineNumber, text)
		}
		if value > maxValue {
			return nil, fmt.Errorf("endereço fora do espaço permitido na linha %d: %q excede %d bits", lineNumber, text, addrBits)
		}
		addresses = append(addresses, addressInput{Raw: text, Line: lineNumber})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("erro ao ler o arquivo: %w", err)
	}
	if len(addresses) == 0 {
		return nil, fmt.Errorf("arquivo de entrada vazio ou sem endereços válidos")
	}

	return addresses, nil
}
