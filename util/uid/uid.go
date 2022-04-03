package uid

import (
	"strconv"
	"strings"
)

type UID int64

func (u UID) MarshalJSON() ([]byte, error) {
	id := u.Encode()
	return []byte(strconv.Itoa(int(id))), nil
}

func (u *UID) UnmarshalJSON(b []byte) error {
	id, _ := strconv.Atoi(string(b))
	uid := UID(id)
	deId := uid.Decode()
	(*u) = deId
	return nil
}
func (u UID) Encode() int64 {
	if u == 0 {
		return int64(u)
	}
	var start int64 = 100000
	seq := []string{"0", "4", "2", "9", "1", "6", "7", "5", "8", "3"}
	i := int64(u) + start
	str := strconv.Itoa(int(i))

	arrStr := make([]string, 0, len(str))
	//arrStr = append(arrStr, seq[time.Now().Nanosecond()%9+1])

	for i, _ := range str {
		index := str[i : i+1]
		i, _ := strconv.Atoi(index)
		arrStr = append(arrStr, seq[i])
	}
	_newStr := strings.Join(arrStr, "")
	newI, _ := strconv.Atoi(_newStr)
	return int64(newI)
}

func (u UID) Decode() UID {
	if u == 0 {
		return u
	}
	var start int64 = 100000
	seq := []string{"0", "4", "2", "9", "1", "6", "7", "5", "8", "3"}
	str := strconv.Itoa(int(u))
	arrStr := make([]string, 0, len(str))

	tempMap := make(map[string]int, 0)
	for i, v := range seq {
		tempMap[v] = i
	}

	for i, _ := range str {
		index := str[i : i+1]
		s := strconv.Itoa(tempMap[index])
		arrStr = append(arrStr, s)
	}
	_newStr := strings.Join(arrStr, "")
	newId, _ := strconv.Atoi(_newStr)
	newId = newId - int(start)
	return UID(newId)
}
