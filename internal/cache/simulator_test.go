package cache

import (
	"io"
	"testing"
)

func TestCalculateAddressFields(t *testing.T) {
	layout := BitLayout{OffsetBits: 4, IndexBits: 4, TagBits: 8, NumLines: 16, NumSets: 16}
	fields, err := CalculateAddressFields(addressInput{Raw: "0x12A5", Line: 1}, layout, 16)
	if err != nil {
		t.Fatalf("CalculateAddressFields retornou erro: %v", err)
	}
	if fields.Tag != 0x12 || fields.Index != 0xA || fields.Offset != 0x5 {
		t.Fatalf("campos incorretos: tag=%X index=%X offset=%X", fields.Tag, fields.Index, fields.Offset)
	}
	if fields.Binary != "0001001010100101" {
		t.Fatalf("binário incorreto: %s", fields.Binary)
	}
	if fields.TagBits != "00010010" || fields.IndexBits != "1010" || fields.OffsetBits != "0101" {
		t.Fatalf("particionamento incorreto: tag=%s index=%s offset=%s", fields.TagBits, fields.IndexBits, fields.OffsetBits)
	}
}

func TestSimulateDirectMappingCountsHitsAndMisses(t *testing.T) {
	config := Config{
		CacheSize:     256,
		BlockSize:     16,
		Associativity: 1,
		AddressBits:   16,
		Policy:        PolicyLRU,
		InputPath:     "teste",
	}
	layout := BitLayout{OffsetBits: 4, IndexBits: 4, TagBits: 8, NumLines: 16, NumSets: 16}
	addresses := []addressInput{
		{Raw: "0x0010", Line: 1},
		{Raw: "0x0014", Line: 2},
		{Raw: "0x0110", Line: 3},
		{Raw: "0x0010", Line: 4},
	}

	result, err := SimulateDirectMapping(config, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateDirectMapping retornou erro: %v", err)
	}
	if result.TotalAccesses != 4 || result.Hits != 1 || result.Misses != 3 {
		t.Fatalf("resultado incorreto: acessos=%d hits=%d misses=%d", result.TotalAccesses, result.Hits, result.Misses)
	}
}

func TestSetAssociativeLRUAndFIFOReplacement(t *testing.T) {
	layout := BitLayout{OffsetBits: 0, IndexBits: 0, TagBits: 8, NumLines: 2, NumSets: 1}
	addresses := []addressInput{
		{Raw: "0", Line: 1},
		{Raw: "1", Line: 2},
		{Raw: "0", Line: 3},
		{Raw: "2", Line: 4},
		{Raw: "1", Line: 5},
	}
	baseConfig := Config{
		CacheSize:     2,
		BlockSize:     1,
		Associativity: 2,
		AddressBits:   8,
		InputPath:     "teste",
	}

	lruConfig := baseConfig
	lruConfig.Policy = PolicyLRU
	lruResult, err := SimulateSetAssociative(lruConfig, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateSetAssociative LRU retornou erro: %v", err)
	}
	if lruResult.Hits != 1 || lruResult.Misses != 4 {
		t.Fatalf("resultado LRU incorreto: hits=%d misses=%d", lruResult.Hits, lruResult.Misses)
	}

	fifoConfig := baseConfig
	fifoConfig.Policy = PolicyFIFO
	fifoResult, err := SimulateSetAssociative(fifoConfig, layout, addresses, io.Discard)
	if err != nil {
		t.Fatalf("SimulateSetAssociative FIFO retornou erro: %v", err)
	}
	if fifoResult.Hits != 2 || fifoResult.Misses != 3 {
		t.Fatalf("resultado FIFO incorreto: hits=%d misses=%d", fifoResult.Hits, fifoResult.Misses)
	}
}
