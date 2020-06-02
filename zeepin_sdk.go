

//Ontolog sdk in golang. Using for operation with ontology
package Zeepin_Go_Sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/zeepin/Zeepin-Go-Sdk/bip32"
	"github.com/zeepin/Zeepin-Go-Sdk/bip44"
	"github.com/zeepin/ZeepinChain/smartcontract/event"
	"github.com/tyler-smith/go-bip39"
	"io"
	"math/rand"
	"time"

	 "github.com/zeepin/ZeepinChain-Crypto/keypair"
	"github.com/zeepin/Zeepin-Go-Sdk/client"
	"github.com/zeepin/Zeepin-Go-Sdk/utils"
	"github.com/zeepin/ZeepinChain/common"
	common2 "github.com/zeepin/ZeepinChain/common"
	"github.com/zeepin/ZeepinChain/common/constants"
	"github.com/zeepin/ZeepinChain/core/payload"
	"github.com/zeepin/ZeepinChain/core/types"
	s "github.com/zeepin/ZeepinChain-Crypto/signature"
	"github.com/zeepin/Zeepin-Go-Sdk/contract"
	"github.com/zeepin/Zeepin-Go-Sdk/account"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

//ZeepinSdk is the main struct for user
type ZeepinSdk struct {
	client.ClientMgr
	Native *contract.NativeContract
	WasmVM  *contract.WasmVMContract
}

//NewZeepinSdk return ZeepinSdk.
func NewZeepinSdk() *ZeepinSdk {
	zptSdk := &ZeepinSdk{}
	native := contract.NewNativeContract(zptSdk)
	zptSdk.Native = native
	wasmVM := contract.NewWasmVMContract(zptSdk)
	zptSdk.WasmVM = wasmVM
	return zptSdk
}

//CreateWallet return a new wallet
func (this *ZeepinSdk) CreateWallet(walletFile string) (*account.Wallet, error) {
	if utils.IsFileExist(walletFile) {
		return nil, fmt.Errorf("wallet:%s has already exist", walletFile)
	}
	return account.NewWallet(walletFile), nil
}

//OpenWallet return a wallet instance
func (this *ZeepinSdk) OpenWallet(walletFile string) (*account.Wallet, error) {
	return account.OpenWallet(walletFile)
}

func (this *ZeepinSdk) NewAccountFromPrivateKey(privateKey []byte, signatureScheme s.SignatureScheme) (*account.Account, error){
	return account.NewAccountFromPrivateKey(privateKey, signatureScheme)
}

func ParseNativeTxPayload(raw []byte) (map[string]interface{}, error) {
	tx, err := types.TransactionFromRawBytes(raw)
	if err != nil {
		return nil, err
	}
	invokeCode, ok := tx.Payload.(*payload.InvokeCode)
	if !ok {
		return nil, fmt.Errorf("error payload")
	}
	code := invokeCode.Code
	return ParsePayload(code)
}

func ParsePayload(code []byte) (map[string]interface{}, error) {
	codeHex := common.ToHexString(code)
	l := len(code)
	if l > 44 && string(code[l-22:]) == "Ontology.Native.Invoke" {
		if l > 54 && string(code[l-46-8:l-46]) == "transfer" {
			source := common.NewZeroCopySource(code)
			err := ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			from, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res := make(map[string]interface{})
			res["functionName"] = "transfer"
			res["from"] = from.ToBase58()
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			to, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["to"] = to.ToBase58()
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			var amount = uint64(0)
			if string(codeHex[source.Pos()*2]) == "5" || string(codeHex[source.Pos()*2]) == "6" {
				data, eof := source.NextByte()
				if eof {
					return nil, io.ErrUnexpectedEOF
				}
				b := common.BigIntFromEmbeddedBytes([]byte{data})
				amount = b.Uint64() - 0x50
			} else {
				amountBytes, _, irregular, eof := source.NextVarBytes()
				if irregular || eof {
					return nil, io.ErrUnexpectedEOF
				}
				amount = common.BigIntFromEmbeddedBytes(amountBytes).Uint64()
			}

			res["amount"] = amount
			if common.ToHexString(common2.ToArrayReverse(code[l-25-20:l-25])) == contract.ZPT_CONTRACT_ADDRESS.ToHexString() {
				res["asset"] = "zpt"
			} else if common.ToHexString(common2.ToArrayReverse(code[l-25-20:l-25])) == contract.GALA_CONTRACT_ADDRESS.ToHexString() {
				res["asset"] = "gala"
			} else {
				return nil, fmt.Errorf("not zpt or gala contractAddress")
			}
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			//method name
			_, _, irregular, eof := source.NextVarBytes()
			if irregular || eof {
				return nil, io.ErrUnexpectedEOF
			}
			//contract address
			contractAddress, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["contractAddress"] = contractAddress
			return res, nil
		} else if l > 58 && string(code[l-46-12:l-46]) == "transferFrom" {
			res := make(map[string]interface{})
			res["functionName"] = "transferFrom"
			source := common.NewZeroCopySource(code)
			err := ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			sender, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["sender"] = sender.ToBase58()

			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			from, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["from"] = from.ToBase58()
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			to, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["to"] = to.ToBase58()
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			var amount = uint64(0)
			if string(codeHex[source.Pos()*2]) == "5" || string(codeHex[source.Pos()*2]) == "6" {
				//read amount
				data, eof := source.NextByte()
				if eof {
					return nil, io.ErrUnexpectedEOF
				}
				b := common.BigIntFromEmbeddedBytes([]byte{data})
				amount = b.Uint64() - 0x50
			} else {
				amountBytes, _, irregular, eof := source.NextVarBytes()
				if irregular || eof {
					return nil, io.ErrUnexpectedEOF
				}
				amount = common.BigIntFromEmbeddedBytes(amountBytes).Uint64()
			}
			res["amount"] = amount
			if common.ToHexString(common2.ToArrayReverse(code[l-25-20:l-25])) == contract.ZPT_CONTRACT_ADDRESS.ToHexString() {
				res["asset"] = "zpt"
			} else if common.ToHexString(common2.ToArrayReverse(code[l-25-20:l-25])) == contract.GALA_CONTRACT_ADDRESS.ToHexString() {
				res["asset"] = "gala"
				res["amount"] = amount
			}
			err = ignoreOpCode(source)
			if err != nil {
				return nil, err
			}
			source.BackUp(1)
			//method name
			_, _, irregular, eof := source.NextVarBytes()
			if irregular || eof {
				return nil, io.ErrUnexpectedEOF
			}
			//contract address
			contractAddress, err := readAddress(source)
			if err != nil {
				return nil, err
			}
			res["contractAddress"] = contractAddress
			return res, nil
		}
	}
	return nil, fmt.Errorf("not native transfer and transferFrom transaction")
}

