package db

import (
	"go_game_server/server/include"
	"go_game_server/server/util"
)

type Acct include.Acct

func SaveAcct(acct *Acct) {
	if acct == nil {
		return
	}
	saveSQL := "replace into t_acct(acct_name, password)VALUES(?, ?)"
	ExecDB(UserDBType, -1, saveSQL,
		acct.AcctName, acct.Password)
	acct = GetAcct(acct.AcctName)
}

func GetAcct(acctName string) *Acct {
	selSQL := `SELECT acct_id, acct_name, password FROM t_acct WHERE acct_name = ?`
	rows, err := DB.Query(selSQL, acctName)
	util.CheckErr(err)
	for rows.Next() {
		tmp := &Acct{}
		err := rows.Scan(&tmp.AcctID, &tmp.AcctName, &tmp.Password)
		util.CheckErr(err)
		return tmp
	}
	return nil
}

func JudgeAcctName(acctName string) bool {
	exist := false
	err := DB.QueryRow("SELECT EXISTS (SELECT acct_name FROM t_acct WHERE acct_name = ?)", acctName).Scan(&exist)
	//fmt.Println(err)
	util.CheckErr(err)
	return exist
}
