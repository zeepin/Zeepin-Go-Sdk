package test

import (
	"fmt"
	"github.com/zeepin/ZeepinChain-Crypto/keypair"
	"github.com/zeepin/ZeepinChain/common"
	sdk "github.com/zeepin/Zeepin-Go-Sdk"
	//"github.com/zeepin/Zeepin-Go-Sdk/utils"
	"time"
)


func main(){
	fmt.Println("hello")
	zptSdk := sdk.NewZeepinSdk();
	zptSdk.NewRpcClient().SetAddress("http://localhost:20336")

	//testCreateAccount(zptSdk, "11")

	//testNewAccountFromWifPrivateKey(zptSdk, "", "11")
	//addr := "ZK4xgvBom4D33F9YAmgg89fJW18iVss3tV"

	//Address, _:= utils.AddressFromBase58(addr)
	//fmt.Println(Address)
	//testBalance(zptSdk, Address)


	//testTransfer(zptSdk)
	//testWithDrawGala(zptSdk, Address)

	testMutiAddr(zptSdk)

	//testWasmTransaction(zptSdk)
	//testGetHash(zptSdk, "a9dc36cdc1f5459532816ff7cff891824392d52e1ea57d899387c376f5ae8a61")
}

// wasm合约调用
func testWasmTransaction(zptSdk *sdk.ZeepinSdk){

	contrctAddr := "0f27a43a74c963e07c0b633aff49ebb269e6d727"

	pri, err := common.HexToBytes("2cf804f021d94c33a3a288d6fc0d74f19854f6ef01de20f3ad8b19166b221d90")
	logError(err)

	acct, err:= zptSdk.NewAccountFromPrivateKey(pri, 1)
	logError(err)

	to := "ZQGwVawooDQs7WpGMNT1tkEYXEnV7QUb2M"
	amount := "10"
	txhash, err := zptSdk.WasmVM.SendWasmTransaction(1, 20000, acct, contrctAddr,
		[]interface{}{acct.Address.ToBase58(), to, amount}, "transfer")

	if err != nil {
		fmt.Println("989898989009-09-9-9-9-9-09-9-9")
	}

	fmt.Println(txhash.ToHexString())
	//zptSdk.WaitForGenerateBlock(6*time.Second, 1)
	//testGetHash(zptSdk, txhash.ToHexString())

}


// zpt,gala交易
func testTransfer(zptSdk *sdk.ZeepinSdk){
	fmt.Println("=-=-=-=-=-=-=-=-=-==--=-===-=-=-")

	pri, err := common.HexToBytes("2cf804f021d94c33a3a288d6fc0d74f19854f6ef01de20f3ad8b19166b221d90")
	logError(err)

	acct, err:= zptSdk.NewAccountFromPrivateKey(pri, 1)
	logError(err)

	pri1, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	logError(err)
	acc1, err := zptSdk.NewAccountFromPrivateKey(pri1, 1)
	logError(err)

	fmt.Println(acc1.Address.ToBase58())

	txHash, err := zptSdk.Native.Gala.Transfer(1, 20000 , acct, acc1.Address, 10)
	logError(err)

	zptSdk.WaitForGenerateBlock(5*time.Second, 1)

	testGetHash(zptSdk, txHash.ToHexString())

}


