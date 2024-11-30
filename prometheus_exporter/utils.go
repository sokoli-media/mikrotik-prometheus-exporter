package prometheus_exporter

import (
	"errors"
	"gopkg.in/routeros.v2/proto"
	"strconv"
)

func getKey(reply proto.Sentence, key string) (string, error) {
	for _key, value := range reply.Map {
		if _key == key {
			return value, nil
		}
	}
	return "", errors.New("key not found")
}

func getKeyAsFloat(reply proto.Sentence, key string) (float64, error) {
	stringValue, err := getKey(reply, key)
	if err != nil {
		return 0, err
	}

	floatValue, err := strconv.ParseFloat(stringValue, 32)
	if err != nil {
		return 0, err
	}

	return floatValue, nil
}