func readAddress(source *common.ZeroCopySource) (common.Address, error) {
	senderBytes, _, irregular, eof := source.NextVarBytes()
	if irregular || eof {
		return common.ADDRESS_EMPTY, io.ErrUnexpectedEOF
	}
	sender, err := utils.AddressParseFromBytes(senderBytes)
	if err != nil {
		return common.ADDRESS_EMPTY, err
	}
	return sender, nil
}
func ignoreOpCode(source *common.ZeroCopySource) error {
	s := source.Size()
	for {
		if source.Pos() >= s {
			return nil
		}
		by, eof := source.NextByte()
		if eof {
			return io.EOF
		}
		if contract.OPCODE_IN_PAYLOAD[by] {
			continue
		} else {
			return nil
		}
	}
}

func (this *ZeepinSdk) GenerateMnemonicCodesStr() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func (this *ZeepinSdk) GetPrivateKeyFromMnemonicCodesStrBip44(mnemonicCodesStr string, index uint32) ([]byte, error) {
	if mnemonicCodesStr == "" {
		return nil, fmt.Errorf("mnemonicCodesStr should not be nil")
	}
	//address_index
	if index < 0 {
		return nil, fmt.Errorf("index should be bigger than 0")
	}
	seed := bip39.NewSeed(mnemonicCodesStr, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return nil, err
	}
	//m / purpose' / coin_type' / account' / change / address_index
	//coin type 1024'
	coin := 0x80000400
	//account 0'
	account := 0x80000000
	key, err := bip44.NewKeyFromMasterKey(masterKey, uint32(coin), uint32(account), 0, index)
	if err != nil {
		return nil, err
	}
	keyBytes, err := key.Serialize()
	if err != nil {
		return nil, err
	}
	return keyBytes[46:78], nil
}

//NewInvokeTransaction return smart contract invoke transaction
func (this *ZeepinSdk) NewInvokeTransaction(gasPrice, gasLimit uint64, invokeCode []byte) *types.MutableTransaction {
	invokePayload := &payload.InvokeCode{
		Code: invokeCode,
	}
	tx := &types.MutableTransaction{
		GasPrice: gasPrice,
		GasLimit: gasLimit,
		TxType:   types.Invoke,
		Nonce:    rand.Uint32(),
		Payload:  invokePayload,
		Attributes: 1,
		Sigs:     make([]types.Sig, 0, 0),
	}
	return tx
}

func (this *ZeepinSdk) SignToTransaction(tx *types.MutableTransaction, signer account.Signer) error {
	if tx.Payer == common.ADDRESS_EMPTY {
		account, ok := signer.(*account.Account)
		if ok {
			tx.Payer = account.Address
		}
	}
	for _, sigs := range tx.Sigs {
		if utils.PubKeysEqual([]keypair.PublicKey{signer.GetPublicKey()}, sigs.PubKeys ){
			//have already signed
			return nil
		}
	}
	txHash := tx.Hash()
	sigData, err := signer.Sign(txHash.ToArray())
	if err != nil {
		return fmt.Errorf("sign error:%s", err)
	}
	if tx.Sigs == nil {
		tx.Sigs = make([]types.Sig, 0)
	}
	tx.Sigs = append(tx.Sigs, types.Sig{
		SigData: [][]byte{sigData},
		PubKeys: []keypair.PublicKey{signer.GetPublicKey()},
		M:       1,

	})
	return nil
}

