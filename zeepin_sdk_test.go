package Zeepin_Go_Sdk

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/tyler-smith/go-bip39"
	"github.com/zeepin/Zeepin-Go-Sdk/utils"
	"github.com/zeepin/ZeepinChain-Crypto/signature"
	"github.com/zeepin/ZeepinChain/common"
	"github.com/zeepin/ZeepinChain/core/payload"
	"github.com/zeepin/ZeepinChain/core/validation"
	"github.com/zeepin/ZeepinChain/smartcontract/event"
)

var (
	TestZptSdk   *ZeepinSdk
	TestWallet   *Wallet
	TestPasswd   = []byte("11")
	TestDefAcc   *Account
	TestGasPrice = uint64(1)
	TestGasLimit = uint64(20000)
)

func TestZeepinSdk_TrabsferFrom(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	payloadHex := "00c66b1421ab6ece5c9e44fa5e35261ef42cc6bc31d98e9c6a7cc814c1d2d106f9d2276b383958973b9fca8e4f48cc966a7cc80400e1f5056a7cc86c51c1087472616e736665721400000000000000000000000000000000000000020068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	payloadBytes, err := common.HexToBytes(payloadHex)
	assert.Nil(t, err)
	res, err := ParsePayload(payloadBytes)
	assert.Nil(t, err)
	fmt.Println("res:", res)

	//java sdk,  transferFrom
	//amount =100
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc801646a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	//amount =10
	//payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc85a6a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//amount = 1000000000
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc8149018fbdfe16d5b1054165ab892b0e040919bd1ca6a7cc8143e7c40c2a2a98e3f95adace19b12ef4a1d7a35066a7cc80400ca9a3b6a7cc86c0c7472616e7366657246726f6d1400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//java sdk, transfer
	//amount = 100
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc801646a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	//amount = 10
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc85a6a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"
	//amount = 1000000000
	payloadHex = "00c66b14d2c124dd088190f709b684e0bc676d70c41b37766a7cc814d2c124dd088190f709b684e0bc676d70c41b37766a7cc80400ca9a3b6a7cc86c51c1087472616e736665721400000000000000000000000000000000000000010068164f6e746f6c6f67792e4e61746976652e496e766f6b65"

	payloadBytes, err = common.HexToBytes(payloadHex)
	assert.Nil(t, err)
	res, err = ParsePayload(payloadBytes)
	assert.Nil(t, err)
	fmt.Println("res:", res)
}

