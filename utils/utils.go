package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()
func ParseJson(r *http.Request,payload any) error{
if r.Body == nil {
	fmt.Errorf("Empty request")
}
  return json.NewDecoder(r.Body).Decode(payload);
}

func WriteJson (w http.ResponseWriter, status int , v any) error{
  w.Header().Add("Content-Type","application/json")
  w.WriteHeader(status)
  return json.NewEncoder(w).Encode(v)
}

func WriteError(w http.ResponseWriter, status int,err error){
	WriteJson(w , status,map[string]string{"error":err.Error()})
}

func StringToUint(str string) (uint, error) {
    value, err := strconv.ParseUint(str, 10, 32)
    if err != nil {
        return 0, err
    }
    return uint(value), nil
}
