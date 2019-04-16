//go_utils package library with helper functions
package go_utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	TimeLayoutYYYYMMDD_HHMMSS = "2006-01-02 15:04:05"
	TimeFormatLogFileName     = "2006-01-02T15:00"
)

type ErrorResponse struct {
	Status      string `json:"status"`
	Description string `json:"description"`
}

//Set a unique id for every request
func RequestIdMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewV4()
		c.Writer.Header().Set("X-Request-Id", id.String())
		c.Next()
	}
}

func SetErrorResponse(description string) ErrorResponse {
	return ErrorResponse{
		Status:      "error",
		Description: description,
	}
}

func SetWarningResponse(description string) ErrorResponse {
	return ErrorResponse{
		Status:      "warning",
		Description: description,
	}
}

//GetCurrentDir gets directory of the go binary file and returns the path of it
func GetCurrentDir() (error, string) {
	ex, err := os.Executable()
	if err != nil {
		return err, ""
	}
	exPath := filepath.Dir(ex)

	return nil, exPath
}

func GetNowStd() string {
	return time.Now().Format(TimeLayoutYYYYMMDD_HHMMSS)
}

//set and configure gin framework Logger
func SetLogger() {
	if err := os.Mkdir("Log", os.ModePerm); err != nil && !os.IsExist(err) {
		Logger("Error on r.Run: ", fmt.Sprint(err))
		panic(err.Error())
	}
	if f, err := GetLogFile(); err != nil {
		Logger("Error on r.Run: ", err.Error())
		panic(err.Error())
	} else {
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
}

//Logger gets msg string and print it as JSON with a dateStamp
func Logger(msg ...string) {
	var buffer bytes.Buffer
	buffer.WriteString(`{'dateStamp':'`)
	buffer.WriteString(GetNowStd())
	buffer.WriteString(`','message':'`)
	for i := 0; i <= len(msg)-1; i++ {
		buffer.WriteString(msg[i])
	}
	buffer.WriteString(`'}`)
	if gin.IsDebugging() || gin.Mode() == "test" {
		fmt.Println(buffer.String())
	} else {
		//production logging actions
	}
}

//CreateJson gets an interface{} and marshal it to JSON and returns as []byte
func CreateJson(i interface{}) []byte {
	if jsonReply, err := json.Marshal(&i); err != nil {
		return nil
	} else {
		return jsonReply
	}
}

//LogAsJson takes multiple arguments and returns them as a JSON formatted string
func LogAsJson(logMsg ...interface{}) string {
	var res []interface{}
	for _, msg := range logMsg {
		switch msg.(type) {
		case string:
			res = append(res, msg.(string))
			break
		default:
			res = append(res, string(CreateJson(msg)))
		}
	}
	return fmt.Sprint(res...)
}

//GetLogFile returns log file from GetCurrentDir()/Log/currentDir/TimeFormatLogFileName.gin.log
//if file doesn't exist creates one
//it makes one log file every hour
func GetLogFile() (*os.File, error) {
	t := time.Now()
	err, currentDir := GetCurrentDir()
	if err != nil {
		return nil, err
	}

	currentDir += "/Log"
	strings.TrimRight(currentDir, "/")
	var fileName string = fmt.Sprintf("%s/%s.gin.log", currentDir, t.Format(TimeFormatLogFileName))
	if fileName == "" {
		return nil, errors.New("empty file name")
	}

	if strings.Index(fileName, "/") == -1 {
		fileName = currentDir + "/Log/" + fileName
	}
	logFile, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}

	return logFile, nil
}
