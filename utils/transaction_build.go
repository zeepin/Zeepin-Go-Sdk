package utils

import (
	"bytes"
	"encoding/json"
	"github.com/zeepin/ZeepinChain/common"
	util "github.com/zeepin/Zeepin-Go-Sdk/common"
)

const(
	TYPE     string   =   "type"
	VALUE    string   =   "value"
	PARAM    string   =   "Params"
)



func BuildWasmVMInvokeCode(method ,contractAddr string, params []interface{}) ([]byte, error) {

	conaddr, err:= BuildWasmVMContractAddrParam(contractAddr)
	if err != nil{
		return nil, err
	}

	args, err:= BuildWasmVMParam(params)
	if err != nil{
		return nil, err
	}
	builder := new(bytes.Buffer)
	builder.WriteString("1")
	builder.Write(conaddr)
	builder.Write([]byte{util.IntToByte(len(method))})
	builder.WriteString(method)
	builder.Write([]byte{util.IntToByte(len(args))})
	builder.WriteString(args)

	return builder.Bytes(), nil
}


//buildWasmVMParamInter build wasmvm invoke param code
func BuildWasmVMParam(smartContractParams []interface{}) (args string ,err error) {
	smartContractParamsLen := len(smartContractParams)
	paramArray := make([]interface{}, smartContractParamsLen)
	for i :=  0; i < smartContractParamsLen; i++ {
		paramsMap := make(map[string]interface{})
		switch v := smartContractParams[i].(type) {
		case byte:
			paramsMap[TYPE] = "byte"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case int:
			paramsMap[TYPE] = "int"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case uint:
			paramsMap[TYPE] = "uint"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case int32:
			paramsMap[TYPE] = "int32"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case uint32:
			paramsMap[TYPE] = "uint32"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case int64:
			paramsMap[TYPE] = "int64"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case uint64:
			paramsMap[TYPE] = "uint64"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case string:
			paramsMap[TYPE] = "string"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		case []byte:
			paramsMap[TYPE] = "byte"
			paramsMap[VALUE] = v
			paramArray[i] = paramsMap
		default:
			return "" , err
		}
	}

	param := make(map[string]interface{})
	param[PARAM] = paramArray
	paramsJson, _ := json.Marshal(param)
	args = string(paramsJson)
	return args, nil
}


// hextobyte for contractAddr params, and reverse output
func BuildWasmVMContractAddrParam(contractAddr string) ([]byte, error){

	conaddr,err := common.HexToBytes(contractAddr)
	if err != nil{
		return nil, err
	}
	conaddr = util.Reverse(conaddr)
	return conaddr, nil
}