//transferFrom
func TestZeepinSdk_ParseNativeTxPayload2(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	var err error
	assert.Nil(t, err)
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	pri2, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb8cf")
	assert.Nil(t, err)

	pri3, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb9cf")
	assert.Nil(t, err)
	acc, err = NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	acc2, err := NewAccountFromPrivateKey(pri2, signature.SHA256withECDSA)

	acc3, err := NewAccountFromPrivateKey(pri3, signature.SHA256withECDSA)
	amount := 1000000000
	txFrom, err := TestZptSdk.Native.Zpt.NewTransferFromTransaction(500, 20000, acc.Address, acc2.Address, acc3.Address, uint64(amount))
	assert.Nil(t, err)
	tx, err := txFrom.IntoImmutable()
	assert.Nil(t, err)
	invokeCode, ok := tx.Payload.(*payload.InvokeCode)
	assert.True(t, ok)
	code := invokeCode.Code
	res, err := ParsePayload(code)
	assert.Nil(t, err)
	assert.Equal(t, acc.Address.ToBase58(), res["sender"].(string))
	assert.Equal(t, acc2.Address.ToBase58(), res["from"].(string))
	assert.Equal(t, uint64(amount), res["amount"].(uint64))
	assert.Equal(t, "transferFrom", res["functionName"].(string))
	fmt.Println("res:", res)
}
func TestZeepinSdk_ParseNativeTxPayload(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	var err error
	assert.Nil(t, err)
	pri, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb7cf")
	acc, err := NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	pri2, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb8cf")
	assert.Nil(t, err)

	pri3, err := common.HexToBytes("75de8489fcb2dcaf2ef3cd607feffde18789de7da129b5e97c81e001793cb9cf")
	assert.Nil(t, err)
	acc, err = NewAccountFromPrivateKey(pri, signature.SHA256withECDSA)

	acc2, err := NewAccountFromPrivateKey(pri2, signature.SHA256withECDSA)

	acc3, err := NewAccountFromPrivateKey(pri3, signature.SHA256withECDSA)
	y, _ := common.HexToBytes(acc.Address.ToHexString())

	fmt.Println("acc:", common.ToHexString(common.ToArrayReverse(y)))
	assert.Nil(t, err)

	amount := uint64(1000000000)
	tx, err := TestZptSdk.Native.Zpt.NewTransferTransaction(500, 20000, acc.Address, acc2.Address, amount)
	assert.Nil(t, err)

	tx2, err := tx.IntoImmutable()
	assert.Nil(t, err)
	res, err := ParseNativeTxPayload(tx2.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", res)
	assert.Equal(t, acc.Address.ToBase58(), res["from"].(string))
	assert.Equal(t, acc2.Address.ToBase58(), res["to"].(string))
	assert.Equal(t, amount, res["amount"].(uint64))
	assert.Equal(t, "transfer", res["functionName"].(string))

	transferFrom, err := TestZptSdk.Native.Zpt.NewTransferFromTransaction(500, 20000, acc.Address, acc2.Address, acc3.Address, 10)
	transferFrom2, err := transferFrom.IntoImmutable()
	r, err := ParseNativeTxPayload(transferFrom2.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", r)
	assert.Equal(t, r["sender"], acc.Address.ToBase58())
	assert.Equal(t, r["from"], acc2.Address.ToBase58())
	assert.Equal(t, r["to"], acc3.Address.ToBase58())
	assert.Equal(t, r["amount"], uint64(10))

	galaTransfer, err := TestZptSdk.Native.Gala.NewTransferTransaction(uint64(500), uint64(20000), acc.Address, acc2.Address, 100000000)
	assert.Nil(t, err)
	galaTx, err := galaTransfer.IntoImmutable()
	assert.Nil(t, err)
	res, err = ParseNativeTxPayload(galaTx.ToArray())
	assert.Nil(t, err)
	fmt.Println("res:", res)
}

func TestZeepinSdk_GenerateMnemonicCodesStr2(t *testing.T) {
	mnemonic := make(map[string]bool)
	TestZptSdk = NewZeepinSdk()
	for i := 0; i < 100000; i++ {
		mnemonicStr, err := TestZptSdk.GenerateMnemonicCodesStr()
		assert.Nil(t, err)
		if mnemonic[mnemonicStr] == true {
			panic("there is the same mnemonicStr ")
		} else {
			mnemonic[mnemonicStr] = true
		}
	}
}

func TestZeepinSdk_GenerateMnemonicCodesStr(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	for i := 0; i < 1000; i++ {
		mnemonic, err := TestZptSdk.GenerateMnemonicCodesStr()
		assert.Nil(t, err)
		private, err := TestZptSdk.GetPrivateKeyFromMnemonicCodesStrBip44(mnemonic, 0)
		assert.Nil(t, err)
		acc, err := NewAccountFromPrivateKey(private, signature.SHA256withECDSA)
		assert.Nil(t, err)
		si, err := signature.Sign(acc.SigScheme, acc.PrivateKey, []byte("test"), nil)
		boo := signature.Verify(acc.PublicKey, []byte("test"), si)
		assert.True(t, boo)

		tx, err := TestZptSdk.Native.Zpt.NewTransferTransaction(0, 0, acc.Address, acc.Address, 10)
		assert.Nil(t, err)
		TestZptSdk.SignToTransaction(tx, acc)
		tx2, err := tx.IntoImmutable()
		assert.Nil(t, err)
		res := validation.VerifyTransaction(tx2)
		assert.Equal(t, "not an error", res.Error())
	}
}

func TestGenerateMemory(t *testing.T) {
	expectedPrivateKey := []string{"915f5df65c75afe3293ed613970a1661b0b28d0cb711f21c489d8785977df0cd", "dbf1090889ba8b19aa01fa31c8b1ce29828bd2fa664afd95cc62e6055b74e112",
		"1487a8e53e4f4e2e1991781bcd14b3d334d3b2965cb48c976b234da29d7cf242", "79f85da015f079469c6e04aa0fc23523187d0f72c29450073d858ddeed272617"}
	entropy, _ := bip39.NewEntropy(128)
	mnemonic, _ := bip39.NewMnemonic(entropy)
	mnemonic = "ecology cricket napkin scrap board purpose picnic toe bean heart coast retire"
	TestZptSdk = NewZeepinSdk()
	for i := 0; i < len(expectedPrivateKey); i++ {
		privk, err := TestZptSdk.GetPrivateKeyFromMnemonicCodesStrBip44(mnemonic, uint32(i))
		assert.Nil(t, err)
		assert.Equal(t, expectedPrivateKey[i], common.ToHexString(privk))
	}
}

func TestZeepinSdk_CreateWallet(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	wal, err := TestZptSdk.CreateWallet("./wallet3.dat")
	assert.Nil(t, err)
	_, err = wal.NewDefaultSettingAccount([]byte("11"))
	assert.Nil(t, err)
	wal.Save()
}

func TestNewZeepinSdk(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	TestWallet, _ = TestZptSdk.OpenWallet("./wallet.dat")
	event := &event.NotifyEventInfo{
		ContractAddress: common.ADDRESS_EMPTY,
		States:          []interface{}{"transfer", "Abc3UVbyL1kxd9sK6N9hzAT2u91ftbpoXT", "AFmseVrdL9f9oyCzZefL9tG6UbviEH9ugK", uint64(10000000)},
	}
	e, err := TestZptSdk.ParseNaitveTransferEvent(event)
	assert.Nil(t, err)
	fmt.Println(e)
}

func TestZeepinSdk_GetTxData(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	TestWallet, _ = TestZptSdk.OpenWallet("./wallet2.dat")
	acc, _ := TestWallet.GetAccountByAddress("ZaD9GckZt7cPVRL8iTvVB1iYzAUYhvNd8x", TestPasswd)
	tx, _ := TestZptSdk.Native.Zpt.NewTransferTransaction(500, 10000, acc.Address, acc.Address, 100)

	TestZptSdk.SignToTransaction(tx, acc)
	tx2, _ := tx.IntoImmutable()
	var buffer bytes.Buffer
	tx2.Serialize(&buffer)
	txData := hex.EncodeToString(buffer.Bytes())
	tx3, _ := TestZptSdk.GetMutableTx(txData)
	assert.Equal(t, tx, tx3)
}

func Init() {
	TestZptSdk = NewZeepinSdk()
	TestZptSdk.NewRpcClient().SetAddress("http://localhost:20336")

	var err error
	var wallet *Wallet
	if !common.FileExisted("./wallet.dat") {
		wallet, err = TestZptSdk.CreateWallet("./wallet.dat")
		if err != nil {
			fmt.Println("[CreateWallet] error:", err)
			return
		}
	} else {
		wallet, err = TestZptSdk.OpenWallet("./wallet.dat")
		if err != nil {
			fmt.Println("[CreateWallet] error:", err)
			return
		}
	}
	_, err = wallet.NewDefaultSettingAccount(TestPasswd)
	if err != nil {
		fmt.Println("")
		return
	}
	wallet.Save()
	TestWallet, err = TestZptSdk.OpenWallet("./wallet.dat")
	if err != nil {
		fmt.Printf("account.Open error:%s\n", err)
		return
	}
	TestDefAcc, err = TestWallet.GetDefaultAccount(TestPasswd)
	if err != nil {
		fmt.Printf("GetDefaultAccount error:%s\n", err)
		return
	}

	ws := TestZptSdk.NewWebSocketClient()
	err = ws.Connect("ws://localhost:20335")
	if err != nil {
		fmt.Printf("Connect ws error:%s", err)
		return
	}
}

func TestZpt_Transfer(t *testing.T) {
	TestZptSdk = NewZeepinSdk()
	TestWallet, _ = TestZptSdk.OpenWallet("./walle2.dat")
	toaddr, _ := utils.AddressFromHexString("ZPdTk4zzhvhDA5sbgnK1bJ9qkYegy2QvrH")
	TestZptSdk.NewRpcClient().SetAddress("http://localhost:20336")

	acc, _ := TestWallet.GetAccountByAddress("ZaD9GckZt7cPVRL8iTvVB1iYzAUYhvNd8x", TestPasswd)

	txHash, err := TestZptSdk.Native.Zpt.Transfer(TestGasPrice, TestGasLimit, acc, toaddr, 1)
	if err != nil {
		t.Errorf("NewTransferTransaction error:%s", err)
		return
	}
	TestZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	evts, err := TestZptSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("GetSmartCZptractEvent error:%s", err)
		return
	}
	fmt.Printf("TxHash:%s\n", txHash.ToHexString())
	fmt.Printf("State:%d\n", evts.State)
	fmt.Printf("GasConsume:%d\n", evts.GasConsumed)
	for _, notify := range evts.Notify {
		fmt.Printf("CZptractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}
}

func TestGALA_WithDrawGALA(t *testing.T) {
	Init()
	unboundGALA, err := TestZptSdk.Native.Gala.UnboundGala(TestDefAcc.Address.ToBase58())
	if err != nil {
		t.Errorf("UnboundGALA error:%s", err)
		return
	}
	fmt.Printf("Address:%s UnboundGALA:%d\n", TestDefAcc.Address.ToBase58(), unboundGALA)
	_, err = TestZptSdk.Native.Gala.WithdrawGala(0, 20000, TestDefAcc, unboundGALA)
	if err != nil {
		t.Errorf("WithDrawGALA error:%s", err)
		return
	}
	fmt.Printf("Address:%s WithDrawGALA amount:%d success\n", TestDefAcc.Address.ToBase58(), unboundGALA)
}

func TestGlobalParam_GetGlobalParams(t *testing.T) {
	Init()
	gasPrice := "gasPrice"
	params := []string{gasPrice}
	results, err := TestZptSdk.Native.GlobalParams.GetGlobalParams(params)
	if err != nil {
		t.Errorf("GetGlobalParams:%+v error:%s", params, err)
		return
	}
	fmt.Printf("Params:%s Value:%v\n", gasPrice, results[gasPrice])
}

func TestGlobalParam_SetGlobalParams(t *testing.T) {
	Init()
	gasPrice := "gasPrice"
	globalParams, err := TestZptSdk.Native.GlobalParams.GetGlobalParams([]string{gasPrice})
	if err != nil {
		t.Errorf("GetGlobalParams error:%s", err)
		return
	}
	gasPriceValue, err := strconv.Atoi(globalParams[gasPrice])
	if err != nil {
		t.Errorf("Get prama value error:%s", err)
		return
	}
	_, err = TestZptSdk.Native.GlobalParams.SetGlobalParams(TestGasPrice, TestGasLimit, TestDefAcc, map[string]string{gasPrice: strconv.Itoa(gasPriceValue + 1)})
	if err != nil {
		t.Errorf("SetGlobalParams error:%s", err)
		return
	}
	TestZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	globalParams, err = TestZptSdk.Native.GlobalParams.GetGlobalParams([]string{gasPrice})
	if err != nil {
		t.Errorf("GetGlobalParams error:%s", err)
		return
	}
	gasPriceValueAfter, err := strconv.Atoi(globalParams[gasPrice])
	if err != nil {
		t.Errorf("Get prama value error:%s", err)
		return
	}
	fmt.Printf("After set params gasPrice:%d\n", gasPriceValueAfter)
}

func TestWsScribeEvent(t *testing.T) {
	Init()
	wsClient := TestZptSdk.ClientMgr.GetWebSocketClient()
	err := wsClient.SubscribeEvent()
	if err != nil {
		t.Errorf("SubscribeTxHash error:%s", err)
		return
	}
	defer wsClient.UnsubscribeTxHash()

	actionCh := wsClient.GetActionCh()
	timer := time.NewTimer(time.Minute * 3)
	for {
		select {
		case <-timer.C:
			return
		case action := <-actionCh:
			fmt.Printf("Action:%s\n", action.Action)
			fmt.Printf("Result:%s\n", action.Result)
		}
	}
}

func TestWsTransfer(t *testing.T) {
	Init()
	wsClient := TestZptSdk.ClientMgr.GetWebSocketClient()
	TestZptSdk.ClientMgr.SetDefaultClient(wsClient)
	txHash, err := TestZptSdk.Native.Zpt.Transfer(TestGasPrice, TestGasLimit, TestDefAcc, TestDefAcc.Address, 1)
	if err != nil {
		t.Errorf("NewTransferTransaction error:%s", err)
		return
	}
	TestZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	evts, err := TestZptSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("GetSmartCZptractEvent error:%s", err)
		return
	}
	fmt.Printf("TxHash:%s\n", txHash.ToHexString())
	fmt.Printf("State:%d\n", evts.State)
	fmt.Printf("GasConsume:%d\n", evts.GasConsumed)
	for _, notify := range evts.Notify {
		fmt.Printf("CZptractAddress:%s\n", notify.ContractAddress)
		fmt.Printf("States:%+v\n", notify.States)
	}
}
