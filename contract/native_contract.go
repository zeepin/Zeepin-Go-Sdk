package contract

import (
	"bytes"
	"encoding/hex"
	"fmt"

	"github.com/zeepin/ZeepinChain-Crypto/keypair"
	"github.com/zeepin/ZeepinChain/common"
	"github.com/zeepin/ZeepinChain/common/serialization"
	"github.com/zeepin/ZeepinChain/core/types"
	cutils "github.com/zeepin/ZeepinChain/http/base/common"
	"github.com/zeepin/ZeepinChain/smartcontract/service/native/global_params"
	"github.com/zeepin/ZeepinChain/smartcontract/service/native/zpt"
	sdkcom "github.com/zeepin/Zeepin-Go-Sdk/common"
	"github.com/zeepin/Zeepin-Go-Sdk/utils"
	sdk "github.com/zeepin/Zeepin-Go-Sdk"
	"github.com/zeepin/Zeepin-Go-Sdk/account"
)

var (
	ZPT_CONTRACT_ADDRESS, _           = utils.AddressFromHexString("0100000000000000000000000000000000000000")
	GALA_CONTRACT_ADDRESS, _          = utils.AddressFromHexString("0200000000000000000000000000000000000000")
	ZPT_ID_CONTRACT_ADDRESS, _        = utils.AddressFromHexString("0300000000000000000000000000000000000000")
	GLOABL_PARAMS_CONTRACT_ADDRESS, _ = utils.AddressFromHexString("0400000000000000000000000000000000000000")
	AUTH_CONTRACT_ADDRESS, _          = utils.AddressFromHexString("0600000000000000000000000000000000000000")
	GOVERNANCE_CONTRACT_ADDRESS, _    = utils.AddressFromHexString("0700000000000000000000000000000000000000")
)

var (
	ZPT_CONTRACT_VERSION           = byte(0)
	GALA_CONTRACT_VERSION          = byte(0)
	ZPT_ID_CONTRACT_VERSION        = byte(0)
	GLOBAL_PARAMS_CONTRACT_VERSION = byte(0)
	AUTH_CONTRACT_VERSION          = byte(0)
	GOVERNANCE_CONTRACT_VERSION    = byte(0)
)

var OPCODE_IN_PAYLOAD = map[byte]bool{0x00: true, 0xc6: true, 0x6b: true, 0x6a: true, 0xc8: true, 0x6c: true, 0x68: true, 0x67: true,
	0x7c: true, 0x51: true, 0xc1: true}

type NativeContract struct {
	zptSdk       *sdk.ZeepinSdk
	Zpt          *Zpt
	Gala         *Gala
	ZptId        *ZptId
	GlobalParams *GlobalParam
	Auth         *Auth
}

func NewNativeContract(zptSdk *sdk.ZeepinSdk) *NativeContract {
	native := &NativeContract{zptSdk: zptSdk}
	native.Zpt = &Zpt{native: native, zptSdk: zptSdk}
	native.Gala = &Gala{native: native, zptSdk: zptSdk}
	native.ZptId = &ZptId{native: native, zptSdk: zptSdk}
	native.GlobalParams = &GlobalParam{native: native, zptSdk: zptSdk}
	native.Auth = &Auth{native: native, zptSdk: zptSdk}
	return native
}

func (this *NativeContract) NewNativeInvokeTransaction(
	gasPrice,
	gasLimit uint64,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (*types.MutableTransaction, error) {
	if params == nil {
		params = make([]interface{}, 0, 1)
	}
	//Params cannot empty, if params is empty, fulfil with empty string
	if len(params) == 0 {
		params = append(params, "")
	}
	invokeCode, err := cutils.BuildNativeInvokeCode(contractAddress, version, method, params)
	if err != nil {
		return nil, fmt.Errorf("BuildNativeInvokeCode error:%s", err)
	}
	return this.zptSdk.NewInvokeTransaction(gasPrice, gasLimit, invokeCode), nil
}

func (this *NativeContract) InvokeNativeContract(
	gasPrice,
	gasLimit uint64,
	singer *account.Account,
	version byte,
	contractAddress common.Address,
	method string,
	params []interface{},
) (common.Uint256, error) {
	tx, err := this.NewNativeInvokeTransaction(gasPrice, gasLimit, version, contractAddress, method, params)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, singer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *NativeContract) PreExecInvokeNativeContract(
	contractAddress common.Address,
	version byte,
	method string,
	params []interface{},
) (*sdkcom.PreExecResult, error) {
	tx, err := this.NewNativeInvokeTransaction(0, 0, version, contractAddress, method, params)
	if err != nil {
		return nil, err
	}
	return this.zptSdk.PreExecTransaction(tx)
}

type Zpt struct {
	zptSdk *sdk.ZeepinSdk
	native *NativeContract
}

func (this *Zpt) NewTransferTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.NewMultiTransferTransaction(gasPrice, gasLimit, []*zpt.State{state})
}

