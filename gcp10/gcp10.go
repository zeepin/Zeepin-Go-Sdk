/*

 */
package gcp10

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"unsafe"

	"github.com/zeepin/ZeepinChain-Crypto/keypair"
	"github.com/zeepin/ZeepinChain/common"
	"github.com/zeepin/ZeepinChain/core/types"
	"github.com/zeepin/Zeepin-Go-Sdk"
	scomm "github.com/zeepin/Zeepin-Go-Sdk/common"
	"github.com/zeepin/Zeepin-Go-Sdk/utils"
	"github.com/zeepin/Zeepin-Go-Sdk/account"
)

type Gcp10 struct {
	ContractAddress string
	sdk             *zeepin_go_sdk.ZeepinSdk
}

func NewGcp10(address string, sdk *zeepin_go_sdk.ZeepinSdk) *Gcp10 {
	return &Gcp10{
		ContractAddress: address,
		sdk:             sdk,
	}
}

func (this *Gcp10) Name(account *account.Account) (string, error) {
	preResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(1, 20000, account, this.ContractAddress,
		[]interface{}{},"name")
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Gcp10) Symbol(account *account.Account) (string, error) {
	preResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(1, 20000, account, this.ContractAddress,
		[]interface{}{},"symbol")
	if err != nil {
		return "", err
	}
	return preResult.Result.ToString()
}

func (this *Gcp10) Decimals(account *account.Account) (*big.Int, error) {
	preResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(1, 20000, account, this.ContractAddress,
		[]interface{}{},"decimals")
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Gcp10) TotalSupply(account *account.Account) (*big.Int, error) {
	preResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(1, 20000, account, this.ContractAddress,
		[]interface{}{},"totalSupply")
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Gcp10) BalanceOf(account *account.Account, address string) (*big.Int, error) {
	preResult, err := this.sdk.WasmVM.PreExecInvokeWasmVMContract(1, 20000, account, this.ContractAddress,
		[]interface{}{address},"balanceOf")
	if err != nil {
		return nil, err
	}
	return preResult.Result.ToInteger()
}

func (this *Gcp10) Transfer(from *account.Account, to string, amount string, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	result, err := this.sdk.WasmVM.SendWasmTransaction(1, 20000, from, this.ContractAddress,
		[]interface{}{from.Address.ToBase58(), to, amount},"transfer")
	if err != nil{
		return common.UINT256_EMPTY, err
	}
	return result, nil
}

