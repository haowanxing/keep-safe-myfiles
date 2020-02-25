package enc

import "testing"

func TestAESEncrypt(t *testing.T) {
	orign := []byte("Hello World!")
	key, iv := []byte("GoodPeopleGoodMe"), []byte("0000000000000000")
	crypt, err := AESEncrypt(orign, key, iv)
	if err != nil {
		t.Error(err)
	}
	t.Log(crypt, len(crypt))
	o2, err := AESDecrypt(crypt, key, iv)
	if err != nil {
		t.Error(err)
	}
	t.Log(o2)
}
