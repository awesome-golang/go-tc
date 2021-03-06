package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaIptUnspec = iota
	tcaIptTable
	tcaIptHook
	tcaIptIndex
	tcaIptCnt
	tcaIptTm
	tcaIptTarg
	tcaIptPad
)

// Ipt contains attribute of the ipt discipline
type Ipt struct {
	Table string
	Hook  uint32
	Index uint32
	Cnt   *IptCnt
	Tm    *Tcft
}

// IptCnt as tc_cnt from include/uapi/linux/pkt_cls.h
type IptCnt struct {
	RefCnt  uint32
	BindCnt uint32
}

// unmarshalIpt parses the ipt-encoded data and stores the result in the value pointed to by info.
func unmarshalIpt(data []byte, info *Ipt) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaIptTm:
			tcft := &Tcft{}
			if err := unmarshalStruct(ad.Bytes(), tcft); err != nil {
				return err
			}
			info.Tm = tcft
		case tcaIptTable:
			info.Table = ad.String()
		case tcaIptHook:
			info.Hook = ad.Uint32()
		case tcaIptIndex:
			info.Index = ad.Uint32()
		case tcaIptCnt:
			tmp := &IptCnt{}
			if err := unmarshalStruct(ad.Bytes(), tmp); err != nil {
				return err
			}
			info.Cnt = tmp
		case tcaIptPad:
			// padding does not contain data, we just skip it
		default:
			return fmt.Errorf("UnmarshalIpt()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalIpt returns the binary encoding of Ipt
func marshalIpt(info *Ipt) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Ipt: %w", ErrNoArg)
	}
	// TODO: improve logic and check combinations
	if info.Tm != nil {
		return []byte{}, ErrNoArgAlter
	}
	if len(info.Table) > 0 {
		options = append(options, tcOption{Interpretation: vtString, Type: tcaIptTable, Data: info.Table})
	}
	options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIptHook, Data: info.Hook})
	options = append(options, tcOption{Interpretation: vtUint32, Type: tcaIptIndex, Data: info.Index})
	options = append(options, tcOption{Interpretation: vtString, Type: tcaIptTable, Data: info.Table})
	if info.Cnt != nil {
		data, err := marshalStruct(info.Cnt)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaIptCnt, Data: data})
	}

	return marshalAttributes(options)
}
