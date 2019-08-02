package service

import (
	"github.com/integration-system/isp-journal/client"
	"github.com/integration-system/isp-lib/http"
)

var (
	JournalServiceClient = client.NewJournalServiceClient(http.NewJsonRestClient())
)
