package aesencryption

import "testing"

func TestGenerateSalt(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Errorf("GenerateSalt() returns error: ", err.Error())
	}
	if len(salt) != saltSize {
		t.Errorf("Salt should be of length 8, it is ", len(salt))
	}
	salt2, err := GenerateSalt()
	if err != nil {
		t.Errorf("GenerateSalt() returns error: ", err.Error())
	}
	e := true
	for i := 0; i < len(salt); i++ {
		if salt[i] != salt2[i] {
			e = false
		}
	}
	if e {
		t.Errorf("salt should not be the same as salt2")
	}
}

func TestExtractFunctions(t *testing.T) {
	// write these later
}

func TestMakeAES256Key(t *testing.T) {
	salt, err := GenerateSalt()
	if err != nil {
		t.Errorf("GenerateSalt() returned error: ", err)
	}
	aeskey, err := MakeAES256Key("test passphrase", salt)
	if err != nil {
		t.Errorf("MakeAES256Key() returned error: ", err.Error())
	}
	if len(aeskey.key) != 32 {
		t.Errorf("AES key should be of length 32 (256-bit) it is not")
	}
	if len(aeskey.salt) != saltSize {
		t.Errorf("AES key salt should be of length 8, it is not")
	}
	aeskey2, err := MakeAES256Key("test passphrase", salt)
	if err != nil {
		t.Errorf("MakeAES256Key() returned error: ", err.Error())
	}
	e := false
	for i := 0; i < len(aeskey.key); i++ {
		if aeskey.key[i] != aeskey2.key[i] {
			e = true
		}
	}
	if e {
		t.Errorf("MakeAES256Key() should be deterministic for same passphrase and salt")
	}
	aeskey3, err := MakeAES256Key("test passphrase", nil)
	if err != nil {
		t.Errorf("MakeAES256Key() returned error: ", err.Error())
	}
	e = true
	for i := 0; i < len(aeskey.key); i++ {
		if aeskey.key[i] != aeskey3.key[i] {
			e = false
		}
	}
	if e {
		t.Errorf("MakeAES256Key() with nil salt parameter should generate new key")
	}
}

func TestEncryption(t *testing.T) {
	plainData := []byte("This is my super secret secret. Keep safe pls. Ty.")
	passphrase := "test passphrase"
	s, _ := GenerateSalt()

	aeskey, err := MakeAES256Key(passphrase, s)

	// Make sure there are now errors with MakeAES256Key()
	if err != nil {
		t.Errorf("MakeAES256Key() gave error: ", err)
	}

	msg, err := EncryptAES256(*aeskey, plainData)
	// Make sure there are no errors with EncryptAES256()
	if err != nil {
		t.Errorf("EncryptAES256() returned error: ", err)
	}
	ciphertext, err := ExtractMsg(msg)
	// Make sure there are no errors with ExtractMsg()
	if err != nil {
		t.Errorf("ExtractMsg() returned error: ", err)
	}
	plaintext, err := DecryptAES256(*aeskey, ciphertext)

	// Check that the length of plaintext and plainData are the same
	if len(plaintext) != len(plainData) {
		t.Errorf("plaintext should have same length as plainData")
	}

	// Check that plaintext and plainData are indeed idenitcal
	e := false
	for i := 0; i < len(plaintext); i++ {
		if plaintext[i] != plainData[i] {
			e = true
		}
	}
	if e {
		t.Errorf("Decryption should return plainData, it does not")
	}

}