func (this *ZeepinSdk) MultiSignToTransaction(tx *types.MutableTransaction, m uint16, pubKeys []keypair.PublicKey, signer account.Signer) error {
	pkSize := len(pubKeys)
	if m == 0 || int(m) > pkSize || pkSize > constants.MULTI_SIG_MAX_PUBKEY_SIZE {
		return fmt.Errorf("both m and number of pub key must larger than 0, and small than %d, and m must smaller than pub key number", constants.MULTI_SIG_MAX_PUBKEY_SIZE)
	}
	validPubKey := false
	for _, pk := range pubKeys {
		if keypair.ComparePublicKey(pk, signer.GetPublicKey()) {
			validPubKey = true
			break
		}
	}
	if !validPubKey {
		return fmt.Errorf("invalid signer")
	}
	if tx.Payer == common.ADDRESS_EMPTY {
		payer, err := types.AddressFromMultiPubKeys(pubKeys, int(m))
		if err != nil {
			return fmt.Errorf("AddressFromMultiPubKeys error:%s", err)
		}
		tx.Payer = payer
	}
	txHash := tx.Hash()
	if len(tx.Sigs) == 0 {
		tx.Sigs = make([]types.Sig, 0)
	}
	sigData, err := signer.Sign(txHash.ToArray())
	if err != nil {
		return fmt.Errorf("sign error:%s", err)
	}
	hasMutilSig := false
	for i, sigs := range tx.Sigs {
		if utils.PubKeysEqual(sigs.PubKeys, pubKeys) {
			hasMutilSig = true
			if utils.HasAlreadySig(txHash.ToArray(), signer.GetPublicKey(), sigs.SigData) {
				break
			}
			sigs.SigData = append(sigs.SigData, sigData)
			tx.Sigs[i] = sigs
			break
		}
	}
	if !hasMutilSig {
		tx.Sigs = append(tx.Sigs, types.Sig{
			PubKeys: pubKeys,
			M:       m,
			SigData: [][]byte{sigData},
		})
	}
	return nil
}

func (this *ZeepinSdk) GetTxData(tx *types.MutableTransaction) (string, error) {
	txData, err := tx.IntoImmutable()
	if err != nil {
		return "", fmt.Errorf("IntoImmutable error:%s", err)
	}
	sink := common2.ZeroCopySink{}
	txData.Serialization(&sink)
	rawtx := hex.EncodeToString(sink.Bytes())
	return rawtx, nil
}

type TransferEvent struct {
	FuncName string
	From     string
	To       string
	Amount   uint64
}

func (this *ZeepinSdk) ParseNaitveTransferEvent(event *event.NotifyEventInfo) (*TransferEvent, error) {
	if event == nil {
		return nil, fmt.Errorf("event is nil")
	}
	state, ok := event.States.([]interface{})
	if !ok {
		return nil, fmt.Errorf("state.States is not []interface")
	}
	if len(state) != 4 {
		return nil, fmt.Errorf("state length is not 4")
	}
	funcName, ok := state[0].(string)
	if !ok {
		return nil, fmt.Errorf("state.States[0] is not string")
	}
	if funcName != "transfer" {
		return nil, fmt.Errorf("funcName is not transfer")
	} else {
		from, ok := state[1].(string)
		if !ok {
			return nil, fmt.Errorf("state[1] is not string")
		}
		to, ok := state[2].(string)
		if !ok {
			return nil, fmt.Errorf("state[2] is not string")
		}
		amount, ok := state[3].(uint64)
		if !ok {
			return nil, fmt.Errorf("state[3] is not uint64")
		}
		return &TransferEvent{
			FuncName: "transfer",
			From:     from,
			To:       to,
			Amount:   uint64(amount),
		}, nil
	}
}

func (this *ZeepinSdk) GetMutableTx(rawTx string) (*types.MutableTransaction, error) {
	txData, err := hex.DecodeString(rawTx)
	if err != nil {
		return nil, fmt.Errorf("RawTx hex decode error:%s", err)
	}
	tx, err := types.TransactionFromRawBytes(txData)
	if err != nil {
		return nil, fmt.Errorf("TransactionFromRawBytes error:%s", err)
	}
	mutTx, err := tx.IntoMutable()
	if err != nil {
		return nil, fmt.Errorf("[ZPT]IntoMutable error:%s", err)
	}
	return mutTx, nil
}

func (this *ZeepinSdk) GetMultiAddr(pubkeys []keypair.PublicKey, m int) (string, error) {
	addr, err := types.AddressFromMultiPubKeys(pubkeys, m)
	if err != nil {
		return "", fmt.Errorf("GetMultiAddrs error:%s", err)
	}
	return addr.ToBase58(), nil
}

func (this *ZeepinSdk) GetAdddrByPubKey(pubKey keypair.PublicKey) string {
	address := types.AddressFromPubKey(pubKey)
	return address.ToBase58()
}
