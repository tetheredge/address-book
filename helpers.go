package main

type jsonError struct {
	Code    int
	Message string
	Err     string
}

type jsonData struct {
	Code int
	Byte []byte
}
