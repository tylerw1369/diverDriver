package client

import (
	"math/rand"
	"testing"

	"github.com/tylerw1369/fpgago"
)

const (
	TRYTE_CHARS     = "9ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	MWM         int = 14
)

var transaction = "999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999A9RGRKVGWMWMKOLVMDFWJUHNUNYWZTJADGGPZGXNLERLXYWJE9WQHWWBMCPZMVVMJUMWWBLZLNMLDCGDJ999999999999999999999999999999999999999999999999999999YGYQIVD99999999999999999999TXEFLKNPJRBYZPORHZU9CEMFIFVVQBUSTDGSJCZMBTZCDTTJVUFPTCCVHHORPMGCURKTH9VGJIXUQJVHK999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"

func TestLocalPOW(t *testing.T) {
	var localPath = "/tmp/diverDriver.sock"
	diverClient := Initialize(localPath, 500, 5000)

	serverVersion, powType, powVersion, err := diverClient.GetPowInfo()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("ServerVersion: %v, PowType: %v, PowVersion: %v", serverVersion, powType, powVersion)

	// test transaction data
	randomTrytes := make([]rune, 256)

	for i := 0; i < 10000; i++ {
		for j := 0; j < 256; j++ {
			randomTrytes[j] = rune(TRYTE_CHARS[rand.Intn(len(TRYTE_CHARS))])
		}

		data, err := giota.ToTrytes(string(randomTrytes) + transaction[256:])
		if err != nil {
			t.Error(err)
			continue
		}

		response, err := diverClient.PowFunc(data, MWM)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("Client received: %v", response)
	}
}

func TestRemotePOW(t *testing.T) {
	var remotePath = "http://192.168.178.101:16000"
	diverClient := Initialize(remotePath, 0, 0)

	serverVersion, powType, powVersion, err := diverClient.GetPowInfo()
	if err != nil {
		t.Error(err)
		return
	}

	t.Logf("ServerVersion: %v, PowType: %v, PowVersion: %v", serverVersion, powType, powVersion)

	// test transaction data
	randomTrytes := make([]rune, 256)

	for i := 0; i < 100; i++ {
		for j := 0; j < 256; j++ {
			randomTrytes[j] = rune(TRYTE_CHARS[rand.Intn(len(TRYTE_CHARS))])
		}

		data, err := giota.ToTrytes(string(randomTrytes) + transaction[256:])
		if err != nil {
			t.Error(err)
			continue
		}

		response, err := diverClient.PowFunc(data, MWM)
		if err != nil {
			t.Error(err)
			continue
		}
		t.Logf("Client received: %v", response)
	}
}
