package types

type RevIP4 struct {
	Byte   [4]byte
	String string
}

//AddrNames includes address and its names array
type AddrNames struct {
	Addr  string
	Names []string
}
