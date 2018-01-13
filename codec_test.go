package zero

import "testing"

func TestCodec(t *testing.T) {
	// test encode
	msg1 := NewMessage(1, []byte("message codec test..."))

	data, err := Encode(msg1)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg1)

	// test decode
	// The first four bytes is size for socket read
	msg2, err := Decode(data[4:])
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ID=%d, Data=%s", msg2.GetID(), string(msg2.GetData()))
}
