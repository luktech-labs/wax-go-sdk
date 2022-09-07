package wax

type GetTableRowsPayload struct {
	Json          bool   `json:"json"`
	Code          string `json:"code"`
	Scope         string `json:"scope"`
	Table         string `json:"table"`
	LowerBound    string `json:"lower_bound"`
	UpperBound    string `json:"upper_bound"`
	IndexPosition int    `json:"index_position"`
	KeyType       string `json:"key_type,omitempty"`
	Limit         string `json:"limit"`
}

type GetInfoResponse struct {
	ServerVersion             string `json:"server_version"`
	ChainID                   string `json:"chain_id"`
	HeadBlockNum              int    `json:"head_block_num"`
	LastIrreversibleBlockNum  int    `json:"last_irreversible_block_num"`
	LastIrreversibleBlockID   string `json:"last_irreversible_block_id"`
	HeadBlockID               string `json:"head_block_id"`
	HeadBlockTime             string `json:"head_block_time"`
	HeadBlockProducer         string `json:"head_block_producer"`
	VirtualBlockCpuLimit      int    `json:"virtual_block_cpu_limit"`
	VirtualBlockNetLimit      int    `json:"virtual_block_net_limit"`
	BlockCpuLimit             int    `json:"block_cpu_limit"`
	BlockNetLimit             int    `json:"block_net_limit"`
	ServerVersionString       string `json:"server_version_string"`
	ForkDBHeadBlockNum        int    `json:"fork_db_head_block_num"`
	ForkDBHeadBlockID         string `json:"fork_db_head_block_id"`
	ServerFullVersionString   string `json:"server_full_version_string"`
	LastIrreversibleBlockTime string `json:"last_irreversible_block_time"`
}
