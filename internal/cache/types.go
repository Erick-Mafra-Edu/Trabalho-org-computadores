package cache

// representa a estratégia usada quando um conjunto cheio precisa receber uma nova linha.
type ReplacementPolicy string

const (
	PolicyLRU  ReplacementPolicy = "LRU"
	PolicyFIFO ReplacementPolicy = "FIFO"
)

// agrupa todos os parâmetros necessários para configurar a simulação da cache.
type Config struct {
	CacheSize     uint64
	BlockSize     uint64
	Associativity uint64
	AddressBits   uint
	Policy        ReplacementPolicy
	InputPath     string
	Verbose       bool
}

// guarda a quantidade de bits usada por cada campo do endereço.
type BitLayout struct {
	TagBits    uint
	IndexBits  uint
	OffsetBits uint
	NumLines   uint64
	NumSets    uint64
}

// guarda o endereço decomposto usado na busca da cache.
type AddressFields struct {
	Original   string
	Value      uint64
	Binary     string
	TagBits    string
	IndexBits  string
	OffsetBits string
	Tag        uint64
	Index      uint64
	Offset     uint64
}

// representa uma linha dentro de um conjunto da cache.
type CacheLine struct {
	Valid    bool
	Tag      uint64
	LastUsed uint64
	LoadedAt uint64
}

// agrupa linhas que compartilham o mesmo índice.
type CacheSet struct {
	Lines    []CacheLine
	NextFIFO int
}

// representa a cache simulada completa.
type Cache struct {
	Sets   []CacheSet
	Config Config
	Layout BitLayout
	Clock  uint64
}

// guarda contadores e metadados exibidos no relatório final.
type SimulationResult struct {
	Config        Config
	Layout        BitLayout
	TotalAccesses int
	Hits          int
	Misses        int
}

// guarda os dados de entrada para o endereço.
type addressInput struct {
	Raw  string
	Line int
}