func (this *Zpt) Transfer(gasPrice, gasLimit uint64, from *account.Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Zpt) NewMultiTransferTransaction(gasPrice, gasLimit uint64, states []*zpt.State) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_CONTRACT_VERSION,
		ZPT_CONTRACT_ADDRESS,
		zpt.TRANSFER_NAME,
		[]interface{}{states})
}

func (this *Zpt) MultiTransfer(gasPrice, gasLimit uint64, states []*zpt.State, signer *account.Account) (common.Uint256, error) {
	tx, err := this.NewMultiTransferTransaction(gasPrice, gasLimit, states)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Zpt) NewTransferFromTransaction(gasPrice, gasLimit uint64, sender, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.TransferFrom{
		Sender: sender,
		From:   from,
		To:     to,
		Value:  amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_CONTRACT_VERSION,
		ZPT_CONTRACT_ADDRESS,
		zpt.TRANSFERFROM_NAME,
		[]interface{}{state},
	)
}

func (this *Zpt) TransferFrom(gasPrice, gasLimit uint64, sender *account.Account, from, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferFromTransaction(gasPrice, gasLimit, sender.Address, from, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, sender)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Zpt) NewApproveTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_CONTRACT_VERSION,
		ZPT_CONTRACT_ADDRESS,
		zpt.APPROVE_NAME,
		[]interface{}{state},
	)
}

func (this *Zpt) Approve(gasPrice, gasLimit uint64, from *account.Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewApproveTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Zpt) Allowance(from, to common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.ALLOWANCE_NAME,
		[]interface{}{&allowanceStruct{From: from, To: to}},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Zpt) Symbol() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.SYMBOL_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Zpt) BalanceOf(address common.Address) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.BALANCEOF_NAME,
		[]interface{}{address[:]},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Zpt) Name() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.NAME_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Zpt) Decimals() (byte, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.DECIMALS_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	decimals, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return byte(decimals.Uint64()), nil
}

func (this *Zpt) TotalSupply() (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_CONTRACT_ADDRESS,
		ZPT_CONTRACT_VERSION,
		zpt.TOTAL_SUPPLY_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

type Gala struct {
	zptSdk *sdk.ZeepinSdk
	native *NativeContract
}

func (this *Gala) NewTransferTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.NewMultiTransferTransaction(gasPrice, gasLimit, []*zpt.State{state})
}

func (this *Gala) Transfer(gasPrice, gasLimit uint64, from *account.Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Gala) NewMultiTransferTransaction(gasPrice, gasLimit uint64, states []*zpt.State) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GALA_CONTRACT_VERSION,
		GALA_CONTRACT_ADDRESS,
		zpt.TRANSFER_NAME,
		[]interface{}{states})
}

func (this *Gala) MultiTransfer(gasPrice, gasLimit uint64, states []*zpt.State, signer *account.Account) (common.Uint256, error) {
	tx, err := this.NewMultiTransferTransaction(gasPrice, gasLimit, states)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Gala) NewTransferFromTransaction(gasPrice, gasLimit uint64, sender, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.TransferFrom{
		Sender: sender,
		From:   from,
		To:     to,
		Value:  amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GALA_CONTRACT_VERSION,
		GALA_CONTRACT_ADDRESS,
		zpt.TRANSFERFROM_NAME,
		[]interface{}{state},
	)
}

