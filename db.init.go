package main

import (
	"fmt"
)

func (s *db_store) addBaseData() {
	const agent_id = "92aa258f95834a8bb35f74d5c21787d8"
	query := fmt.Sprintf(
		`INSERT INTO agent (ID, Name, Level, APIKey, Category, ThirdParty) VALUES ('%s', 'AB_0001', 0, 'MTVkYjll', '0,1,2,3', '{"scheme":"","host":"","path":"Api","key":"8f1W32f8gT9"}');`,
		agent_id,
	)
	e := s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}

	query = fmt.Sprintf(`CALL generate_wallet('%s');`, agent_id)
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}

	query = fmt.Sprintf(`CALL generate_user('%s');`, agent_id)
	e = s.execSQLText(query)
	if e != nil {
		s.Panic(e)
		return
	}
}
