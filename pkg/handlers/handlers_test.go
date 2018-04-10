package handlers

import (
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/gobuffalo/pop"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
)

type HandlerSuite struct {
	suite.Suite
	db           *pop.Connection
	logger       *zap.Logger
	filesToClose []*os.File
}

func (suite *HandlerSuite) SetupTest() {
	suite.db.TruncateAll()
}

func (suite *HandlerSuite) mustSave(model interface{}) {
	verrs, err := suite.db.ValidateAndSave(model)
	if err != nil {
		suite.T().Fatalf("Errors encountered saving %v: %v", model, err)
	}
	if verrs.HasAny() {
		suite.T().Fatalf("Validation errors encountered saving %v: %v", model, verrs)
	}
}

func (suite *HandlerSuite) checkErrorResponse(resp middleware.Responder, code int, name string) {
	errResponse, ok := resp.(*ErrResponse)
	if !ok || errResponse.code != code {
		suite.T().Fatalf("Expected %s Response: %v", name, resp)
	}
}

func (suite *HandlerSuite) checkResponseBadRequest(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusBadRequest, "BadRequest")
}

func (suite *HandlerSuite) checkResponseUnauthorized(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusUnauthorized, "Unauthorized")
}

func (suite *HandlerSuite) checkResponseForbidden(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusForbidden, "Forbidden")
}

func (suite *HandlerSuite) checkResponseNotFound(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusNotFound, "NotFound")
}

func (suite *HandlerSuite) checkResponseInternalServerError(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusInternalServerError, "InternalServerError")
}

func (suite *HandlerSuite) checkResponseTeapot(resp middleware.Responder) {
	suite.checkErrorResponse(resp, http.StatusTeapot, "Teapot")
}

func (suite *HandlerSuite) fixture(name string) *runtime.File {
	fixtureDir := "fixtures"
	cwd, err := os.Getwd()
	if err != nil {
		suite.T().Fatal(err)
	}

	fixturePath := path.Join(cwd, fixtureDir, name)

	info, err := os.Stat(fixturePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	header := multipart.FileHeader{
		Filename: name,
		Size:     info.Size(),
	}
	data, err := os.Open(fixturePath)
	if err != nil {
		suite.T().Fatal(err)
	}
	suite.closeFile(data)
	return &runtime.File{
		Header: &header,
		Data:   data,
	}
}

func (suite *HandlerSuite) AfterTest() {
	for _, file := range suite.filesToClose {
		file.Close()
	}
}

func (suite *HandlerSuite) closeFile(file *os.File) {
	suite.filesToClose = append(suite.filesToClose, file)
}

func TestHandlerSuite(t *testing.T) {
	configLocation := "../../config"
	pop.AddLookupPaths(configLocation)
	db, err := pop.Connect("test")
	if err != nil {
		log.Panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Panic(err)
	}
	hs := &HandlerSuite{db: db, logger: logger}
	suite.Run(t, hs)
}