func (this *Gala) TransferFrom(gasPrice, gasLimit uint64, sender *account.Account, from, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewTransferFromTransaction(gasPrice, gasLimit, sender.Address, from, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, sender)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Gala) NewWithdrawGalaTransaction(gasPrice, gasLimit uint64, address common.Address, amount uint64) (*types.MutableTransaction, error) {
	return this.NewTransferFromTransaction(gasPrice, gasLimit, address, ZPT_CONTRACT_ADDRESS, address, amount)
}

func (this *Gala) WithdrawGala(gasPrice, gasLimit uint64, address *account.Account, amount uint64) (common.Uint256, error) {
	tx, err := this.NewWithdrawGalaTransaction(gasPrice, gasLimit, address.Address, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, address)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Gala) NewApproveTransaction(gasPrice, gasLimit uint64, from, to common.Address, amount uint64) (*types.MutableTransaction, error) {
	state := &zpt.State{
		From:  from,
		To:    to,
		Value: amount,
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GALA_CONTRACT_VERSION,
		GALA_CONTRACT_ADDRESS,
		zpt.APPROVE_NAME,
		[]interface{}{state},
	)
}

func (this *Gala) Approve(gasPrice, gasLimit uint64, from *account.Account, to common.Address, amount uint64) (common.Uint256, error) {
	tx, err := this.NewApproveTransaction(gasPrice, gasLimit, from.Address, to, amount)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, from)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Gala) Allowance(from, to common.Address) (uint64, error) {
	type allowanceStruct struct {
		From common.Address
		To   common.Address
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.ALLOWANCE_NAME,
		[]interface{}{&allowanceStruct{From: from, To: to}},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Gala) UnboundGala(address common.Address) (uint64, error) {
	return this.Allowance(ZPT_CONTRACT_ADDRESS, address)
}

func (this *Gala) Symbol() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.SYMBOL_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Gala) BalanceOf(address common.Address) (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.BALANCEOF_NAME,
		[]interface{}{address[:]},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

func (this *Gala) Name() (string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.NAME_NAME,
		[]interface{}{},
	)
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Gala) Decimals() (byte, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.DECIMALS_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	decimals, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return byte(decimals.Uint64()), nil
}

func (this *Gala) TotalSupply() (uint64, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GALA_CONTRACT_ADDRESS,
		GALA_CONTRACT_VERSION,
		zpt.TOTAL_SUPPLY_NAME,
		[]interface{}{},
	)
	if err != nil {
		return 0, err
	}
	balance, err := preResult.Result.ToInteger()
	if err != nil {
		return 0, err
	}
	return balance.Uint64(), nil
}

type ZptId struct {
	zptSdk *sdk.ZeepinSdk
	native *NativeContract
}

func (this *ZptId) NewRegIDWithPublicKeyTransaction(gasPrice, gasLimit uint64, zptId string, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type regIDWithPublicKey struct {
		ZptId  string
		PubKey []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"regIDWithPublicKey",
		[]interface{}{
			&regIDWithPublicKey{
				ZptId:  zptId,
				PubKey: keypair.SerializePublicKey(pubKey),
			},
		},
	)
}

