package cache

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateParametersExamplesFromPDF(t *testing.T) {
	inputPath := createTempAddressFile(t, "0x0000\n")
	tests := []struct {
		name       string
		config     Config
		wantOffset uint
		wantIndex  uint
		wantTag    uint
	}{
		{
			name: "256-16-direta-16bits",
			config: Config{
				CacheSize:     256,
				BlockSize:     16,
				Associativity: 1,
				AddressBits:   16,
				Policy:        PolicyLRU,
				InputPath:     inputPath,
			},
			wantOffset: 4,
			wantIndex:  4,
			wantTag:    8,
		},
		{
			name: "1024-32-2way-16bits",
			config: Config{
				CacheSize:     1024,
				BlockSize:     32,
				Associativity: 2,
				AddressBits:   16,
				Policy:        PolicyLRU,
				InputPath:     inputPath,
			},
			wantOffset: 5,
			wantIndex:  4,
			wantTag:    7,
		},
		{
			name: "512-8-4way-16bits",
			config: Config{
				CacheSize:     512,
				BlockSize:     8,
				Associativity: 4,
				AddressBits:   16,
				Policy:        PolicyLRU,
				InputPath:     inputPath,
			},
			wantOffset: 3,
			wantIndex:  4,
			wantTag:    9,
		},
		{
			name: "2048-64-8way-32bits",
			config: Config{
				CacheSize:     2048,
				BlockSize:     64,
				Associativity: 8,
				AddressBits:   32,
				Policy:        PolicyLRU,
				InputPath:     inputPath,
			},
			wantOffset: 6,
			wantIndex:  2,
			wantTag:    24,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			layout, err := ValidateParameters(tt.config)
			if err != nil {
				t.Fatalf("ValidateParameters retornou erro: %v", err)
			}
			if layout.OffsetBits != tt.wantOffset || layout.IndexBits != tt.wantIndex || layout.TagBits != tt.wantTag {
				t.Fatalf("campos incorretos: offset=%d index=%d tag=%d", layout.OffsetBits, layout.IndexBits, layout.TagBits)
			}
		})
	}
}

func TestValidateParametersRejectsInvalidConfig(t *testing.T) {
	inputPath := createTempAddressFile(t, "0x0000\n")
	tests := []Config{
		{CacheSize: 300, BlockSize: 16, Associativity: 1, AddressBits: 16, Policy: PolicyLRU, InputPath: inputPath},
		{CacheSize: 256, BlockSize: 24, Associativity: 1, AddressBits: 16, Policy: PolicyLRU, InputPath: inputPath},
		{CacheSize: 256, BlockSize: 16, Associativity: 3, AddressBits: 16, Policy: PolicyLRU, InputPath: inputPath},
		{CacheSize: 256, BlockSize: 16, Associativity: 1, AddressBits: 16, Policy: "RANDOM", InputPath: inputPath},
		{CacheSize: 256, BlockSize: 16, Associativity: 1, AddressBits: 16, Policy: PolicyLRU, InputPath: filepath.Join(t.TempDir(), "nao-existe.txt")},
	}

	for _, config := range tests {
		if _, err := ValidateParameters(config); err == nil {
			t.Fatalf("ValidateParameters deveria rejeitar a configuração: %+v", config)
		}
	}
}

func createTempAddressFile(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "enderecos.txt")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("não foi possível criar arquivo temporário: %v", err)
	}
	return path
}
