package routes

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/vtmhieu/VCS_SMS/models"
)

// type MockRepository struct {
// 	mock.Mock
// }

// func (mock *MockRepository) create(server models.Server) error {
// 	arg := mock.Called()
// 	return arg.Error(0)
// }

// func (mock *MockRepository) find(id string) (models.Server, error) {
// 	args := mock.Called()
// 	result := args.Get(0)
// 	return result.(models.Server), args.Error(1)

// }

// func TestCreateServer(t *testing.T) {
// 	mocktest := new(MockRepository)
// 	server := models.Create_server{
// 		Server_id:   "12,s",
// 		Server_name: "test",
// 		Status:      "Offline",
// 		Ipv4:        "112.221.1332.122",
// 	}

// 	t.Run("Successfully", func(t *testing.T) {
// 		mocktest.On("find").Return(models.Server{}, gorm.ErrRecordNotFound).Once()
// 		mocktest.On("create").Return(nil).Once()

// 		NewServer := controllers.CreateServer(mocktest)

// 	})
// }

func SetUpRouter() *gin.Engine {
	server := gin.Default()
	router := server.Group("/api")
	return
}

func TestCreateServer(t *testing.T) {
	r := SetUpRouter()
	var c Server_Route_Controller
	r.POST("/", c.servercontroller.CreateServer)
	now := time.Now()
	server := models.Server{
		Server_id:    "demo",
		Server_name:  "demo_name",
		Status:       "demo_status",
		Created_time: now,
		Last_updated: now,
		Ipv4:         "1222.111.1220",
	}
	jsonValue, _ := json.Marshal(server)
	req, _ := http.NewRequest("POST", "/", bytes.NewBuffer(jsonValue))
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}