func (this *Gcp10) MultiSignTransfer(fromAccounts []*account.Account, m int, to string, amount string,
	gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range fromAccounts {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	fromAddr, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.WasmVM.MakeWasmTransaction(1, 20000, this.ContractAddress, []interface{}{fromAddr, to, amount}, "transfer")
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range fromAccounts {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}


func (this *Gcp10) Approve(owner *account.Account, spender string, amount string, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	return this.sdk.WasmVM.SendWasmTransaction(1, 20000, owner, this.ContractAddress,
		[]interface{}{owner.Address.ToBase58(), spender, amount},"approve")
}

func (this *Gcp10) MultiSignApprove(ownerAccounts []*account.Account, m int, spender string,
	amount string, gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range ownerAccounts {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	owner, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.WasmVM.MakeWasmTransaction(1, 20000, this.ContractAddress, []interface{}{owner, spender, amount}, "approve")
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range ownerAccounts {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

func (this *Gcp10) TransferFrom(spender *account.Account, from, to string, amount string, gasPrice,
	gasLimit uint64) (common.Uint256, error) {
	return this.sdk.WasmVM.SendWasmTransaction(1, 20000, spender, this.ContractAddress,
		[]interface{}{spender.Address.ToBase58(), from, to, amount},"transferFrom")
}

func (this *Gcp10) MultiSignTransferFrom(spenders []*account.Account, m int, from, to string,
	amount string, gasPrice, gasLimit uint64) (common.Uint256, error) {
	pubKeys := make([]keypair.PublicKey, 0)
	for _, acc := range spenders {
		pubKeys = append(pubKeys, acc.PublicKey)
	}
	spender, err := types.AddressFromMultiPubKeys(pubKeys, m)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("generate multi-sign address failed, err: %s", err)
	}
	mutableTx, err := this.sdk.WasmVM.MakeWasmTransaction(1, 20000, this.ContractAddress, []interface{}{spender, from, to, amount}, "transferFrom")
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("construct tx failed, err: %s", err)
	}
	for _, signer := range spenders {
		err = this.sdk.MultiSignToTransaction(mutableTx, uint16(m), pubKeys, signer)
		if err != nil {
			return common.UINT256_EMPTY, fmt.Errorf("multi sign failed, err: %s", err)
		}
	}
	return this.sdk.SendTransaction(mutableTx)
}

func (this *Gcp10) FetchTxTransferEvent(hash string) ([]*Gcp10TransferEvent, error) {
	contractEvt, err := this.sdk.GetSmartContractEvent(hash)
	if err != nil {
		return nil, err
	}
	return this.parseTransferEvent(contractEvt), nil
}

// TODO: fetch approve event

func (this *Gcp10) FetchBlockTransferEvent(height uint32) ([]*Gcp10TransferEvent, error) {
	contractEvt, err := this.sdk.GetSmartContractEventByBlock(height)
	if err != nil {
		return nil, err
	}
	result := make([]*Gcp10TransferEvent, 0)
	for _, evt := range contractEvt {
		result = append(result, this.parseTransferEvent(evt)...)
	}
	return result, nil
}

func (this *Gcp10) parseTransferEvent(contractEvt *scomm.SmartContactEvent) []*Gcp10TransferEvent {
	result := make([]*Gcp10TransferEvent, 0)
	for _, notify := range contractEvt.Notify {
		addr, _ := utils.AddressFromHexString(notify.ContractAddress)
		if addr.ToBase58() == this.ContractAddress {
			selfEvt, err := parseGcp10TransferEvent(notify)
			if err == nil {
				result = append(result, selfEvt)
			}
		}
	}
	return result
}

func parseGcp10TransferEvent(notify *scomm.NotifyEventInfo) (*Gcp10TransferEvent, error) {
	state, ok := notify.States.([]interface{})
	if !ok {
		return nil, fmt.Errorf("state.States is not []interface")
	}
	if len(state) != 4 {
		return nil, fmt.Errorf("state length is not 4")
	}
	eventName, ok := state[0].(string)
	if !ok {
		return nil, fmt.Errorf("state.States[0] is not string")
	}
	from, ok := state[1].(string)
	if !ok {
		return nil, fmt.Errorf("state[1] is not string")
	}
	to, ok := state[2].(string)
	if !ok {
		return nil, fmt.Errorf("state[2] is not string")
	}
	amount, ok := state[3].(string)
	if !ok {
		return nil, fmt.Errorf("state[3] is not uint64")
	}
	evt, err := hex.DecodeString(eventName)
	if err != nil {
		return nil, fmt.Errorf("decode event name failed, err: %s", err)
	}
	fromAddr, err := utils.AddressFromHexString(from)
	if err != nil {
		return nil, fmt.Errorf("decode from failed, err: %s", err)
	}
	toAddr, err := utils.AddressFromHexString(to)
	if err != nil {
		return nil, fmt.Errorf("decode to failed, err: %s", err)
	}
	value, err := hex.DecodeString(amount)
	if err != nil {
		return nil, fmt.Errorf("decode value failed, err: %s", err)
	}
	return &Gcp10TransferEvent{
		Name:   string(evt),
		From:   fromAddr.ToBase58(),
		To:     toAddr.ToBase58(),
		Amount: (*string)(unsafe.Pointer(&value)),
	}, nil
}
