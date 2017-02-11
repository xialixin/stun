package stun

import (
	"net"
	"testing"
)

func BenchmarkFingerprint_AddTo(b *testing.B) {
	b.ReportAllocs()
	m := new(Message)
	s := NewSoftware("software")
	addr := &XORMappedAddress{
		IP: net.IPv4(213, 1, 223, 5),
	}
	addAttr(b, m, addr)
	addAttr(b, m, s)
	b.SetBytes(int64(len(m.Raw)))
	for i := 0; i < b.N; i++ {
		Fingerprint.AddTo(m)
		m.WriteLength()
		m.Length -= attributeHeaderSize + fingerprintSize
		m.Raw = m.Raw[:m.Length+messageHeaderSize]
		m.Attributes = m.Attributes[:len(m.Attributes)-1]
	}
}

func TestFingerprint_Check(t *testing.T) {
	m := new(Message)
	addAttr(t, m, NewSoftware("software"))
	m.WriteHeader()
	Fingerprint.AddTo(m)
	m.WriteHeader()
	if err := Fingerprint.Check(m); err != nil {
		t.Error(err)
	}
}

func BenchmarkFingerprint_Check(b *testing.B) {
	b.ReportAllocs()
	m := new(Message)
	s := NewSoftware("software")
	addr := &XORMappedAddress{
		IP: net.IPv4(213, 1, 223, 5),
	}
	addAttr(b, m, addr)
	addAttr(b, m, s)
	m.WriteHeader()
	Fingerprint.AddTo(m)
	m.WriteHeader()
	b.SetBytes(int64(len(m.Raw)))
	for i := 0; i < b.N; i++ {
		if err := Fingerprint.Check(m); err != nil {
			b.Fatal(err)
		}
	}
}