package sdk

import (
	"log"
	"testing"
)

func TestSign(t *testing.T) {
	osdk := OPPOSDK{
		appKey:    "11",
		appSecret: "22",
		pkgName:   "com.oppo.testgame",
	}
	log.Fatal(osdk.getVerifyUrlParam("ccb5ae1cf111de847db508edb90fa633"))
}

// appKey=11&appSecret=11&pkgName=com.oppo.testgame&timeStamp=1526304757000&token=ccb5ae1cf111de847db508edb90fa633
// C26228F50092BC9CA251AF41C4F70022
