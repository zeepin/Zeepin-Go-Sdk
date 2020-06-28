package Zeepin_Go_Sdk

import (
	"encoding/hex"
	"fmt"
	"github.com/zeepin/ZeepinChain-Crypto/keypair"
	"testing"
	"time"
)

var (
	testZptSdk   *ZeepinSdk
	testWallet   *Wallet
	testPasswd   = []byte("11")
	testDefAcc   *Account
	testGasPrice = uint64(1)
	testGasLimit = uint64(20000)
)

func TestZptId_RegIDWithPublicKey(t *testing.T) {
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey GetControllerByIndex error:%s", err)
		return
	}
	txHash, err := testZptSdk.Native.ZptId.RegIDWithPublicKey(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey RegIDWithPublicKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	event, err := testZptSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey GetSmartContractEvent error:%s", err)
		return
	}
	fmt.Printf("TestZptId_RegIDWithPublicKey Event: %+v\n", event)

	ddo, err := testZptSdk.Native.ZptId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey GetDDO error:%s", err)
		return
	}
	fmt.Printf("TestZptId_RegIDWithPublicKey DDO:%+v\n", ddo)
}

func TestZptId_RegIDWithAttributes(t *testing.T) {
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey GetControllerByIndex error:%s", err)
		return
	}
	attributes := make([]*DDOAttribute, 0)
	attr1 := &DDOAttribute{
		Key:       []byte("Hello"),
		Value:     []byte("World"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr1)
	attr2 := &DDOAttribute{
		Key:       []byte("Foo"),
		Value:     []byte("Bar"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr2)
	_, err = testZptSdk.Native.ZptId.RegIDWithAttributes(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefController, attributes)
	if err != nil {
		t.Errorf("TestZptId_RegIDWithPublicKey RegIDWithAttributes error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)

	ddo, err := testZptSdk.Native.ZptId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("GetDDO error:%s", err)
		return
	}

	owners := ddo.Owners
	if owners[0].Value != hex.EncodeToString(keypair.SerializePublicKey(testDefController.GetPublicKey())) {
		t.Errorf("TestZptId_RegIDWithPublicKey pubkey %s != %s", owners[0].Value, hex.EncodeToString(keypair.SerializePublicKey(testDefController.GetPublicKey())))
		return
	}
	attrs := ddo.Attributes
	if len(attributes) != len(attrs) {
		t.Errorf("TestZptId_RegIDWithPublicKey attribute size %d != %d", len(attrs), len(attributes))
		return
	}
	fmt.Printf("Owner:%+v\n", owners[0])
	if string(attr1.Key) != string(attrs[0].Key) ||
		string(attr1.Value) != string(attrs[0].Value) ||
		string(attr1.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestZptId_RegIDWithPublicKey Attribute %s != %s", attrs[0], attr1)
	}
	if string(attr2.Key) != string(attrs[1].Key) ||
		string(attr2.Value) != string(attrs[1].Value) ||
		string(attr2.ValueType) != string(attrs[1].ValueType) {
		t.Errorf("TestZptId_RegIDWithPublicKey Attribute %s != %s", attrs[1], attr2)
	}
}

func TestZptId_Key(t *testing.T) {
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Key NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Key GetControllerByIndex error:%s", err)
		return
	}
	_, err = testZptSdk.Native.ZptId.RegIDWithPublicKey(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Key RegIDWithPublicKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)

	controller1, err := testIdentity.NewDefaultSettingController("2", testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Key NewDefaultSettingController error:%s", err)
		return
	}

	_, err = testZptSdk.Native.ZptId.AddKey(testGasPrice, testGasLimit, testIdentity.ID, testDefAcc, controller1.PublicKey, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Key AddKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)

	owners, err := testZptSdk.Native.ZptId.GetPublicKeys(testIdentity.ID)
	if err != nil {
		t.Errorf("TestZptId_Key GetPublicKeys error:%s", err)
		return
	}

	if len(owners) != 2 {
		t.Errorf("TestZptId_Key owner size:%d != 2", len(owners))
		return
	}

	if owners[0].Value != hex.EncodeToString(keypair.SerializePublicKey(testDefController.PublicKey)) {
		t.Errorf("TestZptId_Key owner index:%d pubkey:%s != %s", owners[0].PubKeyIndex, owners[0].Value, hex.EncodeToString(keypair.SerializePublicKey(testDefController.PublicKey)))
		return
	}

	if owners[1].Value != hex.EncodeToString(keypair.SerializePublicKey(controller1.PublicKey)) {
		t.Errorf("TestZptId_Key owner index:%d pubkey:%s != %s", owners[1].PubKeyIndex, owners[1].Value, hex.EncodeToString(keypair.SerializePublicKey(controller1.PublicKey)))
		return
	}

	_, err = testZptSdk.Native.ZptId.RevokeKey(testGasPrice, testGasLimit, testIdentity.ID, testDefAcc, testDefController.PublicKey, controller1)
	if err != nil {
		t.Errorf("TestZptId_Key RevokeKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)

	owners, err = testZptSdk.Native.ZptId.GetPublicKeys(testIdentity.ID)
	if err != nil {
		t.Errorf("TestZptId_Key GetPublicKeys error:%s", err)
		return
	}

	if len(owners) != 1 {
		t.Errorf("TestZptId_Key owner size:%d != 1 after remove", len(owners))
		return
	}

	state, err := testZptSdk.Native.ZptId.GetKeyState(testIdentity.ID, 1)
	if err != nil {
		t.Errorf("TestZptId_Key GetKeyState error:%s", err)
		return
	}

	if state != KEY_STATUS_REVOKE {
		t.Errorf("TestZptId_Key remove key state != %s", KEY_STATUS_REVOKE)
		return
	}

	state, err = testZptSdk.Native.ZptId.GetKeyState(testIdentity.ID, 2)
	if err != nil {
		t.Errorf("TestZptId_Key GetKeyState error:%s", err)
		return
	}
	if state != KEY_STSTUS_IN_USE {
		t.Errorf("TestZptId_Key GetKeyState state != %s", KEY_STSTUS_IN_USE)
		return
	}
}

func TestZptId_Attribute(t *testing.T) {
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Attribute NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Attribute GetControllerByIndex error:%s", err)
		return
	}
	_, err = testZptSdk.Native.ZptId.RegIDWithPublicKey(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Attribute RegIDWithPublicKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)

	attributes := make([]*DDOAttribute, 0)
	attr1 := &DDOAttribute{
		Key:       []byte("Foo"),
		Value:     []byte("Bar"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr1)
	attr2 := &DDOAttribute{
		Key:       []byte("Hello"),
		Value:     []byte("World"),
		ValueType: []byte("string"),
	}
	attributes = append(attributes, attr2)
	_, err = testZptSdk.Native.ZptId.AddAttributes(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, attributes, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Attribute AddAttributes error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	attrs, err := testZptSdk.Native.ZptId.GetAttributes(testIdentity.ID)
	if len(attributes) != len(attrs) {
		t.Errorf("TestZptId_Attribute GetAttributes len:%d != %d", len(attrs), len(attributes))
		return
	}
	if string(attr1.Key) != string(attrs[0].Key) || string(attr1.Value) != string(attrs[0].Value) || string(attr1.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestZptId_Attribute attribute:%s != %s", attrs[0], attr1)
		return
	}

	_, err = testZptSdk.Native.ZptId.RemoveAttribute(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, attr1.Key, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Attribute RemoveAttribute error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	attrs, err = testZptSdk.Native.ZptId.GetAttributes(testIdentity.ID)
	if len(attrs) != 1 {
		t.Errorf("TestZptId_Attribute GetAttributes len:%d != 1", len(attrs))
		return
	}
	if string(attr2.Key) != string(attrs[0].Key) || string(attr2.Value) != string(attrs[0].Value) || string(attr2.ValueType) != string(attrs[0].ValueType) {
		t.Errorf("TestZptId_Attribute attribute:%s != %s", attrs[0], attr2)
		return
	}
}

func TestZptId_Recovery(t *testing.T) {
	testIdentity, err := testWallet.NewDefaultSettingIdentity(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Recovery NewDefaultSettingIdentity error:%s", err)
		return
	}
	testDefController, err := testIdentity.GetControllerByIndex(1, testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Recovery GetControllerByIndex error:%s", err)
		return
	}
	_, err = testZptSdk.Native.ZptId.RegIDWithPublicKey(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Recovery RegIDWithPublicKey error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	_, err = testZptSdk.Native.ZptId.SetRecovery(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, testDefAcc.Address, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Recovery SetRecovery error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	ddo, err := testZptSdk.Native.ZptId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestZptId_Recovery GetDDO error:%s", err)
		return
	}
	if ddo.Recovery != testDefAcc.Address.ToBase58() {
		t.Errorf("TestZptId_Recovery recovery address:%s != %s", ddo.Recovery, testDefAcc.Address.ToBase58())
		return
	}

	acc1, err := testWallet.NewDefaultSettingAccount(testPasswd)
	if err != nil {
		t.Errorf("TestZptId_Recovery NewDefaultSettingAccount error:%s", err)
		return
	}

	txHash, err := testZptSdk.Native.ZptId.SetRecovery(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, acc1.Address, testDefController)

	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	evt, err := testZptSdk.GetSmartContractEvent(txHash.ToHexString())
	if err != nil {
		t.Errorf("TestZptId_Recovery GetSmartContractEvent:%s error:%s", txHash.ToHexString(), err)
		return
	}
	if evt.State == 1 {
		t.Errorf("TestZptId_Recovery duplicate add recovery should failed")
		return
	}
	_, err = testZptSdk.Native.ZptId.ChangeRecovery(testGasPrice, testGasLimit, testDefAcc, testIdentity.ID, acc1.Address, testDefAcc.Address, testDefController)
	if err != nil {
		t.Errorf("TestZptId_Recovery ChangeRecovery error:%s", err)
		return
	}
	testZptSdk.WaitForGenerateBlock(30*time.Second, 1)
	ddo, err = testZptSdk.Native.ZptId.GetDDO(testIdentity.ID)
	if err != nil {
		t.Errorf("TestZptId_Recovery GetDDO error:%s", err)
		return
	}
	if ddo.Recovery != acc1.Address.ToBase58() {
		t.Errorf("TestZptId_Recovery recovery address:%s != %s", ddo.Recovery, acc1.Address.ToBase58())
		return
	}
}
