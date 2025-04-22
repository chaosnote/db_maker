package main

import (
	"github.com/chaosnote/db_maker/utils"
)

func (ds *db_store) addSPDropTables() {
	query := `DROP PROCEDURE IF EXISTS drop_all_tables ;`
	e := ds.execSQLText(query)
	if e != nil {
		ds.Panic(utils.LogFields{"error": e.Error()})
	}
	query = `
	CREATE PROCEDURE drop_all_tables()
	BEGIN
	DECLARE done INT DEFAULT FALSE;
	DECLARE tableName VARCHAR(255);
	DECLARE cursor_tables CURSOR FOR
		SELECT table_name FROM information_schema.tables 
		WHERE table_type = 'BASE TABLE' AND table_schema = DATABASE();
	DECLARE CONTINUE HANDLER FOR NOT FOUND SET done = TRUE;
	OPEN cursor_tables;
	read_loop: LOOP
		FETCH cursor_tables INTO tableName;
		IF done THEN
		LEAVE read_loop;
		END IF;
		SET @drop_sql = CONCAT('DROP TABLE IF EXISTS ', tableName, ';');
		PREPARE stmt_drop FROM @drop_sql;
		EXECUTE stmt_drop;
		DEALLOCATE PREPARE stmt_drop;
	END LOOP;
	CLOSE cursor_tables;
	END;
	`
	e = ds.execSQLText(query)
	if e != nil {
		return
	}
}

func (s *db_store) addSPDropTable() {
	query := `DROP PROCEDURE IF EXISTS drop_table ;`
	e := s.execSQLText(query)
	if e != nil {
		return
	}
	query = `
	CREATE PROCEDURE drop_table(IN table_name VARCHAR(255))
	BEGIN
		SET @drop_sql = CONCAT('DROP TABLE IF EXISTS ', table_name, ';');

		PREPARE stmt_drop FROM @drop_sql;
		EXECUTE stmt_drop;
		DEALLOCATE PREPARE stmt_drop;
	END;
	`
	e = s.execSQLText(query)
	if e != nil {
		return
	}
}

func (s *db_store) addSPSearchTables() {
	query := `DROP PROCEDURE IF EXISTS search_tables ;`
	e := s.execSQLText(query)
	if e != nil {
		return
	}
	query = `
	CREATE PROCEDURE search_tables(IN pattern VARCHAR(255))
	BEGIN
		SET @sql = CONCAT('SELECT table_name FROM information_schema.tables WHERE table_name LIKE ''%', pattern, '%''');
		PREPARE stmt FROM @sql;
		EXECUTE stmt;
		DEALLOCATE PREPARE stmt;
	END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		return
	}
}

func (s *db_store) addSPUpsertUser() {
	query := `DROP PROCEDURE IF EXISTS upsert_user;`
	e := s.execSQLText(query)
	if e != nil {
		return
	}
	query = `
	CREATE PROCEDURE upsert_user(
		IN agent_id VARCHAR(36),
		IN their_uid VARCHAR(36),
		IN their_uname VARCHAR(20),
		IN their_ugrant VARCHAR(2),
		IN modified_at DATETIME
	)
	BEGIN
		DECLARE total INT;
	
		SELECT COUNT(*) INTO total FROM user_list WHERE AgentID = agent_id AND TheirUID = their_uid ;
	
		IF total > 0 THEN
			-- 資料存在，執行更新
			UPDATE user_list SET TheirUName = their_uname, ModifiedAt = modified_at WHERE AgentID = agent_id AND TheirUID = their_uid ;
		ELSE
			-- 資料不存在，執行插入
			INSERT INTO user_list (AgentID, TheirUID, TheirUName, TheirUGrant, CreatedAt, ModifiedAt)
			VALUES (agent_id, their_uid, their_uname, their_ugrant, modified_at, modified_at);
		END IF;
	
		SELECT ID FROM user_list WHERE AgentID = agent_id AND TheirUID = their_uid ;
	
	END ;
	`
	e = s.execSQLText(query)
	if e != nil {
		return
	}
}
