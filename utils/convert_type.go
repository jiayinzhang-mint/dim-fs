package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"strconv"

	"github.com/google/uuid"
)

// UintToString convert uint to string
func UintToString(v uint) string {
	return strconv.Itoa(int(v))
}

// ParamsArrayToStringArray convert array in params to []string for sql
func ParamsArrayToStringArray(v string) []string {
	vMap := make(map[string]interface{})
	json.Unmarshal([]byte(v), &vMap)
	fmt.Println(v)

	vArray := make([]string, len(vMap))
	for i, item := range vMap {
		index, _ := strconv.Atoi(i)
		vArray[index] = fmt.Sprint(item)
	}

	return vArray
}

// IntArrayToStringArray convert int array to str array for []int like get params
func IntArrayToStringArray(iArray []string) []string {
	sArray := make([]string, len(iArray))
	for i, item := range iArray {
		x, _ := strconv.Atoi(item)
		sArray[i] = strconv.Itoa(x)
	}
	return sArray
}

// InterfaceToString convert *xxx to string for jsonb use
func InterfaceToString(object interface{}) (string, error) {
	objectJSON, encodeErr := json.Marshal(object)
	if encodeErr != nil {
		return "", fmt.Errorf("encode error")
	}

	var v bytes.Buffer
	v.WriteString(string(objectJSON))

	return v.String(), nil
}

// EncodeUUID helps shorten url
func EncodeUUID(uuid uuid.UUID) string {
	uuidByte, _ := uuid.MarshalBinary()
	return base64.RawURLEncoding.EncodeToString(uuidByte)
}

// DecodeUUID helps shorten url
func DecodeUUID(uuidEncoded string) uuid.UUID {
	decodedByte, _ := base64.RawURLEncoding.DecodeString(uuidEncoded)
	decoded, _ := uuid.FromBytes(decodedByte)

	return decoded
}

// ParseUUID parse string to uuid
func ParseUUID(s string) (u uuid.UUID) {
	u, _ = uuid.Parse(s)
	return u
}

// Float64frombytes convert []uint8 to float64
func Float64frombytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}
