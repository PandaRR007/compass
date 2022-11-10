package matic

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/light"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
	"github.com/mapprotocol/compass/mapprotocol"
	utils "github.com/mapprotocol/compass/shared/ethereum"
	"math/big"
)

type BlockHeader struct {
	ParentHash       []byte         `json:"parentHash"`
	Sha3Uncles       []byte         `json:"sha3Uncles"`
	Miner            common.Address `json:"miner"`
	StateRoot        []byte         `json:"stateRoot"`
	TransactionsRoot []byte         `json:"transactionsRoot"`
	ReceiptsRoot     []byte         `json:"receiptsRoot"`
	LogsBloom        []byte         `json:"logsBloom"`
	Difficulty       *big.Int       `json:"difficulty"`
	Number           *big.Int       `json:"number"`
	GasLimit         *big.Int       `json:"gasLimit"`
	GasUsed          *big.Int       `json:"gasUsed"`
	Timestamp        *big.Int       `json:"timestamp"`
	ExtraData        []byte         `json:"extraData"`
	MixHash          []byte         `json:"mixHash"`
	Nonce            []byte         `json:"nonce"`
	BaseFeePerGas    *big.Int       `json:"baseFeePerGas"`
}

func ConvertHeader(header *types.Header) BlockHeader {
	bloom := make([]byte, 0, len(header.Bloom))
	for _, b := range header.Bloom {
		bloom = append(bloom, b)
	}
	nonce := make([]byte, 0, len(header.Nonce))
	for _, b := range header.Nonce {
		nonce = append(nonce, b)
	}
	return BlockHeader{
		ParentHash:       hashToByte(header.ParentHash),
		Sha3Uncles:       hashToByte(header.UncleHash),
		Miner:            header.Coinbase,
		StateRoot:        hashToByte(header.Root),
		TransactionsRoot: hashToByte(header.TxHash),
		ReceiptsRoot:     hashToByte(header.ReceiptHash),
		LogsBloom:        bloom,
		Difficulty:       header.Difficulty,
		Number:           header.Number,
		GasLimit:         new(big.Int).SetUint64(header.GasLimit),
		GasUsed:          new(big.Int).SetUint64(header.GasUsed),
		Timestamp:        new(big.Int).SetUint64(header.Time),
		ExtraData:        header.Extra,
		MixHash:          hashToByte(header.MixDigest),
		Nonce:            nonce,
		BaseFeePerGas:    header.BaseFee,
	}
}

func hashToByte(h common.Hash) []byte {
	ret := make([]byte, 0, len(h))
	for _, b := range h {
		ret = append(ret, b)
	}
	return ret
}

type ProofData struct {
	Header BlockHeader
	Rp     ReceiptProof
}

type ReceiptProof struct {
	txReceipt *mapprotocol.TxReceipt
	keyIndex  []byte
	proof     [][]byte
}

func AssembleProof(header *types.Header, log types.Log, bridgeAddr common.Address, receipts []*types.Receipt, method string) ([]byte, error) {
	h := ConvertHeader(header)
	txIndex := log.TxIndex
	receipt, err := mapprotocol.GetTxReceipt(receipts[txIndex])
	if err != nil {
		return nil, err
	}

	proof, err := getProof(receipts, txIndex)
	if err != nil {
		return nil, err
	}

	var key []byte
	key = rlp.AppendUint64(key[:0], uint64(txIndex))
	ek := utils.Key2Hex(key, len(proof))

	pd := ProofData{
		Header: h,
		Rp: ReceiptProof{
			txReceipt: receipt,
			keyIndex:  ek,
			proof:     proof,
		},
	}

	input, err := mapprotocol.Matic.Methods[mapprotocol.MethodOfGetBytes].Inputs.Pack(pd)
	if err != nil {
		return nil, err
	}
	pack, err := mapprotocol.PackInput(mapprotocol.LightManger, mapprotocol.MethodVerifyProofData, input)
	if err != nil {
		return nil, err
	}

	return pack, nil
}

func getProof(receipts []*types.Receipt, txIndex uint) ([][]byte, error) {
	tr, err := trie.New(common.Hash{}, trie.NewDatabase(memorydb.New()))
	if err != nil {
		return nil, err
	}

	tr = utils.DeriveTire(receipts, tr)
	ns := light.NewNodeSet()
	key, err := rlp.EncodeToBytes(txIndex)
	if err != nil {
		return nil, err
	}
	if err = tr.Prove(key, 0, ns); err != nil {
		return nil, err
	}

	proof := make([][]byte, 0, len(ns.NodeList()))
	for _, v := range ns.NodeList() {
		proof = append(proof, v)
	}

	return proof, nil
}