func (this *ZptId) RegIDWithPublicKey(gasPrice, gasLimit uint64, signer *account.Account, zptId string, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewRegIDWithPublicKeyTransaction(gasPrice, gasLimit, zptId, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewRegIDWithAttributesTransaction(gasPrice, gasLimit uint64, zptId string, pubKey keypair.PublicKey, attributes []*account.DDOAttribute) (*types.MutableTransaction, error) {
	type regIDWithAttribute struct {
		ZptId      string
		PubKey     []byte
		Attributes []*account.DDOAttribute
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"regIDWithAttributes",
		[]interface{}{
			&regIDWithAttribute{
				ZptId:      zptId,
				PubKey:     keypair.SerializePublicKey(pubKey),
				Attributes: attributes,
			},
		},
	)
}

func (this *ZptId) RegIDWithAttributes(gasPrice, gasLimit uint64, signer *account.Account, zptId string, controller *account.Controller, attributes []*account.DDOAttribute) (common.Uint256, error) {
	tx, err := this.NewRegIDWithAttributesTransaction(gasPrice, gasLimit, zptId, controller.PublicKey, attributes)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) GetDDO(zptId string) (*account.DDO, error) {
	result, err := this.native.PreExecInvokeNativeContract(
		ZPT_ID_CONTRACT_ADDRESS,
		ZPT_ID_CONTRACT_VERSION,
		"getDDO",
		[]interface{}{zptId},
	)
	if err != nil {
		return nil, err
	}
	data, err := result.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(data)
	keyData, err := serialization.ReadVarBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("key ReadVarBytes error:%s", err)
	}
	owners, err := this.getPublicKeys(zptId, keyData)
	if err != nil {
		return nil, fmt.Errorf("getPublicKeys error:%s", err)
	}
	attrData, err := serialization.ReadVarBytes(buf)
	attrs, err := this.getAttributes(zptId, attrData)
	if err != nil {
		return nil, fmt.Errorf("getAttributes error:%s", err)
	}
	recoveryData, err := serialization.ReadVarBytes(buf)
	if err != nil {
		return nil, fmt.Errorf("recovery ReadVarBytes error:%s", err)
	}
	var addr string
	if len(recoveryData) != 0 {
		address, err := common.AddressParseFromBytes(recoveryData)
		if err != nil {
			return nil, fmt.Errorf("AddressParseFromBytes error:%s", err)
		}
		addr = address.ToBase58()
	}

	ddo := &account.DDO{
		ZptId:      zptId,
		Owners:     owners,
		Attributes: attrs,
		Recovery:   addr,
	}
	return ddo, nil
}

func (this *ZptId) NewAddKeyTransaction(gasPrice, gasLimit uint64, zptId string, newPubKey, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addKey struct {
		ZptId     string
		NewPubKey []byte
		PubKey    []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"addKey",
		[]interface{}{
			&addKey{
				ZptId:     zptId,
				NewPubKey: keypair.SerializePublicKey(newPubKey),
				PubKey:    keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *ZptId) AddKey(gasPrice, gasLimit uint64, zptId string, signer *account.Account, newPubKey keypair.PublicKey, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewAddKeyTransaction(gasPrice, gasLimit, zptId, newPubKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewRevokeKeyTransaction(gasPrice, gasLimit uint64, zptId string, removedPubKey, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type removeKey struct {
		ZptId      string
		RemovedKey []byte
		PubKey     []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"removeKey",
		[]interface{}{
			&removeKey{
				ZptId:      zptId,
				RemovedKey: keypair.SerializePublicKey(removedPubKey),
				PubKey:     keypair.SerializePublicKey(pubKey),
			},
		},
	)
}

func (this *ZptId) RevokeKey(gasPrice, gasLimit uint64, zptId string, signer *account.Account, removedPubKey keypair.PublicKey, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewRevokeKeyTransaction(gasPrice, gasLimit, zptId, removedPubKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewSetRecoveryTransaction(gasPrice, gasLimit uint64, zptId string, recovery common.Address, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addRecovery struct {
		ZptId    string
		Recovery common.Address
		Pubkey   []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"addRecovery",
		[]interface{}{
			&addRecovery{
				ZptId:    zptId,
				Recovery: recovery,
				Pubkey:   keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *ZptId) SetRecovery(gasPrice, gasLimit uint64, signer *account.Account, zptId string, recovery common.Address, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewSetRecoveryTransaction(gasPrice, gasLimit, zptId, recovery, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewChangeRecoveryTransaction(gasPrice, gasLimit uint64, zptId string, newRecovery, oldRecovery common.Address) (*types.MutableTransaction, error) {
	type changeRecovery struct {
		ZptId       string
		NewRecovery common.Address
		OldRecovery common.Address
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"changeRecovery",
		[]interface{}{
			&changeRecovery{
				ZptId:       zptId,
				NewRecovery: newRecovery,
				OldRecovery: oldRecovery,
			},
		})
}

func (this *ZptId) ChangeRecovery(gasPrice, gasLimit uint64, signer *account.Account, zptId string, newRecovery, oldRecovery common.Address, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewChangeRecoveryTransaction(gasPrice, gasLimit, zptId, newRecovery, oldRecovery)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewAddAttributesTransaction(gasPrice, gasLimit uint64, zptId string, attributes []*account.DDOAttribute, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type addAttributes struct {
		ZptId      string
		Attributes []*account.DDOAttribute
		PubKey     []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"addAttributes",
		[]interface{}{
			&addAttributes{
				ZptId:      zptId,
				Attributes: attributes,
				PubKey:     keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *ZptId) AddAttributes(gasPrice, gasLimit uint64, signer *account.Account, zptId string, attributes []*account.DDOAttribute, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewAddAttributesTransaction(gasPrice, gasLimit, zptId, attributes, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}

	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) NewRemoveAttributeTransaction(gasPrice, gasLimit uint64, zptId string, key []byte, pubKey keypair.PublicKey) (*types.MutableTransaction, error) {
	type removeAttribute struct {
		ZptId  string
		Key    []byte
		PubKey []byte
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"removeAttribute",
		[]interface{}{
			&removeAttribute{
				ZptId:  zptId,
				Key:    key,
				PubKey: keypair.SerializePublicKey(pubKey),
			},
		})
}

func (this *ZptId) RemoveAttribute(gasPrice, gasLimit uint64, signer *account.Account, zptId string, removeKey []byte, controller *account.Controller) (common.Uint256, error) {
	tx, err := this.NewRemoveAttributeTransaction(gasPrice, gasLimit, zptId, removeKey, controller.PublicKey)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return common.UINT256_EMPTY, err
	}

	return this.zptSdk.SendTransaction(tx)
}

func (this *ZptId) GetAttributes(zptId string) ([]*account.DDOAttribute, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_ID_CONTRACT_ADDRESS,
		ZPT_ID_CONTRACT_VERSION,
		"getAttributes",
		[]interface{}{zptId})
	if err != nil {
		return nil, err
	}
	data, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, fmt.Errorf("ToByteArray error:%s", err)
	}
	return this.getAttributes(zptId, data)
}

func (this *ZptId) getAttributes(zptId string, data []byte) ([]*account.DDOAttribute, error) {
	buf := bytes.NewBuffer(data)
	attributes := make([]*account.DDOAttribute, 0)
	for {
		if buf.Len() == 0 {
			break
		}
		key, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("key ReadVarBytes error:%s", err)
		}
		valueType, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("value type ReadVarBytes error:%s", err)
		}
		value, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("value ReadVarBytes error:%s", err)
		}
		attributes = append(attributes, &account.DDOAttribute{
			Key:       key,
			Value:     value,
			ValueType: valueType,
		})
	}
	//reverse
	for i, j := 0, len(attributes)-1; i < j; i, j = i+1, j-1 {
		attributes[i], attributes[j] = attributes[j], attributes[i]
	}
	return attributes, nil
}

func (this *ZptId) VerifySignature(zptId string, keyIndex int, controller *account.Controller) (bool, error) {
	tx, err := this.native.NewNativeInvokeTransaction(
		0, 0,
		ZPT_ID_CONTRACT_VERSION,
		ZPT_ID_CONTRACT_ADDRESS,
		"verifySignature",
		[]interface{}{zptId, keyIndex})
	if err != nil {
		return false, err
	}
	err = this.zptSdk.SignToTransaction(tx, controller)
	if err != nil {
		return false, err
	}
	preResult, err := this.zptSdk.PreExecTransaction(tx)
	if err != nil {
		return false, err
	}
	return preResult.Result.ToBool()
}

func (this *ZptId) GetPublicKeys(zptId string) ([]*account.DDOOwner, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_ID_CONTRACT_ADDRESS,
		ZPT_ID_CONTRACT_VERSION,
		"getPublicKeys",
		[]interface{}{
			zptId,
		})
	if err != nil {
		return nil, err
	}
	data, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	return this.getPublicKeys(zptId, data)
}

func (this *ZptId) getPublicKeys(zptId string, data []byte) ([]*account.DDOOwner, error) {
	buf := bytes.NewBuffer(data)
	owners := make([]*account.DDOOwner, 0)
	for {
		if buf.Len() == 0 {
			break
		}
		index, err := serialization.ReadUint32(buf)
		if err != nil {
			return nil, fmt.Errorf("index ReadUint32 error:%s", err)
		}
		pubKeyId := fmt.Sprintf("%s#keys-%d", zptId, index)
		pkData, err := serialization.ReadVarBytes(buf)
		if err != nil {
			return nil, fmt.Errorf("PubKey Idenx:%d ReadVarBytes error:%s", index, err)
		}
		pubKey, err := keypair.DeserializePublicKey(pkData)
		if err != nil {
			return nil, fmt.Errorf("DeserializePublicKey Index:%d error:%s", index, err)
		}
		keyType := keypair.GetKeyType(pubKey)
		owner := &account.DDOOwner{
			PubKeyIndex: index,
			PubKeyId:    pubKeyId,
			Type:        account.GetKeyTypeString(keyType),
			Curve:       account.GetCurveName(pkData),
			Value:       hex.EncodeToString(pkData),
		}
		owners = append(owners, owner)
	}
	return owners, nil
}

func (this *ZptId) GetKeyState(zptId string, keyIndex int) (string, error) {
	type keyState struct {
		ZptId    string
		KeyIndex int
	}
	preResult, err := this.native.PreExecInvokeNativeContract(
		ZPT_ID_CONTRACT_ADDRESS,
		ZPT_ID_CONTRACT_VERSION,
		"getKeyState",
		[]interface{}{
			&keyState{
				ZptId:    zptId,
				KeyIndex: keyIndex,
			},
		})
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

type GlobalParam struct {
	zptSdk *sdk.ZeepinSdk
	native *NativeContract
}

func (this *GlobalParam) GetGlobalParams(params []string) (map[string]string, error) {
	preResult, err := this.native.PreExecInvokeNativeContract(
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		global_params.GET_GLOBAL_PARAM_NAME,
		[]interface{}{params})
	if err != nil {
		return nil, err
	}
	results, err := preResult.Result.ToByteArray()
	if err != nil {
		return nil, err
	}
	queryParams := new(global_params.Params)
	err = queryParams.Deserialize(bytes.NewBuffer(results))
	if err != nil {
		return nil, err
	}
	globalParams := make(map[string]string, len(params))
	for _, param := range params {
		index, values := queryParams.GetParam(param)
		if index < 0 {
			continue
		}
		globalParams[param] = values.Value
	}
	return globalParams, nil
}

func (this *GlobalParam) NewSetGlobalParamsTransaction(gasPrice, gasLimit uint64, params map[string]string) (*types.MutableTransaction, error) {
	var globalParams global_params.Params
	for k, v := range params {
		globalParams.SetParam(global_params.Param{Key: k, Value: v})
	}
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.SET_GLOBAL_PARAM_NAME,
		[]interface{}{globalParams})
}

func (this *GlobalParam) SetGlobalParams(gasPrice, gasLimit uint64, signer *account.Account, params map[string]string) (common.Uint256, error) {
	tx, err := this.NewSetGlobalParamsTransaction(gasPrice, gasLimit, params)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewTransferAdminTransaction(gasPrice, gasLimit uint64, newAdmin common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.TRANSFER_ADMIN_NAME,
		[]interface{}{newAdmin})
}

func (this *GlobalParam) TransferAdmin(gasPrice, gasLimit uint64, signer *account.Account, newAdmin common.Address) (common.Uint256, error) {
	tx, err := this.NewTransferAdminTransaction(gasPrice, gasLimit, newAdmin)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewAcceptAdminTransaction(gasPrice, gasLimit uint64, admin common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.ACCEPT_ADMIN_NAME,
		[]interface{}{admin})
}

func (this *GlobalParam) AcceptAdmin(gasPrice, gasLimit uint64, signer *account.Account) (common.Uint256, error) {
	tx, err := this.NewAcceptAdminTransaction(gasPrice, gasLimit, signer.Address)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewSetOperatorTransaction(gasPrice, gasLimit uint64, operator common.Address) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.SET_OPERATOR,
		[]interface{}{operator},
	)
}

func (this *GlobalParam) SetOperator(gasPrice, gasLimit uint64, signer *account.Account, operator common.Address) (common.Uint256, error) {
	tx, err := this.NewSetOperatorTransaction(gasPrice, gasLimit, operator)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *GlobalParam) NewCreateSnapshotTransaction(gasPrice, gasLimit uint64) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		GLOBAL_PARAMS_CONTRACT_VERSION,
		GLOABL_PARAMS_CONTRACT_ADDRESS,
		global_params.CREATE_SNAPSHOT_NAME,
		[]interface{}{},
	)
}

func (this *GlobalParam) CreateSnapshot(gasPrice, gasLimit uint64, signer *account.Account) (common.Uint256, error) {
	tx, err := this.NewCreateSnapshotTransaction(gasPrice, gasLimit)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

type Auth struct {
	zptSdk *sdk.ZeepinSdk
	native *NativeContract
}

func (this *Auth) NewAssignFuncsToRoleTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, adminId, role []byte, funcNames []string, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"assignFuncsToRole",
		[]interface{}{
			contractAddress,
			adminId,
			role,
			funcNames,
			keyIndex,
		})
}

func (this *Auth) AssignFuncsToRole(gasPrice, gasLimit uint64, contractAddress common.Address, signer *account.Account, adminId, role []byte, funcNames []string, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewAssignFuncsToRoleTransaction(gasPrice, gasLimit, contractAddress, adminId, role, funcNames, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Auth) NewDelegateTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, from, to, role []byte, period, level, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"delegate",
		[]interface{}{
			contractAddress,
			from,
			to,
			role,
			period,
			level,
			keyIndex,
		})
}

func (this *Auth) Delegate(gasPrice, gasLimit uint64, signer *account.Account, contractAddress common.Address, from, to, role []byte, period, level, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewDelegateTransaction(gasPrice, gasLimit, contractAddress, from, to, role, period, level, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Auth) NewWithdrawTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, initiator, delegate, role []byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"withdraw",
		[]interface{}{
			contractAddress,
			initiator,
			delegate,
			role,
			keyIndex,
		})
}

func (this *Auth) Withdraw(gasPrice, gasLimit uint64, signer *account.Account, contractAddress common.Address, initiator, delegate, role []byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewWithdrawTransaction(gasPrice, gasLimit, contractAddress, initiator, delegate, role, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Auth) NewAssignOntIDsToRoleTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, admontId, role []byte, persons [][]byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"assignOntIDsToRole",
		[]interface{}{
			contractAddress,
			admontId,
			role,
			persons,
			keyIndex,
		})
}

func (this *Auth) AssignOntIDsToRole(gasPrice, gasLimit uint64, signer *account.Account, contractAddress common.Address, admontId, role []byte, persons [][]byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewAssignOntIDsToRoleTransaction(gasPrice, gasLimit, contractAddress, admontId, role, persons, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Auth) NewTransferTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, newAdminId []byte, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"transfer",
		[]interface{}{
			contractAddress,
			newAdminId,
			keyIndex,
		})
}

func (this *Auth) Transfer(gasPrice, gasLimit uint64, signer *account.Account, contractAddress common.Address, newAdminId []byte, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewTransferTransaction(gasPrice, gasLimit, contractAddress, newAdminId, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *Auth) NewVerifyTokenTransaction(gasPrice, gasLimit uint64, contractAddress common.Address, caller []byte, funcName string, keyIndex int) (*types.MutableTransaction, error) {
	return this.native.NewNativeInvokeTransaction(
		gasPrice,
		gasLimit,
		AUTH_CONTRACT_VERSION,
		AUTH_CONTRACT_ADDRESS,
		"verifyToken",
		[]interface{}{
			contractAddress,
			caller,
			funcName,
			keyIndex,
		})
}

func (this *Auth) VerifyToken(gasPrice, gasLimit uint64, signer *account.Account, contractAddress common.Address, caller []byte, funcName string, keyIndex int) (common.Uint256, error) {
	tx, err := this.NewVerifyTokenTransaction(gasPrice, gasLimit, contractAddress, caller, funcName, keyIndex)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}
