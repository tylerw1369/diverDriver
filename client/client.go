package client

import (
	"github.com/tylerw1369/diverdriver/client/ipcclient"
	"github.com/tylerw1369/diverdriver/client/remoteclient"
	"github.com/tylerw1369/diverdriver/common"
	"github.com/tylerw1369/diverdriver/utils"
)

func Initialize(diverDriverPath string, writeTimeOutMs int64, readTimeOutMs int) *common.DiverClient {
	p := &common.DiverClient{DiverDriverPath: diverDriverPath, WriteTimeOutMs: writeTimeOutMs, ReadTimeOutMs: readTimeOutMs}
	if utils.IsValidRemoteURL(p.DiverDriverPath) {
		p.PowClientImplementation = remoteclient.RemoteClient
	} else {
		p.PowClientImplementation = ipcclient.IpcClient
	}
	return p
}
