package utils

import (
	"fmt"
	"io"
	"net"
	"time"
)

func ReadByte(input net.Conn) (result []byte) {
	input.SetReadDeadline(time.Now().Add(500 * time.Second))
	bufsize := 128
	blocksize := 8192
	data := make([]byte, blocksize)
	len := 0
	var num int
	var err error
	for {
		if (len + bufsize) <= blocksize {
			num, err = input.Read(data[len : len+bufsize])
		} else {
			num, err = input.Read(data[len:blocksize])
		}
		if err == nil && num > 0 {
			len = num
			if len == blocksize {
				len = 0
			}
		}

		if num < bufsize {
			input.SetReadDeadline(time.Now().Add(time.Microsecond * 50))
		}
		if err != nil {
			fmt.Println("===============================")
			result = append(result, data[0:len]...)
			break
		}
	}
	return
}

func ReadBytesWithLen(c net.Conn, bytesLen int64, timeout time.Duration) ([]byte, error) {
	c.SetReadDeadline(time.Now().Add(timeout))
	readBytes := make([]byte, 0, bytesLen)
	for bytesLen > 0 {
		buf := make([]byte, bytesLen, bytesLen)
		n, err := c.Read(buf)
		if err != nil {
			if err == io.EOF {
				return readBytes, nil
			}
			return nil, err
		}

		bytesLen -= int64(n)
		if n == 0 && bytesLen != 0 {
			return nil, fmt.Errorf("cannot read the %d length of bytes", bytesLen)
		}
		readBytes = append(readBytes, buf[0:n]...)
	}
	return readBytes, nil
}

func ConnectWithRetry(network, addr string, timeout time.Duration, retryNum int) (c net.Conn, err error) {
	for i := 0; i < retryNum; i++ {
		c, err = net.DialTimeout(
			network,
			addr,
			timeout,
		)
		if err == nil {
			return
		}
	}
	return
}

func WriteWithTimeout(c net.Conn, requestBytes []byte, timeout time.Duration) error {
	c.SetWriteDeadline(time.Now().Add(time.Duration(timeout)))
	writeLen := len(requestBytes)
	for writeLen > 0 {
		n, err := c.Write(requestBytes)
		if err != nil {
			return err
		}

		if n <= 0 {
			err = fmt.Errorf("write 0 bytes to connection for pab")
			return err
		}
		writeLen -= n
	}
	return nil
}
