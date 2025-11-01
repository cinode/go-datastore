/*
Copyright © 2025 Bartłomiej Święcki (byo)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package testutils

import (
	"encoding/hex"

	"github.com/cinode/go-common/blob"
	"github.com/cinode/go-common/cutl"
)

// nolint:lll // test data vectors
var TestBlobs = []struct {
	Name     *blob.Name
	Data     []byte
	Expected []byte
}{
	{
		cutl.Must(blob.NameFromString("KDc2ijtWc9mGxb5hP29YSBgkMLH8wCWnVimpvP3M6jdAk")),
		cutl.Must(hex.DecodeString("54657374")),
		cutl.Must(hex.DecodeString("54657374")),
	},
	{
		cutl.Must(blob.NameFromString("BG8WaXMAckEfbCuoiHpx2oMAS4zAaPqAqrgf5Q3YNzmHx")),
		cutl.Must(hex.DecodeString("5465737431")),
		cutl.Must(hex.DecodeString("5465737431")),
	},
	{
		cutl.Must(blob.NameFromString("2GLoj4Bk7SvjQngCT85gxWRu2DXCCjs9XWKsSpM85Wq3Ve")),
		cutl.Must(hex.DecodeString("")),
		cutl.Must(hex.DecodeString("")),
	},
	{
		cutl.Must(blob.NameFromString("251SEdnHjwyvUqX1EZnuKruta4yHMkTDed7LGoi3nUJwhx")),
		cutl.Must(hex.DecodeString("00c44bbf343e578f995dc0e8b4d1119c64973003a2ad68a3b6d1ce6ed9a0a79f7b7daace91b391cace11617f9af0bbb83c9d8a18d930d099444f0f78c1c706511339c272a05403e635ca1d237c936878cfcb3c2456498bf2ca4e12e4246f8403f6b8c462dabfeadc0e0000000000000000189afcc9e1dba6ee8e3926cf2b22da225a5a59806eee5cc737c11b6f9a4f19082414b54cfba2f3498045e1a0c232e532a20ba81f98122839d6867de9b929251cf43d5c7c7be78ac0cbebf507c2c601422a2a6c33b1a2ba1525130673d0c040")),
		cutl.Must(hex.DecodeString("c11b6f9a4f19082414b54cfba2f3498045e1a0c232e532a20ba81f98122839d6867de9b929251cf43d5c7c7be78ac0cbebf507c2c601422a2a6c33b1a2ba1525130673d0c040")),
	},
	{
		cutl.Must(blob.NameFromString("27vP1JG4VJNZvQJ4Zfhy3H5xKugurbh89B7rKTcStM9guB")),
		cutl.Must(hex.DecodeString("008266a49daf44cd25f425ad31ffe88a21d160e9da947b51e733271f8ac71652f876882dc5dbc0dd3eead04bd6d49e826a7ed04a039e825b5353d52054209c18d19803cb68c63acc89e11a284b82a7515ef3b6eb64035d1b9052ea198c1874860afa1d98052f4f1304000000000000000118c001789a1e57030e92a620b6f89abf005bb1f7e389ddaa4a7315d299c9af674bfbdf81ba59fb8d99ad5d5e07779153e90bafaf6b2b755c1de90c0f6118bbc49096e090b5b4a6d17cccabfd10ea957b81552a63da648c5cc27d9ac4fecf630f")),
		cutl.Must(hex.DecodeString("7315d299c9af674bfbdf81ba59fb8d99ad5d5e07779153e90bafaf6b2b755c1de90c0f6118bbc49096e090b5b4a6d17cccabfd10ea957b81552a63da648c5cc27d9ac4fecf630f")),
	},
	{
		cutl.Must(blob.NameFromString("e3T1HcdDLc73NHed2SFu5XHUQx5KDwgdAYTMmmEk2Ekqm")),
		cutl.Must(hex.DecodeString("00628297af66c4e51f1f8d7c491e240ced24cb9172e0056ec3f008d781ecefe177d7b068640f5feec1bafa2974c8b3a793bdcc25b6a7032199b6bd690cb47ff576ec37b36f3fdf9787afbc056fbd42e04e3bad2950aabf80d7f5f8e5438ed6b256c4c5be65f8ba5c0d0000000000000002188ccb210ee53218f74d24c14a9417394a767d74529d6a3618462bd895adfbd18dbc2b95074edd121e734edbde1c9c8e10aea5f7cbfedc792030c43539b6b96c7fee06f003db9e62fcb76ed659983ca73d3fca600c9dafc9e282a5")),
		cutl.Must(hex.DecodeString("462bd895adfbd18dbc2b95074edd121e734edbde1c9c8e10aea5f7cbfedc792030c43539b6b96c7fee06f003db9e62fcb76ed659983ca73d3fca600c9dafc9e282a5")),
	},
}

// nolint:lll // test data vectors
var DynamicLinkPropagationData = []struct {
	Name     *blob.Name
	Data     []byte
	Expected []byte
}{
	{
		cutl.Must(blob.NameFromString("GUnL66Lyv2Qs4baxPhy59kF4dsB9HWakvTjMBjNGFLT6g")),
		cutl.Must(hex.DecodeString("0016170c399aa4af538d1fa5c58eb87a48350796f86ab350451ad155dc74d6ca2aab63effca2df99dbd72b548cf522d00e1e53e8de6c3a46f69894e9e4b0c400054a03f8d5857cb1fb55447404b3ab47aea74a7811b98063210c0b2397928e98589e9273c9afaa1100000000000000271018146a18768320ddd9ebc54edc464ccff1fc48dd401dfa4cbc751beb66b186c68be88fe3599d1d36d6974d41221e96d4b1578cf27708f74f79a3514596edeb9f7d100ea923940a5a9566898a746e6833a7f148b8f2e0cb3c4ed9383a225c0a9b")),
		cutl.Must(hex.DecodeString("751beb66b186c68be88fe3599d1d36d6974d41221e96d4b1578cf27708f74f79a3514596edeb9f7d100ea923940a5a9566898a746e6833a7f148b8f2e0cb3c4ed9383a225c0a9b")),
	},
	{
		cutl.Must(blob.NameFromString("GUnL66Lyv2Qs4baxPhy59kF4dsB9HWakvTjMBjNGFLT6g")),
		cutl.Must(hex.DecodeString("0016170c399aa4af538d1fa5c58eb87a48350796f86ab350451ad155dc74d6ca2aab63effca2df99db7a2d9ef0c9cfac858a60942e250c1266ee5ccafb96fe56e1187ed334ab25d398baff9a69e062237b746b931e3f75ff46196e20d99a1ba7ee6ad8649e8e7375090000000000004e2018e6bd8bf59d9de1ab35a134ed9d15eaaa70064964f648cbf14fe66148a3af71d0a4f4a1817dda8f188b13cc03a122b62529a0ac3894f2654ed362e8ff78d24a7d8c3ba5120bfab4e7d7a8c260f3961732d6efdf25fcde89a530ab8e548e2780")),
		cutl.Must(hex.DecodeString("4fe66148a3af71d0a4f4a1817dda8f188b13cc03a122b62529a0ac3894f2654ed362e8ff78d24a7d8c3ba5120bfab4e7d7a8c260f3961732d6efdf25fcde89a530ab8e548e2780")),
	},
	{
		cutl.Must(blob.NameFromString("GUnL66Lyv2Qs4baxPhy59kF4dsB9HWakvTjMBjNGFLT6g")),
		cutl.Must(hex.DecodeString("0016170c399aa4af538d1fa5c58eb87a48350796f86ab350451ad155dc74d6ca2aab63effca2df99db7388571b272b87ba326c5c1f68daac7ebac471986efc9da1bbeb931fc3092753d8ac969c64db5e8d9d8631d8234ac4f3f5174e5fa01dd1818108cadf89b380030000000000004e201872d3ee3f68ff92213de311a2c138c11edc142470b92dd8a07da9d4c95b2d2fcdd3949a68a2844ea6af28d8ac45507eda0a9956503b6fc808e0b4e693f73adb8f50913b33135b9f13a7bc9b9087d71e136d58b6b069d642a8e734cc6bb953c9")),
		cutl.Must(hex.DecodeString("7da9d4c95b2d2fcdd3949a68a2844ea6af28d8ac45507eda0a9956503b6fc808e0b4e693f73adb8f50913b33135b9f13a7bc9b9087d71e136d58b6b069d642a8e734cc6bb953c9")),
	},
}
