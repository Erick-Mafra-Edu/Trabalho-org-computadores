package cache

import "testing"

func TestReadAddressesSupportsHexDecimalCommentsAndBlankLines(t *testing.T) {
	path := createTempAddressFile(t, `
# comentário
0x0010

256
0x00FF # comentário ao lado
`)

	addresses, err := ReadAddresses(path, 16)
	if err != nil {
		t.Fatalf("ReadAddresses retornou erro: %v", err)
	}
	if len(addresses) != 3 {
		t.Fatalf("quantidade de endereços incorreta: got %d want 3", len(addresses))
	}
	want := []string{"0x0010", "256", "0x00FF"}
	for i := range want {
		if addresses[i].Raw != want[i] {
			t.Fatalf("endereço %d incorreto: got %q want %q", i, addresses[i].Raw, want[i])
		}
	}
}

func TestReadAddressesRejectsInvalidAndOutOfRangeAddresses(t *testing.T) {
	tests := []struct {
		name    string
		content string
		bits    uint
	}{
		{name: "invalido", content: "xyz\n", bits: 16},
		{name: "fora-do-espaco", content: "0x10000\n", bits: 16},
		{name: "vazio", content: "# apenas comentario\n\n", bits: 16},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ReadAddresses(createTempAddressFile(t, tt.content), tt.bits); err == nil {
				t.Fatalf("ReadAddresses deveria retornar erro")
			}
		})
	}
}
