package encryption

/*
CreateKey Creates key for encryption/decryption - NOTE: create key, iv and, mackey with this
*/
func CreateKey() ([]byte, error) {
	return nil, nil
}

/*
CreateSid Creates a session ID to manage encryption keys and its data
*/
func CreateSid() ([]byte, error) {
	return nil, nil
}

/*
KeyToString Converts a byte array key to a string
*/
func KeyToString(key []byte) (string, error) {
	return "", nil
}

/*
StringToKey Converts the string given to a byte array key
*/
func StringToKey(str string) ([]byte, error) {
	return nil, nil
}

/*
EncryptAndSign Encrypts the payload and signs it
*/
func EncryptAndSign(key []byte, iv []byte, mackey []byte, data []byte) ([]byte, error) {
	return nil, nil
}

/*
AuthAndDecrypt Authenticates the payload and decrypts it
*/
func AuthAndDecrypt(key []byte, iv []byte, mackey []byte, data []byte) ([]byte, error) {
	return nil, nil
}

/*
Encrypt Encrypts a payload - NOTE: key must be unique and secure
*/
func Encrypt(key []byte, iv []byte, data []byte) ([]byte, error) {
	return nil, nil
}

/*
Decrypt Decrypts an encrypted payload
*/
func Decrypt(key []byte, iv []byte, size int, data []byte) ([]byte, error) {
	return nil, nil
}

/*
Sign Signs a payload and creates a 32 byte-long signature and copies it
*/
func Sign(mackey []byte, data []byte) {
}

/*
Auth Authenticates a signed payload w/ mac key and signature
*/
func Auth(mackey []byte, signature []byte, data []byte) bool {
	return false
}