// 根据交易hash获取交易内容
func testGetHash(zptSdk *sdk.ZeepinSdk, txHash string){

	evts, err := zptSdk.GetSmartContractEvent(txHash)
	logError(err)
	fmt.Printf("TxHash:%s\n", txHash)
	fmt.Printf("State:%d\n", evts.State)
	fmt.Printf("GasConsume:%d\n", evts.GasConsumed)
	for _, notify := range evts.Notify {
		fmt.Printf("CZptractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}

}


// 解绑gala
func testWithDrawGala(zptSdk *sdk.ZeepinSdk, address common.Address){
	// 查询
	unboundGala, err := zptSdk.Native.Gala.UnboundGala(address)
	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(unboundGala)


	pri, err := common.HexToBytes("2cf804f021d94c33a3a288d6fc0d74f19854f6ef01de20f3ad8b19166b221d90")
	logError(err)

	acct, err:= zptSdk.NewAccountFromPrivateKey(pri, 1)
	logError(err)

	// 解绑
	txhash, err := zptSdk.Native.Gala.WithdrawGala(1, 20000, acct, unboundGala)

	if err != nil {
		fmt.Println("error: ", err)
	}
	fmt.Println(txhash.ToHexString())

}

// 多签地址
func testMutiAddr(zptSdk *sdk.ZeepinSdk){
	fmt.Println("testMutiAddr")

	privatekey1 := "49855b16636e70f100cc5f4f42bc20a6535d7414fb8845e7310f8dd065a97221"
	privatekey2 := "1094e90dd7c4fdfd849c14798d725ac351ae0d924b29a279a9ffa77d5737bd96"
	privatekey3 := "bc254cf8d3910bc615ba6bf09d4553846533ce4403bc24f58660ae150a6d64cf"

	pri1, _ := common.HexToBytes(privatekey1)
	pri2, _ := common.HexToBytes(privatekey2)
	pri3, _ := common.HexToBytes(privatekey3)

	acct1, _:= zptSdk.NewAccountFromPrivateKey(pri1, 1)
	acct2, _:= zptSdk.NewAccountFromPrivateKey(pri2, 1)
	acct3, _:= zptSdk.NewAccountFromPrivateKey(pri3, 1)


	pub := []keypair.PublicKey{acct1.PublicKey, acct2.PublicKey, acct3.PublicKey}
	m := 2
	mutiAddr, err := zptSdk.GetMultiAddr(pub, m)
	if err != nil{
		fmt.Println("error: ", err)
	}
	fmt.Println(mutiAddr)


	// transaction

	//from, _:= utils.AddressFromBase58(mutiAddr)
	addr := "ZK4xgvBom4D33F9YAmgg89fJW18iVss3tV"
	//to, _:= utils.AddressFromBase58(addr)
	//tx, err := zptSdk.Native.Gala.NewTransferTransaction(1, 20000, from, to, 20 )
	tx, err := zptSdk.WasmVM.MakeWasmTransaction(1, 20000, "0f27a43a74c963e07c0b633aff49ebb269e6d727", []interface{}{mutiAddr, addr, "10"}, "transfer")
	if err != nil{
		fmt.Println("=-=-=-=-=-=-=-=-=--=-=-=-=-=-=-=-=")
	}
	fmt.Printf("tx: %v \n", tx)

	err = zptSdk.MultiSignToTransaction(tx, uint16(m), pub, acct1 )
	if err != nil{
		fmt.Println("11111111:error")
	}
	err = zptSdk.MultiSignToTransaction(tx, uint16(m), pub, acct2 )
	if err != nil{
		fmt.Println("22222222:error")
	}
	txhash, err := zptSdk.SendTransaction(tx)
	if err != nil{
		fmt.Println("333333333:error")
	}
	fmt.Println(txhash.ToHexString())

}


// 查询zpt，gala余额
func testBalance(zptSdk *sdk.ZeepinSdk, address common.Address){
	fmt.Println(address)

	balance, _ := zptSdk.Native.Zpt.BalanceOf(address)
	fmt.Println(balance)
}


// 创建账号
func testCreateAccount(zptSdk *sdk.ZeepinSdk, password string){

	wal, err := zptSdk.CreateWallet("./wallet3.dat")

	if err != nil{
		fmt.Println("11111: ",err)
	}

	_, err = wal.NewDefaultSettingAccount([]byte(password))
	if err != nil{
		fmt.Println("22222: ",err)
	}
	wal.Save()
}


// 用wif私钥生成钱包
func testNewAccountFromWifPrivateKey(zptSdk *sdk.ZeepinSdk, wif string, password string){
	wal, err := zptSdk.CreateWallet("./wallet1.dat")

	if err != nil{
		fmt.Println("11111: ",err)
	}

	_, err = wal.NewAccountFromWIF([]byte(wif), []byte(password))
	if err != nil{
		fmt.Println("22222: ",err)
	}
	wal.Save()

}

func logError(err error){
	if err != nil{
		fmt.Println("error: ", err)
	}

}




