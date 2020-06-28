package Zeepin_Go_Sdk

import (
	"fmt"
	"github.com/zeepin/ZeepinChain/common"
	"github.com/zeepin/ZeepinChain/core/types"
	"github.com/zeepin/Zeepin-Go-Sdk/utils"
	sdkcom "github.com/zeepin/Zeepin-Go-Sdk/common"
)

type WasmVMContract struct {
	zptSdk       *ZeepinSdk
}

func NewWasmVMContract(zptSdk *ZeepinSdk) *WasmVMContract{
	return  &WasmVMContract{
		zptSdk: zptSdk,
	}
}

func(this *WasmVMContract) NewWasmVMInvokeTransaction(
	method string,
	args []interface{},
	contractAddr string,
	gasPrice,
	gasLimit uint64) (*types.MutableTransaction, error){
	invokeCode, err := utils.BuildWasmVMInvokeCode(method, contractAddr, args)
	if err != nil {
		return nil, err
	}
	return this.zptSdk.NewInvokeTransaction(gasPrice, gasLimit, invokeCode), nil
}

func (this *WasmVMContract) SendWasmTransaction(
	gasPrice,
	gasLimit uint64,
	signer *Account,
	contractAddr string,
	args []interface{},
	method string) (common.Uint256, error) {

	tx, err := this.NewWasmVMInvokeTransaction(method, args, contractAddr, gasPrice, gasLimit)
	if err != nil {
		return common.UINT256_EMPTY, fmt.Errorf("SendWasmTransaction error:%s", err)
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	fmt.Println(tx)
	if err != nil {
		return common.UINT256_EMPTY, err
	}
	return this.zptSdk.SendTransaction(tx)
}

func (this *WasmVMContract) MakeWasmTransaction(
	gasPrice,
	gasLimit uint64,
	contractAddr string,
	args []interface{},
	method string) (*types.MutableTransaction, error) {

	tx, err := this.NewWasmVMInvokeTransaction(method, args, contractAddr, gasPrice, gasLimit)
	if err != nil {
		return nil, fmt.Errorf("MakeWasmTransaction error:%s", err)
	}

	return tx, nil
}

func (this *WasmVMContract) PreExecInvokeWasmVMContract(
	gasPrice,
	gasLimit uint64,
	signer *Account,
	contractAddr string,
	args []interface{},
	method string) (*sdkcom.PreExecResult, error) {

	tx, err := this.NewWasmVMInvokeTransaction(method, args, contractAddr, gasPrice, gasLimit)
	if err != nil {
		return nil, fmt.Errorf("PreExecInvokeWasmVMContract error:%s", err)
	}
	err = this.zptSdk.SignToTransaction(tx, signer)
	fmt.Println(tx)
	if err != nil {
		return nil, err
	}
	return this.zptSdk.PreExecTransaction(tx)
}

