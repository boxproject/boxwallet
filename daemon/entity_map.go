package daemon

import (
	"github.com/boxproject/boxwallet/db/mysql"
)

func ConvertTxInfo(m *TxInfo) *mysql.TxInfo {
	result := &mysql.TxInfo{
		TxId:          m.TxId,
		Confirmations: m.Cnfm,
		Token:         m.Token,
		BCType:        int(m.BCT),
		BlockH:        m.H.Uint64(),
		Fee:           m.Fee.String(),
		TxObj:         &mysql.TxObj{},
		ExtValid:      m.ExtValid,
		Target:        m.Target,
	}
	for _, v := range m.In {
		in := &mysql.AddrAmount{
			Addr: v.Addr,
			Amt:  v.Amt.String(),
		}

		result.TxObj.In = append(result.TxObj.In, in)
	}
	for _, v := range m.Out {
		out := &mysql.AddrAmount{
			Addr: v.Addr,
			Amt:  v.Amt.String(),
		}
		result.TxObj.Out = append(result.TxObj.Out, out)
	}
	for _, v := range m.InExt {
		in := &mysql.AddrAmount{
			Addr: v.Addr,
			Amt:  v.Amt.String(),
		}
		result.TxObj.ExtIn = append(result.TxObj.ExtIn, in)
	}
	for _, v := range m.OutExt {
		out := &mysql.AddrAmount{
			Addr: v.Addr,
			Amt:  v.Amt.String(),
		}
		result.TxObj.ExtOut = append(result.TxObj.ExtOut, out)
	}
	return result
}
